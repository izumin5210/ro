package query

import (
	"github.com/izumin5210/ro/types"
)

const (
	zrange           = "ZRANGE"
	zrevrange        = "ZREVRANGE"
	zrangeByScore    = "ZRANGEBYSCORE"
	zrevrangeByScore = "ZREVRANGEBYSCORE"
	inf              = "+inf"
	neginf           = "-inf"
)

// New creates a query object
func New(key string) types.Query {
	return &Query{
		key:   key,
		limit: -1,
	}
}

// Query builds command and args for operating redis sorted set
type Query struct {
	key     string
	min     interface{}
	max     interface{}
	limit   int
	offset  int
	reverse bool
}

// Build implements the type.Query interface
func (q *Query) Build() (string, []interface{}) {
	return q.command(), q.args()
}

func (q *Query) command() string {
	if q.isWithScore() {
		if q.reverse {
			return zrevrangeByScore
		}
		return zrangeByScore
	}
	if q.reverse {
		return zrevrange
	}
	return zrange
}

func (q *Query) args() []interface{} {
	args := []interface{}{}
	args = append(args, q.key)
	if q.isWithScore() {
		if q.min == nil {
			args = append(args, neginf)
		} else {
			args = append(args, q.min)
		}
		if q.max == nil {
			args = append(args, inf)
		} else {
			args = append(args, q.max)
		}
		if q.offset != 0 || q.limit != -1 {
			args = append(args, "LIMIT", q.offset, q.limit)
		}
	} else {
		end := q.limit
		if end > 0 {
			end = q.offset + q.limit - 1
		}
		args = append(args, q.offset, end)
	}
	return args
}

func (q *Query) isWithScore() bool {
	return q.min != nil || q.max != nil
}
