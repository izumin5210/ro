package store

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
	dockertest "gopkg.in/ory-am/dockertest.v3"

	"github.com/izumin5210/ro/types"
)

// Types
// ================================================================

type TestPost struct {
	ID        uint64 `redis:"id"`
	Title     string `redis:"title"`
	Body      string `redis:"body"`
	UpdatedAt int64  `redis:"updated_at"`
}

func (p *TestPost) GetKeyPrefix() string {
	return ""
}

func (p *TestPost) GetKeySuffix() string {
	return fmt.Sprint(p.ID)
}

// Test funcs
// ================================================================

func TestSet(t *testing.T) {
	defer teardown(t)
	now := time.Now().UTC()
	post := &TestPost{
		ID:        1,
		Title:     "post 1",
		Body:      "This is a post 1.",
		UpdatedAt: now.UnixNano(),
	}

	cnf := &types.StoreConfig{
		ScorerFuncMap: map[string]types.ScorerFunc{
			"recent": func(m types.Model) interface{} {
				return m.(*TestPost).UpdatedAt
			},
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

	cnf := &types.StoreConfig{
		ScorerFuncMap: map[string]types.ScorerFunc{
			"recent": func(m types.Model) interface{} {
				return m.(*TestPost).UpdatedAt
			},
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

func TestGet(t *testing.T) {
	defer teardown(t)
	store, err := New(redisPool.Get, &TestPost{}, &types.StoreConfig{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	post := &TestPost{
		ID:    1,
		Title: "post 1",
		Body:  "This is a post 1.",
	}
	key := "TestPost:1"

	conn := redisPool.Get()
	defer conn.Close()
	conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(post)...)

	gotPost := &TestPost{ID: 1}
	err = store.Get(gotPost)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if got, want := gotPost, post; !reflect.DeepEqual(got, want) {
		t.Errorf("Stored post is %v, want %v", got, want)
	}
}

// Setup and Teardown
// ================================================================

var redisPool *redis.Pool

func TestMain(m *testing.M) {
	var err error
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("redis", "4.0.2-alpine", nil)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	redisPool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(fmt.Sprintf("redis://localhost:%s", resource.GetPort("6379/tcp")))
		},
	}

	if err = pool.Retry(func() error {
		conn := redisPool.Get()
		defer conn.Close()
		_, err := conn.Do("PING")

		return err
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	exitCode := m.Run()

	err = redisPool.Close()
	if err != nil {
		log.Fatalf("Failed to close redis pool: %s", err)
	}
	err = pool.Purge(resource)
	if err != nil {
		log.Fatalf("Failed to purge docker pool: %s", err)
	}

	os.Exit(exitCode)
}

func teardown(t *testing.T) {
	conn := redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("FLUSHALL")
	if err != nil {
		log.Fatalf("Failed to flush redis: %s", err)
	}
}
