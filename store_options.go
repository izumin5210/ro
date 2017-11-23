package ro

import (
	"github.com/izumin5210/ro/types"
)

// WithScorer returns a StoreOption that sets scorer functions for storing into zset.
func WithScorer(m map[string]types.ScorerFunc) types.StoreOption {
	return func(c *types.StoreConfig) *types.StoreConfig {
		c.ScorerFuncMap = m
		return c
	}
}

// WithKeyPrefix returns a StoreOption that specifies key prefix
// If you does not set this option or set an empty string, it will use a model type name as key prefix.
func WithKeyPrefix(prefix string) types.StoreOption {
	return func(c *types.StoreConfig) *types.StoreConfig {
		c.KeyPrefix = prefix
		return c
	}
}
