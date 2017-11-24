# ro: Redis Objects to Go
[![Build Status](https://travis-ci.org/izumin5210/ro.svg?branch=master)](https://travis-ci.org/izumin5210/ro)
[![codecov](https://codecov.io/gh/izumin5210/ro/branch/master/graph/badge.svg)](https://codecov.io/gh/izumin5210/ro)
[![GoDoc](https://godoc.org/github.com/izumin5210/ro?status.svg)](https://godoc.org/github.com/izumin5210/ro)
[![Go Report Card](https://goreportcard.com/badge/github.com/izumin5210/ro)](https://goreportcard.com/report/github.com/izumin5210/ro)
[![Go project version](https://badge.fury.io/go/github.com%2Fizumin5210%2Fro.svg)](https://badge.fury.io/go/github.com%2Fizumin5210%2Fro)
[![license](https://img.shields.io/github/license/izumin5210/ro.svg)](./LICENSE)

## Example

```go
type Post struct {
	ro.Model
	ID        uint64 `redis:"id"`
	UserID    int `redis:"user_id"`
	Title     string `redis:"title"`
	Body      string `redis:"body"`
	CreatedAt uint64 `redis:"created_at"`
}

func (p *Post) GetKeySuffix() string {
	return fmt.Sprint(p.ID)
}

var PostScorerFuncs = []types.ScorerFunc{
	func (m types.Model) (string, interface{}) { return "created_at", m.(*Post).CreatedAt },
	func (m types.Model) (string, interface{}) {
		p := m.(*Post)
		return fmt.Sprintf("user:%d", p.UserID), p.CreatedAt
	},
}

var pool *redis.Pool

func main() {
	pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.DialURL("redis://localhost:6379")
		},
	}

	store := ro.New(pool.Get, &Post{}, ro.WithScorers(PostScorerFuncs))
	now := time.Now()

	// Posts will be stored as Hash, and user:{{userID}} and created_at are stored as OrderedSet
	store.Set([]*Post{
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
	})

	post := &Post{ID: 1}
	_ := store.Get(post)
	fmt.Println("%v", post)
	// Output:
	// Post{ID: 1, Title: "post 1", Body: "This is a post 1"}

	posts := []*Post{}
	_ := store.Select(&posts, store.Query("created_at").GtEq(now.UnixNano()).Reverse())
	fmt.Println("%v", posts[0])
	// Output:
	// Post{ID: 3, UserID: 1, Title: "post 3", Body: "This is a post 3"}
	fmt.Println("%v", posts[1])
	// Output:
	// Post{ID: 1, UserID: 1, Title: "post 1", Body: "This is a post 1"}

	cnt, _ := store.Count(store.Query("user:1").Gt(now.UnixNano()).Reverse())
	fmt.Println(cnt)
	// Output:
	// 1
}
```
