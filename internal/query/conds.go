package query

import (
	"fmt"

	"github.com/izumin5210/ro/types"
)

// Gt implements the types.Query interface.
func (q *Query) Gt(v interface{}) types.Query {
	q.min = fmt.Sprintf("(%v", v)
	return q
}

// GtEq implements the types.Query interface.
func (q *Query) GtEq(v interface{}) types.Query {
	q.min = v
	return q
}

// Lt implements the types.Query interface.
func (q *Query) Lt(v interface{}) types.Query {
	q.max = fmt.Sprintf("(%v", v)
	return q
}

// LtEq implements the types.Query interface.
func (q *Query) LtEq(v interface{}) types.Query {
	q.max = v
	return q
}

// Limit implements the types.Query interface.
func (q *Query) Limit(v int) types.Query {
	q.limit = v
	return q
}

// Offset implements the types.Query interface.
func (q *Query) Offset(v int) types.Query {
	q.offset = v
	return q
}

// Reverse implements the types.Query interface.
func (q *Query) Reverse() types.Query {
	q.reverse = !q.reverse
	return q
}
