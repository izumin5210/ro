package ro

import (
	"github.com/pkg/errors"

	"github.com/izumin5210/ro/rq"
)

// RemoveBy implements the types.Store interface.
func (s *redisStore) RemoveBy(mods ...rq.Modifier) error {
	keys, err := s.selectKeys(mods)
	if err != nil {
		return errors.WithStack(err)
	}
	err = s.removeByKeys(keys)
	if err != nil {
		return errors.Wrapf(err, "failed to remove by keys %v", keys)
	}
	return nil
}
