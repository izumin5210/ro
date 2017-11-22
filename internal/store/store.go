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
	model     types.Model
	modelType reflect.Type
}

// New creates a ConcreteStore instance
func New(getConnFunc types.GetConnFunc, model types.Model, cnf *types.StoreConfig) (types.Store, error) {
	return &ConcreteStore{
		StoreConfig: cnf,
		getConn:     getConnFunc,
		model:       model,
		modelType:   reflect.ValueOf(model).Elem().Type(),
	}, nil
}

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
	if src.Type() != s.modelType && src.Type().Elem() != s.modelType {
		return fmt.Errorf("%s is not a %v", src.Interface(), s.modelType)
	}

	m := src.Interface().(types.Model)

	if err := s.validate(m); err != nil {
		return err
	}

	key := s.getKey(m)

	err := conn.Send("HMSET", redis.Args{}.Add(key).AddFlat(m)...)
	if err != nil {
		return err
	}
	prefix := s.getKeyPrefix(m)
	for k, f := range s.ScorerFuncMap {
		err = conn.Send("ZADD", prefix+scoreDelimiter+k, f(m), key)
		if err != nil {
			return err
		}
	}

	return nil
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
