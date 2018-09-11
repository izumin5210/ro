package store

import (
	"reflect"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"

	"github.com/izumin5210/ro/rq"
)

// Select implements the types.Store interface.
func (s *ConcreteStore) Select(dest interface{}, mods ...rq.Modifier) error {
	dt := reflect.ValueOf(dest)
	if dt.Kind() != reflect.Ptr || dt.IsNil() {
		return errors.New("must pass a slice ptr")
	}
	dt = dt.Elem()
	if dt.Kind() != reflect.Slice {
		return errors.New("must pass a slice ptr")
	}

	keys, err := s.selectKeys(mods)
	if err != nil {
		return errors.Wrap(err, "failed to select query")
	}

	conn := s.getConn()
	defer conn.Close()

	for _, key := range keys {
		err := conn.Send("HGETALL", key)
		if err != nil {
			return errors.Wrapf(err, "failed to send HGETALL %s", key)
		}
	}

	conn.Flush()

	vt := dt.Type().Elem().Elem()

	for _, key := range keys {
		v, err := redis.Values(conn.Receive())
		if err != nil {
			return errors.Wrap(err, "faild to receive or cast redis command result")
		}
		vv := reflect.New(vt)
		err = redis.ScanStruct(v, vv.Interface())
		if err != nil {
			return errors.Wrapf(err, "faild to scan struct %s %x", key, v)
		}
		dt.Set(reflect.Append(dt, vv))
	}

	return nil
}
