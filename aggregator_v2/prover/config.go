package prover

import "github.com/0xPolygonHermez/zkevm-node/config/types"

// Config represents the Prover configuration.
type Config struct {
	// IntervalFrequencyToGetProofGenerationState is the time the aggregator waits until
	// trying to get proof generation status, in case prover client returns PENDING state
	IntervalFrequencyToGetProofGenerationState types.Duration `mapstructure:"IntervalFrequencyToGetProofGenerationState"`
}
