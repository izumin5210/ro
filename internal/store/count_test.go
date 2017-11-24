package store

import (
	"testing"
	"time"

	"github.com/izumin5210/ro/types"
)

func TestCount(t *testing.T) {
	defer teardown(t)

	cnf := &types.StoreConfig{
		ScorerFuncs: []types.ScorerFunc{
			func(m types.Model) (string, interface{}) {
				return "id", m.(*TestPost).ID
			},
			func(m types.Model) (string, interface{}) {
				return "recent", m.(*TestPost).UpdatedAt
			},
		},
	}
	store, err := New(redisPool.Get, &TestPost{}, cnf)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	now := time.Now().UTC()
	posts := []*TestPost{
		{
			ID:        1,
			Title:     "post 1",
			Body:      "This is a post 1.",
			UpdatedAt: now.UnixNano(),
		},
		{
			ID:        2,
			Title:     "post 2",
			Body:      "This is a post 2.",
			UpdatedAt: now.Add(-1 * 60 * 60 * time.Second).UnixNano(),
		},
		{
			ID:        3,
			Title:     "post 3",
			Body:      "This is a post 3.",
			UpdatedAt: now.Add(1 * 60 * 60 * time.Second).UnixNano(),
		},
		{
			ID:        4,
			Title:     "post 4",
			Body:      "This is a post 4.",
			UpdatedAt: now.Add(2 * 60 * 60 * time.Second).UnixNano(),
		},
		{
			ID:        5,
			Title:     "post 5",
			Body:      "This is a post 5.",
			UpdatedAt: now.Add(-2 * 60 * 60 * time.Second).UnixNano(),
		},
	}

	for _, p := range posts {
		err = store.Set(p)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	}

	conn := redisPool.Get()
	defer conn.Close()

	cases := []struct {
		name  string
		q     types.Query
		count int
	}{
		{
			name:  "id with no query params",
			q:     store.Query("id"),
			count: 5,
		},
		{
			name:  "recent with no query params",
			q:     store.Query("recent"),
			count: 5,
		},
		{
			name:  "should ignore reverse",
			q:     store.Query("id").Reverse(),
			count: 5,
		},
		{
			name:  "should ignore limit",
			q:     store.Query("id").Limit(2),
			count: 5,
		},
		{
			name:  "should ignore offset",
			q:     store.Query("id").Offset(3),
			count: 5,
		},
		{
			name:  "with Gt",
			q:     store.Query("recent").Gt(now.UnixNano()),
			count: 2,
		},
		{
			name:  "with GtEq",
			q:     store.Query("recent").GtEq(now.UnixNano()),
			count: 3,
		},
		{
			name:  "with Lt",
			q:     store.Query("recent").Gt(now.UnixNano()),
			count: 2,
		},
		{
			name:  "with LtEq",
			q:     store.Query("recent").GtEq(now.UnixNano()),
			count: 3,
		},
		{
			name: "with GtEq and Lt",
			q: store.Query("recent").
				GtEq(now.Add(-1 * 60 * 60 * time.Second).UnixNano()).
				Lt(now.Add(1 * 60 * 60 * time.Second).UnixNano()),
			count: 2,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cnt, err := store.Count(c.q)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if got, want := cnt, c.count; got != want {
				t.Errorf("Count() returned %d, want %d", got, want)
			}
		})
	}
}
