package store

import (
	"fmt"
	"reflect"

	"github.com/garyburd/redigo/redis"
	"github.com/izumin5210/ro/internal/query"
	"github.com/izumin5210/ro/types"
)

// Select implements the types.Store interface.
func (s *ConcreteStore) Select(dest interface{}, query types.Query) error {
	dt := reflect.ValueOf(dest)
	if dt.Kind() != reflect.Ptr || dt.IsNil() {
		return fmt.Errorf("must pass a slice ptr")
	}
	dt = dt.Elem()
	if dt.Kind() != reflect.Slice {
		return fmt.Errorf("must pass a slice ptr")
	}

	conn := s.getConn()
	defer conn.Close()

	cmd, args := query.Build()
	keys, err := redis.Strings(conn.Do(cmd, args...))
	if err != nil {
		return err
	}

	for _, key := range keys {
		err := conn.Send("HGETALL", key)
		if err != nil {
			return err
		}
	}

	conn.Flush()

	vt := dt.Type().Elem().Elem()

	for _ = range keys {
		v, err := redis.Values(conn.Receive())
		if err != nil {
			return err
		}
		vv := reflect.New(vt)
		err = redis.ScanStruct(v, vv.Interface())
		if err != nil {
			return err
		}
		dt.Set(reflect.Append(dt, vv))
	}

	return nil
}

// Query implements the types.Store interface.
func (s *ConcreteStore) Query(key string) types.Query {
	k := s.getKeyPrefix(s.model) + scoreDelimiter + key
	return query.New(k)
}
