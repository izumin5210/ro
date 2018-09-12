package ro_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/izumin5210/ro"
	"github.com/izumin5210/ro/rq"
	rotesting "github.com/izumin5210/ro/testing"
)

func TestRedisStore_List(t *testing.T) {
	defer teardown(t)

	store := ro.New(pool, &rotesting.Post{})

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
		err := store.Put(context.TODO(), p)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	}

	conn := pool.Get()
	defer conn.Close()

	cases := []struct {
		name  string
		mods  []rq.Modifier
		order []int
	}{
		{
			name:  "id with no query params",
			mods:  []rq.Modifier{rq.Key("id")},
			order: []int{0, 1, 2, 3, 4},
		},
		{
			name:  "recent with no query params",
			mods:  []rq.Modifier{rq.Key("recent")},
			order: []int{4, 1, 0, 2, 3},
		},
		{
			name:  "id with reverse",
			mods:  []rq.Modifier{rq.Key("id"), rq.Reverse()},
			order: []int{4, 3, 2, 1, 0},
		},
		{
			name:  "id with limit",
			mods:  []rq.Modifier{rq.Key("id"), rq.Limit(2)},
			order: []int{0, 1},
		},
		{
			name:  "id with offset",
			mods:  []rq.Modifier{rq.Key("id"), rq.Offset(3)},
			order: []int{3, 4},
		},
		{
			name:  "id with limit and offset",
			mods:  []rq.Modifier{rq.Key("id"), rq.Offset(1), rq.Limit(3)},
			order: []int{1, 2, 3},
		},
		{
			name:  "recent with limit and offset",
			mods:  []rq.Modifier{rq.Key("recent"), rq.Offset(1), rq.Limit(3)},
			order: []int{1, 0, 2},
		},
		{
			name:  "recent with Gt",
			mods:  []rq.Modifier{rq.Key("recent"), rq.Gt(now.UnixNano())},
			order: []int{2, 3},
		},
		{
			name:  "recent with GtEq",
			mods:  []rq.Modifier{rq.Key("recent"), rq.GtEq(now.UnixNano())},
			order: []int{0, 2, 3},
		},
		{
			name:  "recent with Lt",
			mods:  []rq.Modifier{rq.Key("recent"), rq.Lt(now.UnixNano())},
			order: []int{4, 1},
		},
		{
			name:  "recent with LtEq",
			mods:  []rq.Modifier{rq.Key("recent"), rq.LtEq(now.UnixNano())},
			order: []int{4, 1, 0},
		},
		{
			name: "recent with GtEq and Lt",
			mods: []rq.Modifier{
				rq.Key("recent"),
				rq.GtEq(now.Add(-1 * 60 * 60 * time.Second).UnixNano()),
				rq.Lt(now.Add(1 * 60 * 60 * time.Second).UnixNano()),
			},
			order: []int{1, 0},
		},
		{
			name: "recent with GtEq and Lt and Reverse",
			mods: []rq.Modifier{
				rq.Key("recent"),
				rq.Gt(now.Add(-1 * 60 * 60 * time.Second).UnixNano()),
				rq.LtEq(now.Add(1 * 60 * 60 * time.Second).UnixNano()),
				rq.Reverse(),
			},
			order: []int{2, 0},
		},
		{
			name:  "recent with LtEq and Limit",
			mods:  []rq.Modifier{rq.Key("recent"), rq.LtEq(now.UnixNano()), rq.Limit(2)},
			order: []int{4, 1},
		},
		{
			name:  "recent with LtEq and Offset",
			mods:  []rq.Modifier{rq.Key("recent"), rq.LtEq(now.UnixNano()), rq.Offset(1)},
			order: []int{1, 0},
		},
		{
			name: "recent with all conditions",
			mods: []rq.Modifier{
				rq.Key("recent"),
				rq.Gt(now.Add(-2 * 60 * 60 * time.Second).UnixNano()),
				rq.Offset(1),
				rq.Limit(2),
			},
			order: []int{0, 2},
		},
		{
			name: "recent with all conditions",
			mods: []rq.Modifier{
				rq.Key("recent"),
				rq.Gt(now.Add(-2 * 60 * 60 * time.Second).UnixNano()),
				rq.LtEq(now.Add(2 * 60 * 60 * time.Second).UnixNano()),
				rq.Offset(1),
				rq.Limit(2),
				rq.Reverse(),
			},
			order: []int{2, 0},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotPosts := []*rotesting.Post{}
			err := store.List(context.TODO(), &gotPosts, c.mods...)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if got, want := len(gotPosts), len(c.order); got != want {
				t.Errorf("List() returned %d posts, want %d posts", got, want)
				return
			}

			for i, j := range c.order {
				if got, want := gotPosts[i], posts[j]; !reflect.DeepEqual(got, want) {
					cmd, _ := rq.List(c.mods...).Build()
					t.Errorf("List(%v)[%d] is %v, want %v", cmd, i, got, want)
				}
			}
		})
	}
}
