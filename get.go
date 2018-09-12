package ro

import (
	"context"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
)

func (s *redisStore) Get(ctx context.Context, dests ...Model) error {
	conn, err := s.pool.GetContext(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to acquire a connection")
	}
	defer conn.Close()

	keys := make([]string, len(dests), len(dests))

	for i, m := range dests {
		key, err := s.getKey(m)
		if err != nil {
			return errors.Wrap(err, "failed to get key")
		}
		keys[i] = key
		err = conn.Send("HGETALL", key)
		if err != nil {
			return errors.Wrapf(err, "faild to send HGETALL %s", key)
		}
	}

	err = conn.Flush()
	if err != nil {
		return errors.Wrap(err, "faild to flush HGETALL commands")
	}

	for i, d := range dests {
		v, err := redis.Values(conn.Receive())
		if err != nil {
			return errors.Wrap(err, "faild to receive or cast redis command result")
		}
		err = redis.ScanStruct(v, d)
		if err != nil {
			return errors.Wrapf(err, "faild to scan struct %s %x", keys[i], v)
		}
	}

	return nil
}
