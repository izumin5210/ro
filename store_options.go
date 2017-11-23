package ro

import (
	"github.com/izumin5210/ro/types"
)

// WithScorers returns a StoreOption that sets scorer functions for storing into zset.
func WithScorers(funcs []types.ScorerFunc) types.StoreOption {
	return func(c *types.StoreConfig) *types.StoreConfig {
		c.ScorerFuncs = funcs
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

// WithScoreSetKeysKeySuffix returns a StoreOption that specifies a score set suffix key prefix (default: scoreSetKeys).
func WithScoreSetKeysKeySuffix(suffix string) types.StoreOption {
	return func(c *types.StoreConfig) *types.StoreConfig {
		c.ScoreSetKeysKeySuffix = suffix
		return c
	}
}

// WithKeyDelimiter returns a StoreOption that specifies a key delimiter (default: :).
func WithKeyDelimiter(d string) types.StoreOption {
	return func(c *types.StoreConfig) *types.StoreConfig {
		c.KeyDelimiter = d
		return c
	}
}

// WithScoreKeyDelimiter returns a StoreOption that specifies a score key delimiter (default: /).
func WithScoreKeyDelimiter(d string) types.StoreOption {
	return func(c *types.StoreConfig) *types.StoreConfig {
		c.ScoreKeyDelimiter = d
		return c
	}
}
