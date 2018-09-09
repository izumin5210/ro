package store

import (
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/izumin5210/ro/internal/config"
	"github.com/izumin5210/ro/internal/testing"
)

func TestRemoveBy(t *testing.T) {
	defer teardown(t)
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
			UpdatedAt: now.Add(-1 * 60 * time.Second).UnixNano(),
		},
		{
			ID:        3,
			Title:     "post 3",
			Body:      "This is a post 3.",
			UpdatedAt: now.Add(1 * 60 * time.Second).UnixNano(),
		},
		{
			ID:        4,
			Title:     "post 4",
			Body:      "This is a post 4.",
			UpdatedAt: now.Add(-2 * 60 * time.Second).UnixNano(),
		},
		{
			ID:        5,
			Title:     "post 5",
			Body:      "This is a post 5.",
			UpdatedAt: now.Add(1 * 60 * time.Second).UnixNano(),
		},
	}

	cnf, _ := config.New()
	store, err := New(pool.Get, &rotesting.Post{}, cnf)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	err = store.Set(posts)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	q := store.Query("recent").
		Gt(now.Add(-2 * 60 * time.Second).UnixNano()).
		LtEq(now.Add(1 * 60 * time.Second).UnixNano())
	err = store.RemoveBy(q)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	conn := pool.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if got, want := len(keys), 4; err != nil {
		t.Errorf("Stored keys was %d, want %d", got, want)
	}

	v, err := redis.Values(conn.Do("HGETALL", "rotesting.Post:1"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(v) > 0 {
		t.Errorf("Unexpected response: %v", v)
	}

	v, err = redis.Values(conn.Do("HGETALL", "rotesting.Post:2"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(v) > 0 {
		t.Errorf("Unexpected response: %v", v)
	}

	v, err = redis.Values(conn.Do("HGETALL", "rotesting.Post:5"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(v) > 0 {
		t.Errorf("Unexpected response: %v", v)
	}

	keys, err = redis.Strings(conn.Do("ZRANGE", "TestPoset/recent", 0, -1))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if got, want := len(keys), 2; err != nil {
		t.Errorf("Stored keys was %d, want %d", got, want)
	}
}
