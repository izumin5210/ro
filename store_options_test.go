package ro

import (
	"testing"

	"github.com/izumin5210/ro/types"
)

func Test_WithKeyPrefix(t *testing.T) {
	cnf := &types.StoreConfig{}
	if got, want := "", cnf.KeyPrefix; got != want {
		t.Errorf("StoreConfig.KeyPrefix is %q, want %q", got, want)
	}
	prefix := "newprefix"
	cnf = WithKeyPrefix(prefix)(cnf)
	if got, want := prefix, cnf.KeyPrefix; got != want {
		t.Errorf("StoreConfig.KeyPrefix is %q, want %q", got, want)
	}
}

func Test_WithScoreSetKeysKeySuffix(t *testing.T) {
	cnf := &types.StoreConfig{}
	if got, want := "", cnf.ScoreSetKeysKeySuffix; got != want {
		t.Errorf("StoreConfig.ScoreSetKeysKeySuffix is %q, want %q", got, want)
	}
	suffix := "newsuffix"
	cnf = WithScoreSetKeysKeySuffix(suffix)(cnf)
	if got, want := suffix, cnf.ScoreSetKeysKeySuffix; got != want {
		t.Errorf("StoreConfig.ScoreSetKeysKeySuffix is %q, want %q", got, want)
	}
}

func Test_WithKeyDelimiter(t *testing.T) {
	cnf := &types.StoreConfig{}
	if got, want := "", cnf.KeyDelimiter; got != want {
		t.Errorf("StoreConfig.KeyDelimiter is %q, want %q", got, want)
	}
	delimiter := "::"
	cnf = WithKeyDelimiter(delimiter)(cnf)
	if got, want := delimiter, cnf.KeyDelimiter; got != want {
		t.Errorf("StoreConfig.KeyDelimiter is %q, want %q", got, want)
	}
}

func Test_WithScoreKeyDelimiter(t *testing.T) {
	cnf := &types.StoreConfig{}
	if got, want := "", cnf.ScoreKeyDelimiter; got != want {
		t.Errorf("StoreConfig.ScoreKeyDelimiter is %q, want %q", got, want)
	}
	delimiter := "//"
	cnf = WithScoreKeyDelimiter(delimiter)(cnf)
	if got, want := delimiter, cnf.ScoreKeyDelimiter; got != want {
		t.Errorf("StoreConfig.ScoreKeyDelimiter is %q, want %q", got, want)
	}
}

func Test_WithHashStore(t *testing.T) {
	cnf := &types.StoreConfig{}
	if got, want := false, cnf.HashStoreEnabled; got != want {
		t.Errorf("StoreConfig.HashStoreEnabled is %t, want %t", got, want)
	}
	enabled := true
	cnf = WithHashStore(enabled)(cnf)
	if got, want := enabled, cnf.HashStoreEnabled; got != want {
		t.Errorf("StoreConfig.HashStoreEnabled is %t, want %t", got, want)
	}
}
