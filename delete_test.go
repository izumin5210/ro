package ro

import (
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"

	rotesting "github.com/izumin5210/ro/testing"
)

func TestRedisStore_Delete(t *testing.T) {
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
			UpdatedAt: now.Add(-60 * 60 * 24 * time.Second).UnixNano(),
		},
		{
			ID:        3,
			Title:     "post 3",
			Body:      "This is a post 3.",
			UpdatedAt: now.Add(60 * 60 * 24 * time.Second).UnixNano(),
		},
	}

	store := New(pool, &rotesting.Post{})
	err := store.Put(posts)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	err = store.Delete(posts[0])
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	conn := pool.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if got, want := len(keys), 3; err != nil {
		t.Errorf("Stored keys was %d, want %d", got, want)
	}

	v, err := redis.Values(conn.Do("HGETALL", "rotesting.Post:1"))
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

func TestRedisStore_Delete_WithMultipleItems(t *testing.T) {
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
			UpdatedAt: now.Add(-60 * 60 * 24 * time.Second).UnixNano(),
		},
		{
			ID:        3,
			Title:     "post 3",
			Body:      "This is a post 3.",
			UpdatedAt: now.Add(60 * 60 * 24 * time.Second).UnixNano(),
		},
	}

	store := New(pool, &rotesting.Post{})
	err := store.Put(posts)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	err = store.Delete([]*rotesting.Post{posts[0], posts[2]})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	conn := pool.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if got, want := len(keys), 1; err != nil {
		t.Errorf("Stored keys was %d, want %d", got, want)
	}

	v, err := redis.Values(conn.Do("HGETALL", "rotesting.Post:1"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(v) > 0 {
		t.Errorf("Unexpected response: %v", v)
	}

	v, err = redis.Values(conn.Do("HGETALL", "rotesting.Post:3"))
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
	if got, want := len(keys), 1; err != nil {
		t.Errorf("Stored keys was %d, want %d", got, want)
	}
}

func TestRedisStore_Remove_WhenDisableToStoreToHash(t *testing.T) {
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
			UpdatedAt: now.Add(-60 * 60 * 24 * time.Second).UnixNano(),
		},
		{
			ID:        3,
			Title:     "post 3",
			Body:      "This is a post 3.",
			UpdatedAt: now.Add(60 * 60 * 24 * time.Second).UnixNano(),
		},
	}

	store := New(pool, &rotesting.Post{}, WithHashStore(false))
	err := store.Put(posts)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	err = store.Delete(posts[0])
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	conn := pool.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if got, want := len(keys), 3; err != nil {
		t.Errorf("Stored keys was %d, want %d", got, want)
	}

	keys, err = redis.Strings(conn.Do("ZRANGE", "TestPoset/recent", 0, -1))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if got, want := len(keys), 2; err != nil {
		t.Errorf("Stored keys was %d, want %d", got, want)
	}
}
