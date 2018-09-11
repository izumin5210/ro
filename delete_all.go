package ro

import (
	"github.com/pkg/errors"

	"github.com/izumin5210/ro/rq"
)

// DeleteAll implements the types.Store interface.
func (s *redisStore) DeleteAll(mods ...rq.Modifier) error {
	keys, err := s.selectKeys(mods)
	if err != nil {
		return errors.WithStack(err)
	}
	err = s.deleteByKeys(keys)
	if err != nil {
		return errors.Wrapf(err, "failed to remove by keys %v", keys)
	}
	return nil
}
