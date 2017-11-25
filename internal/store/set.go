package store

import (
	"reflect"

	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"
)

// Set implements the types.Store interface.
func (s *ConcreteStore) Set(src interface{}) error {
	var err error

	conn := s.getConn()
	defer conn.Close()

	err = conn.Send("MULTI")
	if err != nil {
		return errors.Wrap(err, "faild to send MULTI command")
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
		return errors.Wrap(err, "faild to send any commands")
	}

	_, err = conn.Do("EXEC")
	if err != nil {
		return errors.Wrap(err, "faild to EXEC commands")
	}
	return nil
}

func (s *ConcreteStore) set(conn redis.Conn, src reflect.Value) error {
	m, err := s.toModel(src)
	if err != nil {
		return errors.Wrap(err, "failed to convert to model")
	}

	key := s.getKey(m)

	err = conn.Send("HMSET", redis.Args{}.Add(key).AddFlat(m)...)
	if err != nil {
		return errors.Wrapf(err, "failed to send HMEST %s %v", key, m)
	}
	zsetKeys := make([]string, len(s.ScorerFuncs), len(s.ScorerFuncs))
	for _, f := range s.ScorerFuncs {
		ks, score := f(m)
		scoreSetKey := s.getScoreSetKey(ks)
		err = conn.Send("ZADD", scoreSetKey, score, key)
		if err != nil {
			return errors.Wrapf(err, "failed to send ZADD %s %v %s", scoreSetKey, score, key)
		}
		zsetKeys = append(zsetKeys, scoreSetKey)
	}

	err = conn.Send("SADD", redis.Args{}.Add(s.getZsetKeysKey(m)).AddFlat(zsetKeys)...)
	if err != nil {
		return errors.Wrapf(err, "failed to send SADD %s %v", s.getZsetKeysKey(m), zsetKeys)
	}

	return nil
}
