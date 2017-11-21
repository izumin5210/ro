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
	Select(dest interface{}, query Query) error
	Query(key string) Query
}

// StoreConfig contains configurations of a store
type StoreConfig struct {
	ScorerFuncMap map[string]ScorerFunc
}

// StoreOption configures a store
type StoreOption func(c *StoreConfig) *StoreConfig

// ScorerFunc is an adapteer to calculate score from given model
type ScorerFunc func(Model) interface{}

// Query is an interface to set conditions to find stored objects
type Query interface {
	Gt(v interface{}) Query
	GtEq(v interface{}) Query
	Lt(v interface{}) Query
	LtEq(v interface{}) Query
	Limit(v int) Query
	Offset(v int) Query
	Reverse() Query
	Build() (string, []interface{})
}
