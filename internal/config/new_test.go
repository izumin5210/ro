package config

import (
	"testing"

	"github.com/izumin5210/ro/types"
)

func TestNew(t *testing.T) {
	cases := []struct {
		name string
		in   []types.StoreOption
		out  *types.StoreConfig
	}{
		{
			name: "default",
			out: &types.StoreConfig{
				ScoreSetKeysKeySuffix: "scoreSetKeys",
				KeyDelimiter:          ":",
				ScoreKeyDelimiter:     "/",
			},
		},
		{
			name: "reset",
			in: []types.StoreOption{
				func(c *types.StoreConfig) *types.StoreConfig {
					return &types.StoreConfig{}
				},
			},
			out: &types.StoreConfig{},
		},
		{
			name: "with an option",
			in: []types.StoreOption{
				func(c *types.StoreConfig) *types.StoreConfig {
					c.KeyDelimiter = "::"
					return c
				},
			},
			out: &types.StoreConfig{
				ScoreSetKeysKeySuffix: "scoreSetKeys",
				KeyDelimiter:          "::",
				ScoreKeyDelimiter:     "/",
			},
		},
		{
			name: "with 2 options",
			in: []types.StoreOption{
				func(c *types.StoreConfig) *types.StoreConfig {
					c.KeyDelimiter = "::"
					return c
				},
				func(c *types.StoreConfig) *types.StoreConfig {
					c.ScoreKeyDelimiter = "//"
					return c
				},
			},
			out: &types.StoreConfig{
				ScoreSetKeysKeySuffix: "scoreSetKeys",
				KeyDelimiter:          "::",
				ScoreKeyDelimiter:     "//",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cnf, err := New(c.in...)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if got, want := cnf.ScoreSetKeysKeySuffix, c.out.ScoreSetKeysKeySuffix; got != want {
				t.Errorf("New().ScoreSetKeysKeySuffix returns %v, want %v", got, want)
			}
			if got, want := cnf.KeyDelimiter, c.out.KeyDelimiter; got != want {
				t.Errorf("New().KeyDelimiter returns %v, want %v", got, want)
			}
			if got, want := cnf.ScoreKeyDelimiter, c.out.ScoreKeyDelimiter; got != want {
				t.Errorf("New().ScoreKeyDelimiter returns %v, want %v", got, want)
			}
		})
	}
}
