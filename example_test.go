package ro

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"

	"github.com/izumin5210/ro/internal/testing"
)

type Post struct {
	Model
	ID        uint64 `redis:"id"`
	UserID    uint64 `redis:"user_id"`
	Title     string `redis:"title"`
	Body      string `redis:"body"`
	CreatedAt int64  `redis:"created_at"`
}

func (p *Post) GetKeySuffix() string {
	return fmt.Sprint(p.ID)
}

func (p *Post) GetScoreMap() map[string]interface{} {
	return map[string]interface{}{
		"recent": p.CreatedAt,
		fmt.Sprintf("user:%d", p.UserID): p.CreatedAt,
	}
}

// Setup and cleanup
// ----------------------------------------------------------------

var (
	pool      *rotesting.Pool
	postStore Store
	now       time.Time
)

func TestMain(m *testing.M) {
	pool = rotesting.MustCreate()
	defer pool.MustClose()
	postStore, _ = New(pool.Get, &Post{})

	os.Exit(m.Run())
}

func setup() {
	now = time.Now()

	postStore.Set([]*Post{
		{
			ID:        1,
			UserID:    1,
			Title:     "post 1",
			Body:      "This is a post 1",
			CreatedAt: now.UnixNano(),
		},
		{
			ID:        2,
			UserID:    2,
			Title:     "post 2",
			Body:      "This is a post 2",
			CreatedAt: now.Add(-24 * 60 * 60 * time.Second).UnixNano(),
		},
		{
			ID:        3,
			UserID:    1,
			Title:     "post 3",
			Body:      "This is a post 3",
			CreatedAt: now.Add(24 * 60 * 60 * time.Second).UnixNano(),
		},
		{
			ID:        4,
			UserID:    1,
			Title:     "post 4",
			Body:      "This is a post 4",
			CreatedAt: now.Add(-24 * 60 * 60 * time.Second).UnixNano(),
		},
	})
}

func cleanup() {
	pool.Cleanup()
}

// Examples
// ----------------------------------------------------------------

func Example_set() {
	defer cleanup()

	postStore.Set([]*Post{
		{
			ID:        1,
			UserID:    1,
			Title:     "post 1",
			Body:      "This is a post 1",
			CreatedAt: now.UnixNano(),
		},
		{
			ID:        2,
			UserID:    1,
			Title:     "post 2",
			Body:      "This is a post 2",
			CreatedAt: now.Add(-24 * 60 * 60 * time.Second).UnixNano(),
		},
	})

	conn := pool.Get()
	defer conn.Close()

	keys, _ := redis.Strings(conn.Do("ZRANGE", "Post/user:1", 0, -1))
	fmt.Println(keys)
	for _, k := range keys {
		post := &Post{}
		v, _ := redis.Values(conn.Do("HGETALL", k))
		redis.ScanStruct(v, post)
		fmt.Println(post.Body)
	}

	// Output:
	// [Post:2 Post:1]
	// This is a post 2
	// This is a post 1
}

func Example_get() {
	setup()
	defer cleanup()

	post := &Post{ID: 1}
	postStore.Get(post)
	fmt.Println(post.Body)
	// Output:
	// This is a post 1
}

func Example_select() {
	setup()
	defer cleanup()

	posts := []*Post{}
	postStore.Select(&posts, postStore.Query("recent").GtEq(now.UnixNano()).Reverse())
	fmt.Println(posts[0].Body)
	fmt.Println(posts[1].Body)
	// Output:
	// This is a post 3
	// This is a post 1
}

func Example_count() {
	setup()
	defer cleanup()

	cnt, _ := postStore.Count(postStore.Query("user:1").LtEq(now.UnixNano()))
	fmt.Println(cnt)
	// Output:
	// 2
}
