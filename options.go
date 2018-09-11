package ro

import (
	"reflect"
)

// Config contains configurations of a store
type Config struct {
	KeyPrefix             string
	ScoreSetKeysKeySuffix string
	KeyDelimiter          string
	ScoreKeyDelimiter     string
	HashStoreEnabled      bool
}

// Option configures a store
type Option func(c *Config)

func createConfig(modelType reflect.Type, opts []Option) *Config {
	cfg := &Config{
		KeyPrefix:             modelType.Name(),
		ScoreSetKeysKeySuffix: "scoreSetKeys",
		KeyDelimiter:          ":",
		ScoreKeyDelimiter:     "/",
		HashStoreEnabled:      true,
	}

	for _, f := range opts {
		f(cfg)
	}

	return cfg
}

// WithKeyPrefix returns a StoreOption that specifies key prefix
// If you does not set this option or set an empty string, it will use a model type name as key prefix.
func WithKeyPrefix(prefix string) Option {
	return func(c *Config) {
		c.KeyPrefix = prefix
	}
}

// WithScoreSetKeysKeySuffix returns a StoreOption that specifies a score set suffix key prefix (default: scoreSetKeys).
func WithScoreSetKeysKeySuffix(suffix string) Option {
	return func(c *Config) {
		c.ScoreSetKeysKeySuffix = suffix
	}
}

// WithKeyDelimiter returns a StoreOption that specifies a key delimiter (default: :).
func WithKeyDelimiter(d string) Option {
	return func(c *Config) {
		c.KeyDelimiter = d
	}
}

// WithScoreKeyDelimiter returns a StoreOption that specifies a score key delimiter (default: /).
func WithScoreKeyDelimiter(d string) Option {
	return func(c *Config) {
		c.ScoreKeyDelimiter = d
	}
}

// WithHashStore returns a StoreOption that enables or disables to store models into redis hash (default: true).
func WithHashStore(enabled bool) Option {
	return func(c *Config) {
		c.HashStoreEnabled = enabled
	}
}
