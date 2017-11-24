package store

import (
	"github.com/garyburd/redigo/redis"

	"github.com/izumin5210/ro/types"
)

// Count implements the types.Store interface.
func (s *ConcreteStore) Count(query types.Query) (int, error) {
	conn := s.getConn()
	defer conn.Close()

	cmd, args := query.BuildForCount()

	return redis.Int(conn.Do(cmd, args...))
}
