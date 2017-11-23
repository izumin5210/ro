package query

import (
	"testing"

	"github.com/izumin5210/ro/types"
)

func TestQuery_Build(t *testing.T) {
	key := "testkey"
	cases := []struct {
		q    types.Query
		cmd  string
		args []interface{}
	}{
		{
			q:    New(key),
			cmd:  "ZRANGE",
			args: []interface{}{key, 0, -1},
		},
		{
			q:    &Query{key: key},
			cmd:  "ZRANGE",
			args: []interface{}{key, 0, 0},
		},
		{
			q:    New(key).Reverse(),
			cmd:  "ZREVRANGE",
			args: []interface{}{key, 0, -1},
		},
		{
			q:    New(key).Limit(10),
			cmd:  "ZRANGE",
			args: []interface{}{key, 0, 9},
		},
		{
			q:    New(key).Offset(10),
			cmd:  "ZRANGE",
			args: []interface{}{key, 10, -1},
		},
		{
			q:    New(key).Limit(10).Offset(15).Reverse(),
			cmd:  "ZREVRANGE",
			args: []interface{}{key, 15, 24},
		},
		{
			q:    New(key).GtEq(10),
			cmd:  "ZRANGEBYSCORE",
			args: []interface{}{key, 10, "+inf"},
		},
		{
			q:    New(key).LtEq(10),
			cmd:  "ZRANGEBYSCORE",
			args: []interface{}{key, "-inf", 10},
		},
		{
			q:    New(key).Gt(10),
			cmd:  "ZRANGEBYSCORE",
			args: []interface{}{key, "(10", "+inf"},
		},
		{
			q:    New(key).Lt(10),
			cmd:  "ZRANGEBYSCORE",
			args: []interface{}{key, "-inf", "(10"},
		},
		{
			q:    New(key).Eq(10),
			cmd:  "ZRANGEBYSCORE",
			args: []interface{}{key, 10, 10},
		},
		{
			q:    New(key).GtEq(6).LtEq(10),
			cmd:  "ZRANGEBYSCORE",
			args: []interface{}{key, 6, 10},
		},
		{
			q:    New(key).LtEq(10).Limit(15),
			cmd:  "ZRANGEBYSCORE",
			args: []interface{}{key, "-inf", 10, "LIMIT", 0, 15},
		},
		{
			q:    New(key).GtEq(10).Offset(15),
			cmd:  "ZRANGEBYSCORE",
			args: []interface{}{key, 10, "+inf", "LIMIT", 15, -1},
		},
		{
			q:    New(key).LtEq(10).Limit(-1).Reverse(),
			cmd:  "ZREVRANGEBYSCORE",
			args: []interface{}{key, 10, "-inf"},
		},
		{
			q:    New(key).GtEq(10).LtEq(15).Limit(20).Offset(15),
			cmd:  "ZRANGEBYSCORE",
			args: []interface{}{key, 10, 15, "LIMIT", 15, 20},
		},
		{
			q:    New(key).GtEq(10).Lt(15).Limit(20).Offset(15).Reverse(),
			cmd:  "ZREVRANGEBYSCORE",
			args: []interface{}{key, "(15", 10, "LIMIT", 15, 20},
		},
	}

	for _, c := range cases {
		cmd, args := c.q.Build()

		if got, want := cmd, c.cmd; got != want {
			t.Errorf("%v.Bulid() returned %q, want %q", c.q, got, want)
		}

		if got, want := len(args), len(c.args); got != want {
			t.Errorf("%v.Bulid() args returned %d items, want %d items", c.q, got, want)
		} else {
			for i, want := range c.args {
				if got := args[i]; got != want {
					t.Errorf("%v.Bulid() args[%d] returned %v, want %v", c.q, i, got, want)
				}
			}
		}
	}
}
