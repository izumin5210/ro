package ro

import (
	"github.com/izumin5210/ro/internal/store"
	"github.com/izumin5210/ro/types"
)

// Store is an interface for providing CRUD operations for objects
type Store interface {
	types.Store
}

// NewStore creates a new store instance for given model objects
func NewStore(getConnFunc types.GetConnFunc, model types.Model, opts ...types.StoreOption) (Store, error) {
	cnf := &types.StoreConfig{}
	for _, opt := range opts {
		cnf = opt(cnf)
	}
	return store.New(getConnFunc, model, cnf)
}
