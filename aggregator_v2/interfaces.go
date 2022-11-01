package aggregatorv2

import "context"

type proverInterface interface {
	ID() string
	IsIdle() bool
	Aggregate(ctx context.Context) error
}
