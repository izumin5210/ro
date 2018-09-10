package types

import (
	"github.com/gomodule/redigo/redis"
	"github.com/izumin5210/ro/rq"
)

// GetConnFunc retrievs a connection object with redis
type GetConnFunc func() redis.Conn

// Model is an interface for redis objects
type Model interface {
	GetKeySuffix() string
	GetScoreMap() map[string]interface{}
}

// Store is an interface for providing CRUD operations for objects
type Store interface {
	Set(src interface{}) error
	Get(dests ...Model) error
	Select(dest interface{}, mods ...rq.Modifier) error
	Count(mods ...rq.Modifier) (int, error)
	Remove(src interface{}) error
	RemoveBy(mods ...rq.Modifier) error
}

// StoreConfig contains configurations of a store
type StoreConfig struct {
	KeyPrefix             string
	ScoreSetKeysKeySuffix string `default:"scoreSetKeys"`
	KeyDelimiter          string `default:":"`
	ScoreKeyDelimiter     string `default:"/"`
	HashStoreEnabled      bool   `default:"true"`
}

// StoreOption configures a store
type StoreOption func(c *StoreConfig) *StoreConfig

// ScorerFunc is an adapteer to calculate score from given model
type ScorerFunc func(Model) (string, interface{})
