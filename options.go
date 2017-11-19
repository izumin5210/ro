package ro

import (
	"github.com/izumin5210/ro/types"
)

func WithScorer(m map[string]types.ScorerFunc) types.StoreOption {
	return func(c *types.StoreConfig) *types.StoreConfig {
		c.ScorerFuncMap = m
		return c
	}
}
