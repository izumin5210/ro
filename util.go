package ro

import (
	"fmt"
	"reflect"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"

	"github.com/izumin5210/ro/rq"
)

func (s *redisStore) getKey(m Model) (string, error) {
	suffix := m.GetKeySuffix()
	if len(suffix) == 0 {
		return "", errors.New("GetKeySuffix() should be present")
	}
	return s.KeyPrefix + s.KeyDelimiter + suffix, nil
}

func (s *redisStore) getScoreSetKey(key string) string {
	return s.KeyPrefix + s.ScoreKeyDelimiter + key
}

func (s *redisStore) getScoreSetKeysKeyByKey(key string) string {
	return key + s.KeyDelimiter + s.ScoreSetKeysKeySuffix
}

func (s *redisStore) toModel(rv reflect.Value) (Model, error) {
	if rv.Type() != s.modelType && rv.Type().Elem() != s.modelType {
		return nil, fmt.Errorf("%s is not a %v", rv.Interface(), s.modelType)
	}

	m, ok := rv.Interface().(Model)
	if !ok {
		return nil, fmt.Errorf("failed to cast %v to ro.IModel", rv.Interface())
	}

	if len(m.GetKeySuffix()) == 0 {
		return nil, fmt.Errorf("%v.GetKeySuffix() should be present", m)
	}

	return m, nil
}

func (s *redisStore) selectKeys(mods []rq.Modifier) ([]string, error) {
	conn := s.pool.Get()
	defer conn.Close()

	cmd, err := s.injectKeyPrefix(rq.List(mods...)).Build()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return redis.Strings(conn.Do(cmd.Name, cmd.Args...))
}

func (s *redisStore) injectKeyPrefix(q *rq.Query) *rq.Query {
	if q.Key.Prefix == "" {
		q.Key.Prefix = s.KeyPrefix
	}
	return q
}
