package query

import (
	"testing"

	"github.com/izumin5210/ro/types"
)

func TestQuery_BuildForCount(t *testing.T) {
	key := "testkey"
	cases := []struct {
		q    types.Query
		cmd  string
		args []interface{}
	}{
		{
			q:    New(key),
			cmd:  "ZCARD",
			args: []interface{}{key},
		},
		{
			q:    New(key).GtEq(10),
			cmd:  "ZCOUNT",
			args: []interface{}{key, 10, "+inf"},
		},
		{
			q:    New(key).LtEq(10),
			cmd:  "ZCOUNT",
			args: []interface{}{key, "-inf", 10},
		},
		{
			q:    New(key).Gt(10),
			cmd:  "ZCOUNT",
			args: []interface{}{key, "(10", "+inf"},
		},
		{
			q:    New(key).Lt(10),
			cmd:  "ZCOUNT",
			args: []interface{}{key, "-inf", "(10"},
		},
		{
			q:    New(key).Eq(10),
			cmd:  "ZCOUNT",
			args: []interface{}{key, 10, 10},
		},
		{
			q:    New(key).GtEq(6).LtEq(10),
			cmd:  "ZCOUNT",
			args: []interface{}{key, 6, 10},
		},
	}

	for _, c := range cases {
		cmd, args := c.q.BuildForCount()

		if got, want := cmd, c.cmd; got != want {
			t.Errorf("%v.BulidForCount() returned %q, want %q", c.q, got, want)
		}

		if got, want := len(args), len(c.args); got != want {
			t.Errorf("%v.BulidForCount() args returned %d items, want %d items", c.q, got, want)
		} else {
			for i, want := range c.args {
				if got := args[i]; got != want {
					t.Errorf("%v.BulidForCount() args[%d] returned %v, want %v", c.q, i, got, want)
				}
			}
		}
	}
}
