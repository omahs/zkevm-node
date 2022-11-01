package aggregatorv2

import "context"

type proverInterface interface {
	ID() string
	IsIdle() bool
	AggregateProofs(ctx context.Context) error
	VerifyBatch(ctx context.Context) error
}
