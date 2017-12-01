package store

import (
	"reflect"
	"testing"

	"github.com/garyburd/redigo/redis"

	"github.com/izumin5210/ro/internal/config"
	"github.com/izumin5210/ro/internal/testing"
)

func TestGet(t *testing.T) {
	defer teardown(t)
	cnf, _ := config.New()
	store, err := New(pool.Get, &rotesting.Post{}, cnf)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	posts := []*rotesting.Post{
		{
			ID:    1,
			Title: "post 1",
			Body:  "This is a post 1.",
		},
		{
			ID:    2,
			Title: "post 2",
			Body:  "This is a post 2.",
		},
	}

	conn := pool.Get()
	defer conn.Close()
	conn.Do("HMSET", redis.Args{}.Add("Post:1").AddFlat(posts[0])...)
	conn.Do("HMSET", redis.Args{}.Add("Post:2").AddFlat(posts[1])...)

	t.Run("single get", func(t *testing.T) {
		gotPost := &rotesting.Post{ID: 2}
		err = store.Get(gotPost)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if got, want := gotPost, posts[1]; !reflect.DeepEqual(got, want) {
			t.Errorf("Stored post is %v, want %v", got, want)
		}
	})

	t.Run("multi get", func(t *testing.T) {
		gotPosts := []*rotesting.Post{{ID: 2}, {ID: 1}}
		err = store.Get(gotPosts[0], gotPosts[1])

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if got, want := gotPosts[0], posts[1]; !reflect.DeepEqual(got, want) {
			t.Errorf("Stored post is %v, want %v", got, want)
		}
		if got, want := gotPosts[1], posts[0]; !reflect.DeepEqual(got, want) {
			t.Errorf("Stored post is %v, want %v", got, want)
		}
	})
}
