package store

import (
	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"

	"github.com/izumin5210/ro/types"
)

// Get implements the types.Store interface.
func (s *ConcreteStore) Get(dests ...types.Model) error {
	var err error

	conn := s.getConn()
	defer conn.Close()

	for _, m := range dests {
		key := s.getKey(m)
		err = conn.Send("HGETALL", key)
		if err != nil {
			return errors.Wrapf(err, "faild to send HGETALL %s", key)
		}
	}

	err = conn.Flush()
	if err != nil {
		return errors.Wrap(err, "faild to flush HGETALL commands")
	}

	for _, d := range dests {
		v, err := redis.Values(conn.Receive())
		if err != nil {
			return errors.Wrap(err, "faild to receive or cast redis command result")
		}
		err = redis.ScanStruct(v, d)
		if err != nil {
			return errors.Wrapf(err, "faild to scan struct %s %x", s.getKey(d), v)
		}
	}

	return err
}
