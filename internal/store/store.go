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
	modelType := reflect.ValueOf(model).Elem().Type()

	if len(cnf.KeyPrefix) == 0 {
		cnf.KeyPrefix = modelType.Name()
	}

	return &ConcreteStore{
		StoreConfig: cnf,
		getConn:     getConnFunc,
		model:       model,
		modelType:   modelType,
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
	m, err := s.toModel(src)
	if err != nil {
		return nil
	}

	key := s.getKey(m)

	err = conn.Send("HMSET", redis.Args{}.Add(key).AddFlat(m)...)
	if err != nil {
		return err
	}
	for k, f := range s.ScorerFuncMap {
		err = conn.Send("ZADD", s.KeyPrefix+scoreDelimiter+k, f(m), key)
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

func (s *ConcreteStore) getKey(m types.Model) string {
	return s.KeyPrefix + keyDelimiter + m.GetKeySuffix()
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
