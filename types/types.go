package types

import (
	"github.com/garyburd/redigo/redis"
)

// GetConnFunc retrievs a connection object with redis
type GetConnFunc func() redis.Conn

// Model is an interface for redis objects
type Model interface {
	GetKeyPrefix() string
	GetKeySuffix() string
}

// Store is an interface for providing CRUD operations for objects
type Store interface {
	Set(src Model) error
	Get(dest Model) error
}
