package store

import (
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"

	"github.com/izumin5210/ro/types"
)

// Count implements the types.Store interface.
func (s *ConcreteStore) Count(query types.Query) (int, error) {
	conn := s.getConn()
	defer conn.Close()

	cmd, args := query.BuildForCount()

	cnt, err := redis.Int(conn.Do(cmd, args...))
	if err != nil {
		return 0, errors.Wrapf(err, "faild to execute %s %v", cmd, args)
	}
	return cnt, nil
}
