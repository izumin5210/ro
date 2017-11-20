package ro

import (
	"github.com/izumin5210/ro/internal/query"
	"github.com/izumin5210/ro/types"
)

// Query creates new query object
func Query(key string) types.Query {
	return query.New(key)
}
