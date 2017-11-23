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
func WithKeyPrefix(prefix string) types.StoreOption {
	return func(c *types.StoreConfig) *types.StoreConfig {
		c.KeyPrefix = prefix
		return c
	}
}
