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

	keys, err := s.keys(query)
	if err != nil {
		return err
	}

	conn := s.getConn()
	defer conn.Close()

	for _, key := range keys {
		err := conn.Send("HGETALL", key)
		if err != nil {
			return err
		}
	}

	conn.Flush()

	vt := dt.Type().Elem().Elem()

	for range keys {
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

// Count implements the types.Store interface.
func (s *ConcreteStore) Count(query types.Query) (int, error) {
	keys, err := s.keys(query)
	if err != nil {
		return 0, err
	}
	return len(keys), nil
}

// Query implements the types.Store interface.
func (s *ConcreteStore) Query(key string) types.Query {
	k := s.KeyPrefix + scoreDelimiter + key
	return query.New(k)
}

func (s *ConcreteStore) keys(query types.Query) ([]string, error) {
	conn := s.getConn()
	defer conn.Close()

	cmd, args := query.Build()
	return redis.Strings(conn.Do(cmd, args...))
}
