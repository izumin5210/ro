package ro

import (
	"reflect"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
)

// Remove implements the types.Store interface.
func (s *redisStore) Remove(src interface{}) error {
	keys := []string{}

	rv := reflect.ValueOf(src)
	if rv.Kind() == reflect.Slice {
		for i := 0; i < rv.Len(); i++ {
			m, err := s.toModel(rv.Index(i))
			if err != nil {
				return errors.Wrapf(err, "failed to convert to model %v", rv.Index(i).Interface())
			}
			key, err := s.getKey(m)
			if err != nil {
				return errors.Wrap(err, "failed to get key")
			}
			keys = append(keys, key)
		}
	} else {
		m, err := s.toModel(rv)
		if err != nil {
			return errors.Wrapf(err, "failed to convert to model %v", rv.Interface())
		}
		key, err := s.getKey(m)
		if err != nil {
			return errors.Wrap(err, "failed to get key")
		}
		keys = append(keys, key)
	}

	err := s.removeByKeys(keys)
	if err != nil {
		return errors.Wrapf(err, "failed to remove by keys %v", keys)
	}
	return nil
}

func (s *redisStore) removeByKeys(keys []string) error {
	conn := s.pool.Get()
	defer conn.Close()

	keysByZsetKey := map[string][]string{}
	for _, k := range keys {
		zsetKeys, err := redis.Strings(conn.Do("SMEMBERS", s.getScoreSetKeysKeyByKey(k)))
		if err != nil {
			return errors.Wrapf(err, "failed to execute SMEMBERS %s", s.getScoreSetKeysKeyByKey(k))
		}
		for _, zk := range zsetKeys {
			keysByZsetKey[zk] = append(keysByZsetKey[zk], k)
		}
	}

	err := conn.Send("MULTI")
	if err != nil {
		return errors.Wrap(err, "faild to send MULTI command")
	}

	if len(keys) > 0 {
		err = conn.Send("DEL", redis.Args{}.AddFlat(keys)...)
		if err != nil {
			return errors.Wrapf(err, "faild to send DEL %v", keys)
		}
	}

	for zk, hkeys := range keysByZsetKey {
		err = conn.Send("ZREM", redis.Args{}.Add(zk).AddFlat(hkeys)...)
		if err != nil {
			return errors.Wrapf(err, "faild to send ZREM %s %v", zk, keys)
		}
	}

	_, err = conn.Do("EXEC")
	if err != nil {
		return errors.Wrap(err, "failed to execute EXEC")
	}
	return nil
}
