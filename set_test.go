package ro

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"

	rotesting "github.com/izumin5210/ro/testing"
)

func TestSet(t *testing.T) {
	defer teardown(t)
	now := time.Now().UTC()
	post := &rotesting.Post{
		ID:        1,
		Title:     "post 1",
		Body:      "This is a post 1.",
		UpdatedAt: now.UnixNano(),
	}

	store := New(pool, &rotesting.Post{})
	err := store.Set(post)

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

	err = store.Set(&rotesting.Post{
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

func TestSet_WithMultipleItems(t *testing.T) {
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
	err := store.Set(posts)
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

func TestSet_WhenDisableToStoreToHash(t *testing.T) {
	defer teardown(t)
	now := time.Now().UTC()
	post := &rotesting.Post{
		ID:        1,
		Title:     "post 1",
		Body:      "This is a post 1.",
		UpdatedAt: now.UnixNano(),
	}

	store := New(pool, &rotesting.Post{}, WithHashStore(false))
	err := store.Set(post)

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

func TestSet_WhenKeySuffixIsEmpty(t *testing.T) {
	store := New(pool, &DummyWithEmptyKeySuffix{}, WithHashStore(false))

	dummy := &DummyWithEmptyKeySuffix{}
	err := store.Set(dummy)

	if err == nil {
		t.Error("Set() with an empty key suffix should return an error")
	}

	if got, want := err.Error(), "GetKeySuffix() should be present"; !strings.Contains(got, want) {
		t.Errorf("Set() with an empty key suffix should return an error %q, want to contain %q", got, want)
	}

	conn := pool.Get()
	defer conn.Close()
	keys, _ := redis.Strings(conn.Do("KEYS", "*"))
	if got, want := keys, []string{}; !reflect.DeepEqual(got, want) {
		t.Errorf("Set() with an empty key suffix stores %v, want %v", got, want)
	}
}

type DummyWithNilScoreMap struct {
}

func (d *DummyWithNilScoreMap) GetKeySuffix() string                { return "test" }
func (d *DummyWithNilScoreMap) GetScoreMap() map[string]interface{} { return nil }

func TestSet_WithScoreMapIsNil(t *testing.T) {
	store := New(pool, &DummyWithNilScoreMap{}, WithHashStore(false))
	dummy := &DummyWithNilScoreMap{}
	err := store.Set(dummy)

	if err == nil {
		t.Error("Set() with nil score map should return an error")
	}

	if got, want := err.Error(), "GetScoreMap() should be present"; !strings.Contains(got, want) {
		t.Errorf("Set() with nil score map should return an error %q, want to contain %q", got, want)
	}

	conn := pool.Get()
	defer conn.Close()
	keys, _ := redis.Strings(conn.Do("KEYS", "*"))
	if got, want := keys, []string{}; !reflect.DeepEqual(got, want) {
		t.Errorf("Set() with nil score map stores %v, want %v", got, want)
	}
}

type DummyWithEmptyScoreKey struct {
}

func (d *DummyWithEmptyScoreKey) GetKeySuffix() string { return "test" }
func (d *DummyWithEmptyScoreKey) GetScoreMap() map[string]interface{} {
	return map[string]interface{}{"test": 1, "": 2}
}

func TestSet_WithEmptyScoreKey(t *testing.T) {
	store := New(pool, &DummyWithEmptyScoreKey{}, WithHashStore(false))
	dummy := &DummyWithEmptyScoreKey{}
	err := store.Set(dummy)

	if err == nil {
		t.Error("Set() with empty score key should return an error")
	}

	if got, want := err.Error(), "key in DummyWithEmptyScoreKey:test's GetScoreMap() should be present"; !strings.Contains(got, want) {
		t.Errorf("Set() with empty score key should return an error %q, want to contain %q", got, want)
	}

	conn := pool.Get()
	defer conn.Close()
	keys, _ := redis.Strings(conn.Do("KEYS", "*"))
	if got, want := keys, []string{}; !reflect.DeepEqual(got, want) {
		t.Errorf("Set() with empty score key stores %v, want %v", got, want)
	}
}

type DummyWithNotNumberScore struct {
}

func (d *DummyWithNotNumberScore) GetKeySuffix() string { return "test" }
func (d *DummyWithNotNumberScore) GetScoreMap() map[string]interface{} {
	return map[string]interface{}{"test": 1, "test1": "1.1.1"}
}

func TestSet_WithNotNumberScore(t *testing.T) {
	store := New(pool, &DummyWithNotNumberScore{}, WithHashStore(false))
	dummy := &DummyWithNotNumberScore{}
	err := store.Set(dummy)

	if err == nil {
		t.Error("Set() with not number score should return an error")
	}

	if got, want := err.Error(), "GetScoreMap()[test1] should be number"; !strings.Contains(got, want) {
		t.Errorf("Set() with not number score should return an error %q, want to contain %q", got, want)
	}

	conn := pool.Get()
	defer conn.Close()
	keys, _ := redis.Strings(conn.Do("KEYS", "*"))
	if got, want := keys, []string{}; !reflect.DeepEqual(got, want) {
		t.Errorf("Set() with not number score stores %v, want %v", got, want)
	}
}

type DummyWithTooLargeScore struct {
}

func (d *DummyWithTooLargeScore) GetKeySuffix() string { return "test" }
func (d *DummyWithTooLargeScore) GetScoreMap() map[string]interface{} {
	return map[string]interface{}{"test": 1, "test1": strings.Repeat("2", 309)}
}

func TestSet_WithTooLargeNumberScore(t *testing.T) {
	store := New(pool, &DummyWithTooLargeScore{}, WithHashStore(false))
	dummy := &DummyWithTooLargeScore{}
	err := store.Set(dummy)

	if err == nil {
		t.Error("Set() with not number score should return an error")
	}

	if got, want := err.Error(), "GetScoreMap()[test1] should be number"; !strings.Contains(got, want) {
		t.Errorf("Set() with not number score should return an error %q, want to contain %q", got, want)
	}

	conn := pool.Get()
	defer conn.Close()
	keys, _ := redis.Strings(conn.Do("KEYS", "*"))
	if got, want := keys, []string{}; !reflect.DeepEqual(got, want) {
		t.Errorf("Set() with not number score stores %v, want %v", got, want)
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

func TestSet_WithStringNumberScore(t *testing.T) {
	store := New(pool, &DummyWithStringNumberScore{}, WithHashStore(false))
	dummy := &DummyWithStringNumberScore{}
	err := store.Set(dummy)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	conn := pool.Get()
	defer conn.Close()
	keys, _ := redis.Strings(conn.Do("KEYS", "*"))
	if got, want := len(keys), 5; got != want {
		t.Errorf("Set() with string number score stores %d items, want %d items", got, want)
	}
}
