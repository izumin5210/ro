# ro: map Redis Objects to Go
[![GoDoc](https://godoc.org/github.com/izumin5210/ro?status.svg)](https://godoc.org/github.com/izumin5210/ro)

## Example

```go
type Post struct {
	ro.Base
	ID        uint64 `redis:"id"`
	Title     string `redis:"title"`
	Body      string `redis:"body"`
	UpdatedAt uint64 `redis:"updated_at"`
}

var PostScorerMap = map[string] interface{} {
	"id": func (p *Post) interface{} { return p.ID },
	"updated": func (p *Post) interface{} { return p.UpdatedAt },
}

var pool *redis.Pool

func main() {
	pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.DialURL("redis://localhost:6379")
		},
	}

	store := ro.NewStore(pool, &Post{}, ro.WithScorer(PostScorerMap))

	// Posts will be stored as Hash, and id and updated_at are stored as OrderedSet
	store.Set(&Post{
		ID: 1,
		Title: "post 1",
		Body: "This is a post 1",
	})
	store.Set(&Post{
		ID: 2,
		Title: "post 2",
		Body: "This is a post 2",
	})
	store.Set(&Post{
		ID: 3,
		Title: "post 3",
		Body: "This is a post 3",
	})

	post := &Post{ID: 1}
	_ := store.Get(post)
	fmt.Println("%v", post)
	// Output:
	// Post{ID: 1, Title: "post 1", Body: "This is a post 1"}

	posts := []*Post{}
	_  := store.Select(&posts, ro.Query("id").Gt(2).Lt(10))
	fmt.Println("%v", posts[0])
	// Output:
	// Post{ID: 3, Title: "post 3", Body: "This is a post 3"}
	fmt.Println("%v", posts[1])
	// Output:
	// Post{ID: 2, Title: "post 2", Body: "This is a post 2"}

	cnt, _ := store.Count(ro.Query("id").Gt(2))
	fmt.Println(cnt)
	// Output:
	// 2
}
```
