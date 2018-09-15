package ro

import (
	"context"
	"reflect"

	"github.com/gomodule/redigo/redis"

	"github.com/izumin5210/ro/rq"
)

// Store is an interface for providing CRUD operations for objects
type Store interface {
	List(ctx context.Context, dest interface{}, mods ...rq.Modifier) error
	Get(ctx context.Context, dests ...Model) error
	Put(ctx context.Context, src interface{}) error
	PutConn(conn redis.Conn, src interface{}) error
	Delete(ctx context.Context, src interface{}) error
	DeleteAll(ctx context.Context, mods ...rq.Modifier) error
	Count(ctx context.Context, mods ...rq.Modifier) (int, error)
}

// Pool is a pool of redis connections.
type Pool interface {
	GetContext(context.Context) (redis.Conn, error)
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
