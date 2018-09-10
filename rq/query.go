package rq

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// CommandType represents operation types of a query.
type CommandType int

// Command types
const (
	CommandList CommandType = iota
	CommandCount
)

const (
	zrange           = "ZRANGE"
	zrevrange        = "ZREVRANGE"
	zrangeByScore    = "ZRANGEBYSCORE"
	zrevrangeByScore = "ZREVRANGEBYSCORE"
	zcard            = "ZCARD"
	zcount           = "ZCOUNT"
	inf              = "+inf"
	neginf           = "-inf"
)

// List creates a Query object for list operation.
func List(mods ...Modifier) *Query {
	return createQuery(CommandList, mods)
}

// Count creates a Query object for count operation.
func Count(mods ...Modifier) *Query {
	return createQuery(CommandCount, mods)
}

func createQuery(t CommandType, mods []Modifier) *Query {
	q := &Query{Type: t}
	switch t {
	case CommandList:
		q.Limit = -1
	}
	for _, f := range mods {
		f(q)
	}
	return q
}

// QueryKey contains parameters to build a key for a redis command.
type QueryKey struct {
	Delimiter       string
	Tokens          []interface{}
	Prefix          string
	PrefixDelimiter string
}

// Build creates a key string.
func (q *QueryKey) Build() string {
	delim := q.Delimiter
	if delim == "" {
		delim = DefaultKeyDelimiter
	}
	strs := make([]string, len(q.Tokens))
	for i, t := range q.Tokens {
		strs[i] = fmt.Sprint(t)
	}
	key := strings.Join(strs, delim)
	if q.Prefix != "" {
		prefixDelim := q.PrefixDelimiter
		if prefixDelim == "" {
			prefixDelim = DefaultKeyPrefixDelimiter
		}
		key = q.Prefix + prefixDelim + key
	}
	return key
}

// Query contains parameters to build a redis command.
type Query struct {
	Type    CommandType
	Key     QueryKey
	Min     interface{}
	Max     interface{}
	Limit   int
	Offset  int
	Reverse bool
}

// Build decide a redis command and args from query parameters.
func (q *Query) Build() (*Command, error) {
	name, err := q.decideCommandName()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	args, err := q.buildArgs()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &Command{Name: name, Args: args}, nil
}

func (q *Query) decideCommandName() (string, error) {
	switch q.Type {
	case CommandList:
		if q.isWithScore() {
			if q.Reverse {
				return zrevrangeByScore, nil
			}
			return zrangeByScore, nil
		}
		if q.Reverse {
			return zrevrange, nil
		}
		return zrange, nil
	case CommandCount:
		if q.isWithScore() {
			return zcount, nil
		}
		return zcard, nil
	default:
		return "", fmt.Errorf("unknown query type: %v", q.Type)
	}
}

func (q *Query) buildArgs() ([]interface{}, error) {
	args := make([]interface{}, 0, 10)
	switch q.Type {
	case CommandList:
		args = append(args, q.Key.Build())
		if q.isWithScore() {
			min, max := q.getMinAndMax()
			if q.Reverse {
				args = append(args, max, min)
			} else {
				args = append(args, min, max)
			}
			if q.Offset != 0 || q.Limit != -1 {
				args = append(args, "LIMIT", q.Offset, q.Limit)
			}
		} else {
			end := q.Limit
			if end > 0 {
				end = q.Offset + q.Limit - 1
			}
			args = append(args, q.Offset, end)
		}
	case CommandCount:
		args = append(args, q.Key.Build())
		if q.isWithScore() {
			min, max := q.getMinAndMax()
			args = append(args, min, max)
		}
	default:
		return nil, fmt.Errorf("unknown query type: %v", q.Type)
	}
	return args, nil
}

func (q *Query) isWithScore() bool {
	return q.Min != nil || q.Max != nil
}

func (q *Query) getMinAndMax() (interface{}, interface{}) {
	max, min := q.Max, q.Min
	if max == nil {
		max = inf
	}
	if min == nil {
		min = neginf
	}
	return min, max
}

// Command contains a redis command name and args.
type Command struct {
	Name string
	Args []interface{}
}

func (c *Command) String() string {
	buf := new(bytes.Buffer)
	buf.WriteString(c.Name)
	for _, a := range c.Args {
		buf.WriteString(" " + fmt.Sprint(a))
	}
	return buf.String()
}
