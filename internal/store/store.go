package store

import (
	"reflect"

	"github.com/izumin5210/ro/internal/query"
	"github.com/izumin5210/ro/types"
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

// Query implements the types.Store interface.
func (s *ConcreteStore) Query(key string) types.Query {
	return query.New(s.getScoreSetKey(key))
}
