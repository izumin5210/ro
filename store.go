package ro

import (
	"reflect"

	"github.com/gomodule/redigo/redis"

	"github.com/izumin5210/ro/rq"
)

// Store is an interface for providing CRUD operations for objects
type Store interface {
	List(dest interface{}, mods ...rq.Modifier) error
	Get(dests ...Model) error
	Put(src interface{}) error
	Delete(src interface{}) error
	DeleteAll(mods ...rq.Modifier) error
	Count(mods ...rq.Modifier) (int, error)
}

// Pool is a pool of redis connections.
type Pool interface {
	Get() redis.Conn
}

type redisStore struct {
	*Config
	pool      Pool
	model     Model
	modelType reflect.Type
}

// New creates a redisStore instance
func New(pool Pool, model Model, opts ...Option) Store {
	modelType := reflect.ValueOf(model).Elem().Type()

	return &redisStore{
		Config:    createConfig(modelType, opts),
		pool:      pool,
		model:     model,
		modelType: modelType,
	}
}
