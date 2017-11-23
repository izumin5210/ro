package store

import (
	"reflect"
	"testing"

	"github.com/garyburd/redigo/redis"

	"github.com/izumin5210/ro/types"
)

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
