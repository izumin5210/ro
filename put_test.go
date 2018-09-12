package ro_test

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/izumin5210/ro"
	rotesting "github.com/izumin5210/ro/testing"
)

func TestRedisStore_Put(t *testing.T) {
	defer teardown(t)
	now := time.Now().UTC()
	post := &rotesting.Post{
		ID:        1,
		Title:     "post 1",
		Body:      "This is a post 1.",
		UpdatedAt: now.UnixNano(),
	}

	store := ro.New(pool, &rotesting.Post{})
	err := store.Put(context.TODO(), post)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	conn := pool.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if got, want := len(keys), 2; err != nil {
		t.Errorf("Stored keys was %d, want %d", got, want)
	}

	v, err := redis.Values(conn.Do("HGETALL", "Post:1"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	gotPost := &rotesting.Post{}
	err = redis.ScanStruct(v, gotPost)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if got, want := gotPost, post; !reflect.DeepEqual(got, want) {
		t.Errorf("Stored post is %v, want %v", got, want)
	}

	err = store.Put(context.TODO(), &rotesting.Post{
		ID:        2,
		Title:     "post 1",
		Body:      "This is a post 1.",
		UpdatedAt: now.Add(-60 * 60 * 24 * time.Second).UnixNano(),
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	keys, err = redis.Strings(conn.Do("ZRANGE", "Post/recent", 0, -1))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if got, want := len(keys), 2; err != nil {
		t.Errorf("Stored keys was %d, want %d", got, want)
	}
	if got, want := keys[0], "Post:2"; got != want {
		t.Errorf("Stored key was %q, want %q", got, want)
	}
	if got, want := keys[1], "Post:1"; got != want {
		t.Errorf("Stored key was %q, want %q", got, want)
	}
}

func TestRedisStore_Put_WithMultipleItems(t *testing.T) {
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

	store := ro.New(pool, &rotesting.Post{})
	err := store.Put(context.TODO(), posts)
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

	for _, post := range posts {
		v, err := redis.Values(conn.Do("HGETALL", fmt.Sprintf("Post:%s", post.GetKeySuffix())))
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		gotPost := &rotesting.Post{}
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
	keys, err = redis.Strings(conn.Do("ZRANGE", "Post/recent", 0, -1))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if got, want := len(keys), 3; err != nil {
		t.Errorf("Stored keys was %d, want %d", got, want)
	}
	if got, want := keys[0], "Post:2"; got != want {
		t.Errorf("Stored key was %q, want %q", got, want)
	}
	if got, want := keys[1], "Post:1"; got != want {
		t.Errorf("Stored key was %q, want %q", got, want)
	}
	if got, want := keys[2], "Post:3"; got != want {
		t.Errorf("Stored key was %q, want %q", got, want)
	}
}

func TestRedisStore_Put_WhenDisableToStoreToHash(t *testing.T) {
	defer teardown(t)
	now := time.Now().UTC()
	post := &rotesting.Post{
		ID:        1,
		Title:     "post 1",
		Body:      "This is a post 1.",
		UpdatedAt: now.UnixNano(),
	}

	store := ro.New(pool, &rotesting.Post{}, ro.WithHashStore(false))
	err := store.Put(context.TODO(), post)

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

	v, err := redis.Values(conn.Do("HGETALL", "Post:1"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(v) > 0 {
		t.Errorf("Unexpected response: %v", v)
	}
}

type DummyWithEmptyKeySuffix struct {
}

func (d *DummyWithEmptyKeySuffix) GetKeySuffix() string { return "" }
func (d *DummyWithEmptyKeySuffix) GetScoreMap() map[string]interface{} {
	return map[string]interface{}{}
}

func TestRedisStore_Put_WhenKeySuffixIsEmpty(t *testing.T) {
	store := ro.New(pool, &DummyWithEmptyKeySuffix{}, ro.WithHashStore(false))

	dummy := &DummyWithEmptyKeySuffix{}
	err := store.Put(context.TODO(), dummy)

	if err == nil {
		t.Error("Put() with an empty key suffix should return an error")
	}

	if got, want := err.Error(), "GetKeySuffix() should be present"; !strings.Contains(got, want) {
		t.Errorf("Put() with an empty key suffix should return an error %q, want to contain %q", got, want)
	}

	conn := pool.Get()
	defer conn.Close()
	keys, _ := redis.Strings(conn.Do("KEYS", "*"))
	if got, want := keys, []string{}; !reflect.DeepEqual(got, want) {
		t.Errorf("Put() with an empty key suffix stores %v, want %v", got, want)
	}
}

type DummyWithNilScoreMap struct {
}

func (d *DummyWithNilScoreMap) GetKeySuffix() string                { return "test" }
func (d *DummyWithNilScoreMap) GetScoreMap() map[string]interface{} { return nil }

func TestRedisStore_Put_WithScoreMapIsNil(t *testing.T) {
	store := ro.New(pool, &DummyWithNilScoreMap{}, ro.WithHashStore(false))
	dummy := &DummyWithNilScoreMap{}
	err := store.Put(context.TODO(), dummy)

	if err == nil {
		t.Error("Put() with nil score map should return an error")
	}

	if got, want := err.Error(), "GetScoreMap() should be present"; !strings.Contains(got, want) {
		t.Errorf("Put() with nil score map should return an error %q, want to contain %q", got, want)
	}

	conn := pool.Get()
	defer conn.Close()
	keys, _ := redis.Strings(conn.Do("KEYS", "*"))
	if got, want := keys, []string{}; !reflect.DeepEqual(got, want) {
		t.Errorf("Put() with nil score map stores %v, want %v", got, want)
	}
}

type DummyWithEmptyScoreKey struct {
}

func (d *DummyWithEmptyScoreKey) GetKeySuffix() string { return "test" }
func (d *DummyWithEmptyScoreKey) GetScoreMap() map[string]interface{} {
	return map[string]interface{}{"test": 1, "": 2}
}

func TestRedisStore_Put_WithEmptyScoreKey(t *testing.T) {
	store := ro.New(pool, &DummyWithEmptyScoreKey{}, ro.WithHashStore(false))
	dummy := &DummyWithEmptyScoreKey{}
	err := store.Put(context.TODO(), dummy)

	if err == nil {
		t.Error("Put() with empty score key should return an error")
	}

	if got, want := err.Error(), "key in DummyWithEmptyScoreKey:test's GetScoreMap() should be present"; !strings.Contains(got, want) {
		t.Errorf("Put() with empty score key should return an error %q, want to contain %q", got, want)
	}

	conn := pool.Get()
	defer conn.Close()
	keys, _ := redis.Strings(conn.Do("KEYS", "*"))
	if got, want := keys, []string{}; !reflect.DeepEqual(got, want) {
		t.Errorf("Put() with empty score key stores %v, want %v", got, want)
	}
}

type DummyWithNotNumberScore struct {
}

func (d *DummyWithNotNumberScore) GetKeySuffix() string { return "test" }
func (d *DummyWithNotNumberScore) GetScoreMap() map[string]interface{} {
	return map[string]interface{}{"test": 1, "test1": "1.1.1"}
}

func TestRedisStore_Put_WithNotNumberScore(t *testing.T) {
	store := ro.New(pool, &DummyWithNotNumberScore{}, ro.WithHashStore(false))
	dummy := &DummyWithNotNumberScore{}
	err := store.Put(context.TODO(), dummy)

	if err == nil {
		t.Error("Put() with not number score should return an error")
	}

	if got, want := err.Error(), "GetScoreMap()[test1] should be number"; !strings.Contains(got, want) {
		t.Errorf("Put() with not number score should return an error %q, want to contain %q", got, want)
	}

	conn := pool.Get()
	defer conn.Close()
	keys, _ := redis.Strings(conn.Do("KEYS", "*"))
	if got, want := keys, []string{}; !reflect.DeepEqual(got, want) {
		t.Errorf("Put() with not number score stores %v, want %v", got, want)
	}
}

type DummyWithTooLargeScore struct {
}

func (d *DummyWithTooLargeScore) GetKeySuffix() string { return "test" }
func (d *DummyWithTooLargeScore) GetScoreMap() map[string]interface{} {
	return map[string]interface{}{"test": 1, "test1": strings.Repeat("2", 309)}
}

func TestRedisStore_Put_WithTooLargeNumberScore(t *testing.T) {
	store := ro.New(pool, &DummyWithTooLargeScore{}, ro.WithHashStore(false))
	dummy := &DummyWithTooLargeScore{}
	err := store.Put(context.TODO(), dummy)

	if err == nil {
		t.Error("Put() with not number score should return an error")
	}

	if got, want := err.Error(), "GetScoreMap()[test1] should be number"; !strings.Contains(got, want) {
		t.Errorf("Put() with not number score should return an error %q, want to contain %q", got, want)
	}

	conn := pool.Get()
	defer conn.Close()
	keys, _ := redis.Strings(conn.Do("KEYS", "*"))
	if got, want := keys, []string{}; !reflect.DeepEqual(got, want) {
		t.Errorf("Put() with not number score stores %v, want %v", got, want)
	}
}

type DummyWithStringNumberScore struct {
}

func (d *DummyWithStringNumberScore) GetKeySuffix() string { return "test" }
func (d *DummyWithStringNumberScore) GetScoreMap() map[string]interface{} {
	return map[string]interface{}{
		"test":  1,
		"test1": "100.1",
		"test2": "100",
		"test3": strings.Repeat("1", 309),
	}
}

func TestRedisStore_Put_WithStringNumberScore(t *testing.T) {
	store := ro.New(pool, &DummyWithStringNumberScore{}, ro.WithHashStore(false))
	dummy := &DummyWithStringNumberScore{}
	err := store.Put(context.TODO(), dummy)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	conn := pool.Get()
	defer conn.Close()
	keys, _ := redis.Strings(conn.Do("KEYS", "*"))
	if got, want := len(keys), 5; got != want {
		t.Errorf("Put() with string number score stores %d items, want %d items", got, want)
	}
}
