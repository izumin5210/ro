package store

import (
	"fmt"
	"reflect"

	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"

	"github.com/izumin5210/ro/types"
)

func (s *ConcreteStore) getKey(m types.Model) (string, error) {
	suffix := m.GetKeySuffix()
	if len(suffix) == 0 {
		return "", errors.New("GetKeySuffix() should be present")
	}
	return s.KeyPrefix + s.KeyDelimiter + suffix, nil
}

func (s *ConcreteStore) getScoreSetKey(key string) string {
	return s.KeyPrefix + s.ScoreKeyDelimiter + key
}

func (s *ConcreteStore) getScoreSetKeysKeyByKey(key string) string {
	return key + s.KeyDelimiter + s.ScoreSetKeysKeySuffix
}

func (s *ConcreteStore) toModel(rv reflect.Value) (types.Model, error) {
	if rv.Type() != s.modelType && rv.Type().Elem() != s.modelType {
		return nil, fmt.Errorf("%s is not a %v", rv.Interface(), s.modelType)
	}

	m, ok := rv.Interface().(types.Model)
	if !ok {
		return nil, fmt.Errorf("failed to cast %v to types.Model", rv.Interface())
	}

	if len(m.GetKeySuffix()) == 0 {
		return nil, fmt.Errorf("%v.GetKeySuffix() should be present", m)
	}

	return m, nil
}

func (s *ConcreteStore) selectKeys(query types.Query) ([]string, error) {
	conn := s.getConn()
	defer conn.Close()

	cmd, args := query.Build()
	return redis.Strings(conn.Do(cmd, args...))
}
