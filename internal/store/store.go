package store

import (
	"fmt"
	"reflect"

	"github.com/garyburd/redigo/redis"
	"github.com/izumin5210/ro/types"
)

const (
	keyDelimiter = ":"
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

	_, err = conn.Do("HMSET", redis.Args{}.Add(s.getKey(src)).AddFlat(src)...)
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

func (s *ConcreteStore) getKey(m types.Model) string {
	prefix := m.GetKeyPrefix()
	if len(prefix) == 0 {
		prefix = s.modelType.Name()
	}
	suffix := m.GetKeySuffix()
	return prefix + keyDelimiter + suffix
}

func (s *ConcreteStore) validate(m types.Model) error {
	if len(m.GetKeySuffix()) == 0 {
		return fmt.Errorf("%v.GetKeySuffix() should be present", m)
	}
	return nil
}
