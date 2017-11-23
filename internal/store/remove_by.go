package store

import (
	"github.com/izumin5210/ro/types"
)

// RemoveBy implements the types.Store interface.
func (s *ConcreteStore) RemoveBy(q types.Query) error {
	keys, err := s.keys(q)
	if err != nil {
		return err
	}
	return s.removeByKeys(keys)
}
