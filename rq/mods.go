package rq

import (
	"fmt"
)

var (
	// DefaultKeyDelimiter is used to join tokens for key parameters in default.
	DefaultKeyDelimiter = ":"

	// DefaultKeyPrefixDelimiter is used to join key and prefix in default.
	DefaultKeyPrefixDelimiter = "/"
)

// Modifier modifies a Query object.
type Modifier func(*Query)

// Key specifies a key parameter of a query.
func Key(tokens ...interface{}) Modifier {
	return func(q *Query) {
		q.Key.Tokens = tokens
	}
}

// KeyPrefix specifies a key prefix of a query.
func KeyPrefix(prefix string) Modifier {
	return func(q *Query) {
		q.Key.Prefix = prefix
	}
}

// Gt specifies a value range with `>` for scores of a query.
func Gt(v interface{}) Modifier {
	return func(q *Query) {
		q.Min = fmt.Sprintf("(%v", v)
	}
}

// GtEq specifies a value range with `>=` for scores of a query.
func GtEq(v interface{}) Modifier {
	return func(q *Query) {
		q.Min = v
	}
}

// Lt specifies a value range with `<` for scores of a query.
func Lt(v interface{}) Modifier {
	return func(q *Query) {
		q.Max = fmt.Sprintf("(%v", v)
	}
}

// LtEq specifies a value range with `<=` for scores of a query.
func LtEq(v interface{}) Modifier {
	return func(q *Query) {
		q.Max = v
	}
}

// Eq specifies a value range with `=` for scores of a query.
func Eq(v interface{}) Modifier {
	return func(q *Query) {
		LtEq(v)(q)
		GtEq(v)(q)
	}
}

// Limit specifies returned values count of a query.
func Limit(v int) Modifier {
	return func(q *Query) {
		q.Limit = v
	}
}

// Offset specifies returned values offset of a query.
func Offset(v int) Modifier {
	return func(q *Query) {
		q.Offset = v
	}
}

// Reverse specifies returned values order of a query.
func Reverse() Modifier {
	return func(q *Query) {
		q.Reverse = !q.Reverse
	}
}
