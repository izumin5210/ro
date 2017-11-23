package store

import (
	"reflect"

	"github.com/garyburd/redigo/redis"

	"github.com/izumin5210/ro/internal/query"
	"github.com/izumin5210/ro/types"
)

const (
	keyDelimiter   = ":"
	scoreDelimiter = "/"
)

// ConcreteStore is an implementation of types.Store
type ConcreteStore struct {
	*types.StoreConfig
	getConn   types.GetConnFunc
	model     types.Model
	modelType reflect.Type
}

// New creates a ConcreteStore instance
func New(getConnFunc types.GetConnFunc, model types.Model, cnf *types.StoreConfig) (types.Store, error) {
	modelType := reflect.ValueOf(model).Elem().Type()

	if len(cnf.KeyPrefix) == 0 {
		cnf.KeyPrefix = modelType.Name()
	}

	return &ConcreteStore{
		StoreConfig: cnf,
		getConn:     getConnFunc,
		model:       model,
		modelType:   modelType,
	}, nil
}

// Get implements the types.Store interface.
func (s *ConcreteStore) Get(dest types.Model) error {
	conn := s.getConn()
	defer conn.Close()

	v, err := redis.Values(conn.Do("HGETAll", s.getKey(dest)))
	if err != nil {
		return err
	}
	err = redis.ScanStruct(v, dest)
	if err != nil {
		return err
	}
	return nil
}

// Query implements the types.Store interface.
func (s *ConcreteStore) Query(key string) types.Query {
	k := s.KeyPrefix + scoreDelimiter + key
	return query.New(k)
}
