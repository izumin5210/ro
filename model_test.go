package ro

import (
	"testing"
)

func TestModel_GetKeySuffix(t *testing.T) {
	m := &Model{}
	if got, want := m.GetKeySuffix(), ""; got != want {
		t.Errorf("Model.GetKeySuffix() returns %q, want %q in default", got, want)
	}
}

func TestModel_GetScoreMap(t *testing.T) {
	m := &Model{}
	scoreMap := m.GetScoreMap()
	if scoreMap == nil {
		t.Error("Model.GetScoreMap() should not be nil")
	}
	if got, want := len(scoreMap), 0; got != want {
		t.Errorf("Model.GetScoreMap() has %d items, want %d items in default", got, want)
	}
}
