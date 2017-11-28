package config

import (
	"github.com/creasty/defaults"
	"github.com/pkg/errors"

	"github.com/izumin5210/ro/types"
)

// New creates a new config object from default values and options.
func New(opts ...types.StoreOption) (*types.StoreConfig, error) {
	cnf := &types.StoreConfig{}
	if err := defaults.Set(cnf); err != nil {
		return nil, errors.Wrap(err, "failed to set default config")
	}
	for _, opt := range opts {
		cnf = opt(cnf)
	}
	return cnf, nil
}
