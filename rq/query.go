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
func (q *QueryKey) Build() (string, error) {
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
	if key == "" {
		return "", errors.New("key is required")
	}
	return key, nil
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
	switch q.Type {
	case CommandList:
		cmd, err := q.buildListCommand()
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return cmd, nil
	case CommandCount:
		cmd, err := q.buildCountCommand()
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return cmd, nil
	default:
		return nil, errors.WithStack(newQueryError(q, "unknown query type"))
	}
}

func (q *Query) buildListCommand() (*Command, error) {
	key, err := q.Key.Build()
	if err != nil {
		return nil, errors.WithStack(newQueryError(q, err.Error()))
	}

	cmd := &Command{Args: make([]interface{}, 1, 10)}
	cmd.Args[0] = key

	if q.isWithScore() {
		min, max := q.getMinAndMax()

		if q.Reverse {
			cmd.Name = zrevrangeByScore
			cmd.Args = append(cmd.Args, max, min)
		} else {
			cmd.Name = zrangeByScore
			cmd.Args = append(cmd.Args, min, max)
		}

		if q.Offset != 0 || q.Limit != -1 {
			cmd.Args = append(cmd.Args, "LIMIT", q.Offset, q.Limit)
		}
	} else {
		if q.Reverse {
			cmd.Name = zrevrange
		} else {
			cmd.Name = zrange
		}

		end := q.Limit
		if end > 0 {
			end = q.Offset + q.Limit - 1
		}
		cmd.Args = append(cmd.Args, q.Offset, end)
	}

	return cmd, nil
}

func (q *Query) buildCountCommand() (*Command, error) {
	key, err := q.Key.Build()
	if err != nil {
		return nil, errors.WithStack(newQueryError(q, err.Error()))
	}

	cmd := &Command{Name: zcard, Args: make([]interface{}, 1, 10)}
	cmd.Args[0] = key

	if q.isWithScore() {
		cmd.Name = zcount
		min, max := q.getMinAndMax()
		cmd.Args = append(cmd.Args, min, max)
	}

	return cmd, nil
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
