package store

import (
	"reflect"

	"github.com/garyburd/redigo/redis"
)

// Remove implements the types.Store interface.
func (s *ConcreteStore) Remove(src interface{}) error {
	var err error
	var prefix string

	keys := []string{}

	rv := reflect.ValueOf(src)
	if rv.Kind() == reflect.Slice {
		for i := 0; i < rv.Len(); i++ {
			m, err := s.toModel(rv.Index(i))
			if err != nil {
				return err
			}
			prefix = s.getKeyPrefix(m)

			keys = append(keys, s.getKey(m))
		}
	} else {
		m, err := s.toModel(rv)
		if err != nil {
			return err
		}
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
