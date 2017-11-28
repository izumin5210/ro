package ro

import (
	"github.com/pkg/errors"

	"github.com/izumin5210/ro/internal/config"
	"github.com/izumin5210/ro/internal/store"
	"github.com/izumin5210/ro/types"
)

// Store is an interface for providing CRUD operations for objects
type Store interface {
	types.Store
}

// New creates a new store instance for given model objects
func New(getConnFunc types.GetConnFunc, model types.Model, opts ...types.StoreOption) (Store, error) {
	cnf, err := config.New(opts...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize a config")
	}
	return store.New(getConnFunc, model, cnf)
}
