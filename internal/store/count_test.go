package store

import (
	"testing"
	"time"

	"github.com/izumin5210/ro/internal/config"
	rotesting "github.com/izumin5210/ro/internal/testing"
	"github.com/izumin5210/ro/rq"
)

func TestCount(t *testing.T) {
	defer teardown(t)

	cnf, _ := config.New()
	store, err := New(pool.Get, &rotesting.Post{}, cnf)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	now := time.Now().UTC()
	posts := []*rotesting.Post{
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

	conn := pool.Get()
	defer conn.Close()

	cases := []struct {
		name  string
		mods  []rq.Modifier
		count int
	}{
		{
			name:  "id with no query params",
			mods:  []rq.Modifier{rq.Key("id")},
			count: 5,
		},
		{
			name:  "recent with no query params",
			mods:  []rq.Modifier{rq.Key("recent")},
			count: 5,
		},
		{
			name:  "should ignore reverse",
			mods:  []rq.Modifier{rq.Key("id"), rq.Reverse()},
			count: 5,
		},
		{
			name:  "should ignore limit",
			mods:  []rq.Modifier{rq.Key("id"), rq.Limit(2)},
			count: 5,
		},
		{
			name:  "should ignore offset",
			mods:  []rq.Modifier{rq.Key("id"), rq.Offset(3)},
			count: 5,
		},
		{
			name:  "with Gt",
			mods:  []rq.Modifier{rq.Key("recent"), rq.Gt(now.UnixNano())},
			count: 2,
		},
		{
			name:  "with GtEq",
			mods:  []rq.Modifier{rq.Key("recent"), rq.GtEq(now.UnixNano())},
			count: 3,
		},
		{
			name:  "with Lt",
			mods:  []rq.Modifier{rq.Key("recent"), rq.Lt(now.UnixNano())},
			count: 2,
		},
		{
			name:  "with LtEq",
			mods:  []rq.Modifier{rq.Key("recent"), rq.LtEq(now.UnixNano())},
			count: 3,
		},
		{
			name: "with GtEq and Lt",
			mods: []rq.Modifier{
				rq.Key("recent"),
				rq.GtEq(now.Add(-1 * 60 * 60 * time.Second).UnixNano()),
				rq.Lt(now.Add(1 * 60 * 60 * time.Second).UnixNano()),
			},
			count: 2,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cnt, err := store.Count(c.mods...)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if got, want := cnt, c.count; got != want {
				t.Errorf("Count() returned %d, want %d", got, want)
			}
		})
	}
}
