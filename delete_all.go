package ro

import (
	"context"

	"github.com/pkg/errors"

	"github.com/izumin5210/ro/rq"
)

// DeleteAll implements the types.Store interface.
func (s *redisStore) DeleteAll(ctx context.Context, mods ...rq.Modifier) error {
	keys, err := s.selectKeys(ctx, mods)
	if err != nil {
		return errors.WithStack(err)
	}
	err = s.deleteByKeys(ctx, keys)
	if err != nil {
		return errors.Wrapf(err, "failed to remove by keys %v", keys)
	}
	return nil
}
