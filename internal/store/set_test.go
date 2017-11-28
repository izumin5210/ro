package store

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"

	"github.com/izumin5210/ro/internal/config"
	"github.com/izumin5210/ro/types"
)

func TestSet(t *testing.T) {
	defer teardown(t)
	now := time.Now().UTC()
	post := &TestPost{
		ID:        1,
		Title:     "post 1",
		Body:      "This is a post 1.",
		UpdatedAt: now.UnixNano(),
	}

	cnf, _ := config.New()
	cnf.ScorerFuncs = []types.ScorerFunc{
		func(m types.Model) (string, interface{}) {
			return "recent", m.(*TestPost).UpdatedAt
		},
	}

	store, err := New(redisPool.Get, &TestPost{}, cnf)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	err = store.Set(post)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	conn := redisPool.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if got, want := len(keys), 2; err != nil {
		t.Errorf("Stored keys was %d, want %d", got, want)
	}

	v, err := redis.Values(conn.Do("HGETALL", "TestPost:1"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	gotPost := &TestPost{}
	err = redis.ScanStruct(v, gotPost)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if got, want := gotPost, post; !reflect.DeepEqual(got, want) {
		t.Errorf("Stored post is %v, want %v", got, want)
	}

	err = store.Set(&TestPost{
		ID:        2,
		Title:     "post 1",
		Body:      "This is a post 1.",
		UpdatedAt: now.Add(-60 * 60 * 24 * time.Second).UnixNano(),
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	keys, err = redis.Strings(conn.Do("ZRANGE", "TestPost/recent", 0, -1))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if got, want := len(keys), 2; err != nil {
		t.Errorf("Stored keys was %d, want %d", got, want)
	}
	if got, want := keys[0], "TestPost:2"; got != want {
		t.Errorf("Stored key was %q, want %q", got, want)
	}
	if got, want := keys[1], "TestPost:1"; got != want {
		t.Errorf("Stored key was %q, want %q", got, want)
	}
}

func TestSet_WithMultipleItems(t *testing.T) {
	defer teardown(t)
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
			UpdatedAt: now.Add(-60 * 60 * 24 * time.Second).UnixNano(),
		},
		{
			ID:        3,
			Title:     "post 3",
			Body:      "This is a post 3.",
			UpdatedAt: now.Add(60 * 60 * 24 * time.Second).UnixNano(),
		},
	}

	cnf, _ := config.New()
	cnf.ScorerFuncs = []types.ScorerFunc{
		func(m types.Model) (string, interface{}) {
			return "recent", m.(*TestPost).UpdatedAt
		},
	}

	store, err := New(redisPool.Get, &TestPost{}, cnf)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	err = store.Set(posts)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	conn := redisPool.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if got, want := len(keys), 4; err != nil {
		t.Errorf("Stored keys was %d, want %d", got, want)
	}

	for _, post := range posts {
		v, err := redis.Values(conn.Do("HGETALL", fmt.Sprintf("TestPost:%s", post.GetKeySuffix())))
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		gotPost := &TestPost{}
		err = redis.ScanStruct(v, gotPost)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if got, want := gotPost, post; !reflect.DeepEqual(got, want) {
			t.Errorf("Stored post is %v, want %v", got, want)
		}
	}

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	keys, err = redis.Strings(conn.Do("ZRANGE", "TestPost/recent", 0, -1))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if got, want := len(keys), 3; err != nil {
		t.Errorf("Stored keys was %d, want %d", got, want)
	}
	if got, want := keys[0], "TestPost:2"; got != want {
		t.Errorf("Stored key was %q, want %q", got, want)
	}
	if got, want := keys[1], "TestPost:1"; got != want {
		t.Errorf("Stored key was %q, want %q", got, want)
	}
	if got, want := keys[2], "TestPost:3"; got != want {
		t.Errorf("Stored key was %q, want %q", got, want)
	}
}

func TestSet_WhenDisableToStoreToHash(t *testing.T) {
	defer teardown(t)
	now := time.Now().UTC()
	post := &TestPost{
		ID:        1,
		Title:     "post 1",
		Body:      "This is a post 1.",
		UpdatedAt: now.UnixNano(),
	}

	cnf, _ := config.New()
	cnf.ScorerFuncs = []types.ScorerFunc{
		func(m types.Model) (string, interface{}) {
			return "recent", m.(*TestPost).UpdatedAt
		},
	}
	cnf.HashStoreEnabled = false
	store, err := New(redisPool.Get, &TestPost{}, cnf)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	err = store.Set(post)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	conn := redisPool.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if got, want := len(keys), 1; err != nil {
		t.Errorf("Stored keys was %d, want %d", got, want)
	}

	v, err := redis.Values(conn.Do("HGETALL", "TestPost:1"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(v) > 0 {
		t.Errorf("Unexpected response: %v", v)
	}
}
