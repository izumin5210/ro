package store

import (
	"reflect"

	"github.com/garyburd/redigo/redis"
)

// Set implements the types.Store interface.
func (s *ConcreteStore) Set(src interface{}) error {
	var err error

	conn := s.getConn()
	defer conn.Close()

	err = conn.Send("MULTI")
	if err != nil {
		return err
	}

	rv := reflect.ValueOf(src)
	if rv.Kind() == reflect.Slice {
		for i := 0; i < rv.Len(); i++ {
			err = s.set(conn, rv.Index(i))
			if err != nil {
				break
			}
		}
	} else {
		err = s.set(conn, rv)
	}

	if err != nil {
		conn.Do("DISCARD")
		return err
	}

	_, err = conn.Do("EXEC")
	return err
}

func (s *ConcreteStore) set(conn redis.Conn, src reflect.Value) error {
	m, err := s.toModel(src)
	if err != nil {
		return nil
	}

	key := s.getKey(m)

	err = conn.Send("HMSET", redis.Args{}.Add(key).AddFlat(m)...)
	if err != nil {
		return err
	}
	zsetKeys := make([]string, len(s.ScorerFuncs), len(s.ScorerFuncs))
	for _, f := range s.ScorerFuncs {
		ks, score := f(m)
		scoreSetKey := s.getScoreSetKey(ks)
		err = conn.Send("ZADD", scoreSetKey, score, key)
		if err != nil {
			return err
		}
		zsetKeys = append(zsetKeys, scoreSetKey)
	}

	err = conn.Send("SADD", redis.Args{}.Add(s.getZsetKeysKey(m)).AddFlat(zsetKeys)...)
	if err != nil {
		return err
	}

	return nil
}
