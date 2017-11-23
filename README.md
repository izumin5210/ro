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
	Title     string `redis:"title"`
	Body      string `redis:"body"`
	UpdatedAt uint64 `redis:"updated_at"`
}

var PostScorerFuncs = []types.ScorerFunc{
	func (m types.Model) (string, interface{}) { return "id", m.(*Post).ID },
	func (m types.Model) (string, interface{}) { return "updated_at", m.(*Post).UpdatedAt },
}

var pool *redis.Pool

func main() {
	pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.DialURL("redis://localhost:6379")
		},
	}

	store := ro.New(pool.Get, &Post{}, ro.WithScorers(PostScorerFuncs))

	// Posts will be stored as Hash, and id and updated_at are stored as OrderedSet
	store.Set([]*Post{
		{
			ID: 1,
			Title: "post 1",
			Body: "This is a post 1",
		},
		{
			ID: 2,
			Title: "post 2",
			Body: "This is a post 2",
		},
		{
			ID: 3,
			Title: "post 3",
			Body: "This is a post 3",
		},
	})

	post := &Post{ID: 1}
	_ := store.Get(post)
	fmt.Println("%v", post)
	// Output:
	// Post{ID: 1, Title: "post 1", Body: "This is a post 1"}

	posts := []*Post{}
	_  := store.Select(&posts, store.Query("id").Gt(2).Lt(10))
	fmt.Println("%v", posts[0])
	// Output:
	// Post{ID: 3, Title: "post 3", Body: "This is a post 3"}
	fmt.Println("%v", posts[1])
	// Output:
	// Post{ID: 2, Title: "post 2", Body: "This is a post 2"}

	cnt, _ := store.Count(store.Query("id").Gt(2))
	fmt.Println(cnt)
	// Output:
	// 2
}
```
