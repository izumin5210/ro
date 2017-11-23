package store

import (
	"fmt"
	"reflect"

	"github.com/garyburd/redigo/redis"
	"github.com/izumin5210/ro/types"
)

// Remove implements the types.Store interface.
func (s *ConcreteStore) Remove(src interface{}) error {
	var err error
	var prefix string

	keys := []string{}

	rv := reflect.ValueOf(src)
	if rv.Kind() == reflect.Slice {
		for i := 0; i < rv.Len(); i++ {
			rvv := rv.Index(i)
			if rvv.Type() != s.modelType && rvv.Type().Elem() != s.modelType {
				return fmt.Errorf("%s is not a %v", rvv.Interface(), s.modelType)
			}

			m := rvv.Interface().(types.Model)
			prefix = s.getKeyPrefix(m)

			keys = append(keys, s.getKey(m))
		}
	} else {
		if rv.Type() != s.modelType && rv.Type().Elem() != s.modelType {
			return fmt.Errorf("%s is not a %v", rv.Interface(), s.modelType)
		}

		m := rv.Interface().(types.Model)
		prefix = s.getKeyPrefix(m)

		keys = append(keys, s.getKey(m))
	}

	conn := s.getConn()
	defer conn.Close()

	err = conn.Send("MULTI")
	if err != nil {
		return err
	}

	err = conn.Send("DEL", redis.Args{}.AddFlat(keys)...)
	if err != nil {
		return err
	}

	for k := range s.ScorerFuncMap {
		err = conn.Send("ZREM", redis.Args{}.Add(prefix+scoreDelimiter+k).AddFlat(keys)...)
		if err != nil {
			return err
		}
	}

	_, err = conn.Do("EXEC")
	return err
}
