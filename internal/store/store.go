package store

import (
	"fmt"
	"reflect"

	"github.com/garyburd/redigo/redis"
	"github.com/izumin5210/ro/types"
)

const (
	keyDelimiter   = ":"
	scoreDelimiter = "/"
)

// ConcreteStore is an implementation of types.Store
type ConcreteStore struct {
	*types.StoreConfig
	getConn   types.GetConnFunc
	modelType reflect.Type
}

// New creates a ConcreteStore instance
func New(getConnFunc types.GetConnFunc, model types.Model, cnf *types.StoreConfig) (types.Store, error) {
	return &ConcreteStore{
		StoreConfig: cnf,
		getConn:     getConnFunc,
		modelType:   reflect.ValueOf(model).Elem().Type(),
	}, nil
}

// Set implements the types.Store interface.
func (s *ConcreteStore) Set(src types.Model) error {
	var err error
	if err = s.validate(src); err != nil {
		return err
	}

	conn := s.getConn()
	defer conn.Close()

	key := s.getKey(src)

	err = conn.Send("MULTI")
	if err != nil {
		return err
	}
	err = conn.Send("HMSET", redis.Args{}.Add(key).AddFlat(src)...)
	if err != nil {
		return err
	}
	prefix := s.getKeyPrefix(src)
	for k, f := range s.ScorerFuncMap {
		err = conn.Send("ZADD", prefix+scoreDelimiter+k, f(src), key)
		if err != nil {
			return err
		}
	}
	_, err = conn.Do("EXEC")
	return err
}

// Get implements the types.Store interface.
func (s *ConcreteStore) Get(dest types.Model) error {
	conn := s.getConn()
	defer conn.Close()

	v, err := redis.Values(conn.Do("HGETAll", s.getKey(dest)))
	if err != nil {
		return err
	}
	err = redis.ScanStruct(v, dest)
	if err != nil {
		return err
	}
	return nil
}

func (s *ConcreteStore) getKeyPrefix(m types.Model) string {
	prefix := m.GetKeyPrefix()
	if len(prefix) == 0 {
		prefix = s.modelType.Name()
	}
	return prefix
}

func (s *ConcreteStore) getKey(m types.Model) string {
	return s.getKeyPrefix(m) + keyDelimiter + m.GetKeySuffix()
}

func (s *ConcreteStore) validate(m types.Model) error {
	if len(m.GetKeySuffix()) == 0 {
		return fmt.Errorf("%v.GetKeySuffix() should be present", m)
	}
	return nil
}
