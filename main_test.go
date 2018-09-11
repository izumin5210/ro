package ro

import (
	"log"
	"os"
	"testing"

	rotesting "github.com/izumin5210/ro/testing"
)

var pool *rotesting.Pool

func TestMain(m *testing.M) {
	pool = rotesting.MustCreate()
	defer pool.MustClose()

	os.Exit(m.Run())
}

func teardown(t *testing.T) {
	if err := pool.Cleanup(); err != nil {
		log.Fatalf("Failed to flush redis: %s", err)
	}
}
