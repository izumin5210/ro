package store

import (
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"

	"github.com/izumin5210/ro/rq"
)

// Count implements the types.Store interface.
func (s *ConcreteStore) Count(mods ...rq.Modifier) (int, error) {
	conn := s.getConn()
	defer conn.Close()

	cmd, err := s.injectKeyPrefix(rq.Count(mods...)).Build()
	if err != nil {
		return 0, errors.WithStack(err)
	}

	cnt, err := redis.Int(conn.Do(cmd.Name, cmd.Args...))
	if err != nil {
		return 0, errors.Wrapf(err, "faild to execute %v", cmd)
	}
	return cnt, nil
}
