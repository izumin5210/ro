package store

import (
	"reflect"
	"testing"
	"time"

	"github.com/izumin5210/ro/types"
)

func TestSelect(t *testing.T) {
	defer teardown(t)

	cnf := &types.StoreConfig{
		ScorerFuncMap: map[string]types.ScorerFunc{
			"id": func(m types.Model) interface{} {
				return m.(*TestPost).ID
			},
			"recent": func(m types.Model) interface{} {
				return m.(*TestPost).UpdatedAt
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
		order []int
	}{
		{
			name:  "id with no query params",
			q:     store.Query("id"),
			order: []int{0, 1, 2, 3, 4},
		},
		{
			name:  "recent with no query params",
			q:     store.Query("recent"),
			order: []int{4, 1, 0, 2, 3},
		},
		{
			name:  "id with reverse",
			q:     store.Query("id").Reverse(),
			order: []int{4, 3, 2, 1, 0},
		},
		{
			name:  "id with limit",
			q:     store.Query("id").Limit(2),
			order: []int{0, 1},
		},
		{
			name:  "id with offset",
			q:     store.Query("id").Offset(3),
			order: []int{3, 4},
		},
		{
			name:  "id with limit and offset",
			q:     store.Query("id").Offset(1).Limit(3),
			order: []int{1, 2, 3},
		},
		{
			name:  "recent with limit and offset",
			q:     store.Query("recent").Offset(1).Limit(3),
			order: []int{1, 0, 2},
		},
		{
			name:  "recent with Gt",
			q:     store.Query("recent").Gt(now.UnixNano()),
			order: []int{2, 3},
		},
		{
			name:  "recent with GtEq",
			q:     store.Query("recent").GtEq(now.UnixNano()),
			order: []int{0, 2, 3},
		},
		{
			name:  "recent with Lt",
			q:     store.Query("recent").Gt(now.UnixNano()),
			order: []int{2, 3},
		},
		{
			name:  "recent with LtEq",
			q:     store.Query("recent").GtEq(now.UnixNano()),
			order: []int{0, 2, 3},
		},
		{
			name: "recent with GtEq and Lt",
			q: store.Query("recent").
				GtEq(now.Add(-1 * 60 * 60 * time.Second).UnixNano()).
				Lt(now.Add(1 * 60 * 60 * time.Second).UnixNano()),
			order: []int{1, 0},
		},
		{
			name: "recent with GtEq and Lt and Reverse",
			q: store.Query("recent").
				Gt(now.Add(-1 * 60 * 60 * time.Second).UnixNano()).
				LtEq(now.Add(1 * 60 * 60 * time.Second).UnixNano()).
				Reverse(),
			order: []int{2, 0},
		},
		{
			name:  "recent with LtEq and Limit",
			q:     store.Query("recent").LtEq(now.UnixNano()).Limit(2),
			order: []int{4, 1},
		},
		{
			name:  "recent with LtEq and Offset",
			q:     store.Query("recent").LtEq(now.UnixNano()).Offset(1),
			order: []int{1, 0},
		},
		{
			name: "recent with all conditions",
			q: store.Query("recent").
				Gt(now.Add(-2 * 60 * 60 * time.Second).UnixNano()).
				Offset(1).
				Limit(2),
			order: []int{0, 2},
		},
		{
			name: "recent with all conditions",
			q: store.Query("recent").
				Gt(now.Add(-2 * 60 * 60 * time.Second).UnixNano()).
				LtEq(now.Add(2 * 60 * 60 * time.Second).UnixNano()).
				Offset(1).
				Limit(2).
				Reverse(),
			order: []int{2, 0},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotPosts := []*TestPost{}
			err = store.Select(&gotPosts, c.q)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if got, want := len(gotPosts), len(c.order); got != want {
				t.Errorf("Select() returned %d posts, want %d posts", got, want)
				return
			}

			for i, j := range c.order {
				if got, want := gotPosts[i], posts[j]; !reflect.DeepEqual(got, want) {
					cmd, args := c.q.Build()
					t.Errorf("Select(%v, %v)[%d] is %v, want %v", cmd, args, i, got, want)
				}
			}
		})
	}
}
