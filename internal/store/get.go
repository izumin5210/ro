package store

import (
	"github.com/garyburd/redigo/redis"
	"github.com/izumin5210/ro/types"
)

// Get implements the types.Store interface.
func (s *ConcreteStore) Get(dests ...types.Model) error {
	var err error

	conn := s.getConn()
	defer conn.Close()

	for _, m := range dests {
		err = conn.Send("HGETALL", s.getKey(m))
		if err != nil {
			return err
		}
	}

	err = conn.Flush()
	if err != nil {
		return err
	}

	for _, d := range dests {
		v, err := redis.Values(conn.Receive())
		if err != nil {
			return err
		}
		err = redis.ScanStruct(v, d)
		if err != nil {
			return err
		}
	}

	return err
}
