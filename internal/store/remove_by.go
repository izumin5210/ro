package store

import (
	"github.com/pkg/errors"

	"github.com/izumin5210/ro/types"
)

// RemoveBy implements the types.Store interface.
func (s *ConcreteStore) RemoveBy(q types.Query) error {
	keys, err := s.selectKeys(q)
	if err != nil {
		return errors.Wrapf(err, "failed to select keys %v", q)
	}
	err = s.removeByKeys(keys)
	if err != nil {
		return errors.Wrapf(err, "failed to remove by keys %v", keys)
	}
	return nil
}
