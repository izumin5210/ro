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
