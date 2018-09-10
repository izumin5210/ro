package rq_test

import (
	"reflect"
	"testing"

	"github.com/izumin5210/ro/rq"
)

func TestQuery_Build(t *testing.T) {
	cases := []struct {
		build func(...rq.Modifier) *rq.Query
		mods  []rq.Modifier
		cmd   *rq.Command
		isErr bool
	}{
		{
			build: rq.List,
			mods:  []rq.Modifier{rq.Key("foo")},
			cmd:   &rq.Command{Name: "ZRANGE", Args: []interface{}{"foo", 0, -1}},
		},
		{
			build: rq.List,
			mods:  []rq.Modifier{rq.Key("foo", 1, "bar", 2)},
			cmd:   &rq.Command{Name: "ZRANGE", Args: []interface{}{"foo:1:bar:2", 0, -1}},
		},
		{
			build: rq.List,
			mods:  []rq.Modifier{rq.Key("foo", 1, "bar", 2), rq.KeyPrefix("baz")},
			cmd:   &rq.Command{Name: "ZRANGE", Args: []interface{}{"baz/foo:1:bar:2", 0, -1}},
		},
		{
			build: rq.List,
			mods:  []rq.Modifier{rq.Key("foo"), rq.Reverse()},
			cmd:   &rq.Command{Name: "ZREVRANGE", Args: []interface{}{"foo", 0, -1}},
		},
		{
			build: rq.List,
			mods:  []rq.Modifier{rq.Key("foo"), rq.Limit(10)},
			cmd:   &rq.Command{Name: "ZRANGE", Args: []interface{}{"foo", 0, 9}},
		},
		{
			build: rq.List,
			mods:  []rq.Modifier{rq.Key("foo"), rq.Offset(10)},
			cmd:   &rq.Command{Name: "ZRANGE", Args: []interface{}{"foo", 10, -1}},
		},
		{
			build: rq.List,
			mods:  []rq.Modifier{rq.Key("foo"), rq.Limit(10), rq.Offset(15), rq.Reverse()},
			cmd:   &rq.Command{Name: "ZREVRANGE", Args: []interface{}{"foo", 15, 24}},
		},
		{
			build: rq.List,
			mods:  []rq.Modifier{rq.Key("foo"), rq.GtEq(10)},
			cmd:   &rq.Command{Name: "ZRANGEBYSCORE", Args: []interface{}{"foo", 10, "+inf"}},
		},
		{
			build: rq.List,
			mods:  []rq.Modifier{rq.Key("foo"), rq.LtEq(10)},
			cmd:   &rq.Command{Name: "ZRANGEBYSCORE", Args: []interface{}{"foo", "-inf", 10}},
		},
		{
			build: rq.List,
			mods:  []rq.Modifier{rq.Key("foo"), rq.Gt(10)},
			cmd:   &rq.Command{Name: "ZRANGEBYSCORE", Args: []interface{}{"foo", "(10", "+inf"}},
		},
		{
			build: rq.List,
			mods:  []rq.Modifier{rq.Key("foo"), rq.Lt(10)},
			cmd:   &rq.Command{Name: "ZRANGEBYSCORE", Args: []interface{}{"foo", "-inf", "(10"}},
		},
		{
			build: rq.List,
			mods:  []rq.Modifier{rq.Key("foo"), rq.Eq(10)},
			cmd:   &rq.Command{Name: "ZRANGEBYSCORE", Args: []interface{}{"foo", 10, 10}},
		},
		{
			build: rq.List,
			mods:  []rq.Modifier{rq.Key("foo"), rq.GtEq(6), rq.LtEq(10)},
			cmd:   &rq.Command{Name: "ZRANGEBYSCORE", Args: []interface{}{"foo", 6, 10}},
		},
		{
			build: rq.List,
			mods:  []rq.Modifier{rq.Key("foo"), rq.LtEq(10), rq.Limit(15)},
			cmd:   &rq.Command{Name: "ZRANGEBYSCORE", Args: []interface{}{"foo", "-inf", 10, "LIMIT", 0, 15}},
		},
		{
			build: rq.List,
			mods:  []rq.Modifier{rq.Key("foo"), rq.GtEq(10), rq.Offset(15)},
			cmd:   &rq.Command{Name: "ZRANGEBYSCORE", Args: []interface{}{"foo", 10, "+inf", "LIMIT", 15, -1}},
		},
		{
			build: rq.List,
			mods:  []rq.Modifier{rq.Key("foo"), rq.LtEq(10), rq.Limit(-1), rq.Reverse()},
			cmd:   &rq.Command{Name: "ZREVRANGEBYSCORE", Args: []interface{}{"foo", 10, "-inf"}},
		},
		{
			build: rq.List,
			mods:  []rq.Modifier{rq.Key("foo"), rq.GtEq(10), rq.LtEq(15), rq.Limit(20), rq.Offset(15)},
			cmd:   &rq.Command{Name: "ZRANGEBYSCORE", Args: []interface{}{"foo", 10, 15, "LIMIT", 15, 20}},
		},
		{
			build: rq.List,
			mods:  []rq.Modifier{rq.Key("foo"), rq.GtEq(10), rq.Lt(15), rq.Limit(20), rq.Offset(15), rq.Reverse()},
			cmd:   &rq.Command{Name: "ZREVRANGEBYSCORE", Args: []interface{}{"foo", "(15", 10, "LIMIT", 15, 20}},
		},
		{
			build: rq.Count,
			mods:  []rq.Modifier{rq.Key("foo")},
			cmd:   &rq.Command{Name: "ZCARD", Args: []interface{}{"foo"}},
		},
		{
			build: rq.Count,
			mods:  []rq.Modifier{rq.Key("foo", 1, "bar", 2)},
			cmd:   &rq.Command{Name: "ZCARD", Args: []interface{}{"foo:1:bar:2"}},
		},
		{
			build: rq.Count,
			mods:  []rq.Modifier{rq.Key("foo", 1, "bar", 2), rq.KeyPrefix("baz")},
			cmd:   &rq.Command{Name: "ZCARD", Args: []interface{}{"baz/foo:1:bar:2"}},
		},
		{
			build: rq.Count,
			mods:  []rq.Modifier{rq.Key("foo"), rq.GtEq(10)},
			cmd:   &rq.Command{Name: "ZCOUNT", Args: []interface{}{"foo", 10, "+inf"}},
		},
		{
			build: rq.Count,
			mods:  []rq.Modifier{rq.Key("foo"), rq.LtEq(10)},
			cmd:   &rq.Command{Name: "ZCOUNT", Args: []interface{}{"foo", "-inf", 10}},
		},
		{
			build: rq.Count,
			mods:  []rq.Modifier{rq.Key("foo"), rq.Gt(10)},
			cmd:   &rq.Command{Name: "ZCOUNT", Args: []interface{}{"foo", "(10", "+inf"}},
		},
		{
			build: rq.Count,
			mods:  []rq.Modifier{rq.Key("foo"), rq.Lt(10)},
			cmd:   &rq.Command{Name: "ZCOUNT", Args: []interface{}{"foo", "-inf", "(10"}},
		},
		{
			build: rq.Count,
			mods:  []rq.Modifier{rq.Key("foo"), rq.Eq(10)},
			cmd:   &rq.Command{Name: "ZCOUNT", Args: []interface{}{"foo", 10, 10}},
		},
		{
			build: rq.Count,
			mods:  []rq.Modifier{rq.Key("foo"), rq.GtEq(6), rq.LtEq(10)},
			cmd:   &rq.Command{Name: "ZCOUNT", Args: []interface{}{"foo", 6, 10}},
		},
	}

	for _, c := range cases {
		t.Run(c.cmd.String(), func(t *testing.T) {
			cmd, err := c.build(c.mods...).Build()

			if c.isErr {
				if err == nil {
					t.Error("should return an error")
				}

				if cmd != nil {
					t.Errorf("returned %v, want nil", cmd)
				}
			} else {
				if err != nil {
					t.Errorf("returned %v, want nil", err)
				}

				if got, want := cmd, c.cmd; !reflect.DeepEqual(got, want) {
					t.Errorf("returned %v, want %v", got, want)
				}
			}
		})
	}
}
