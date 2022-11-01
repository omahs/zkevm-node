package aggregatorv2

import "github.com/0xPolygonHermez/zkevm-node/config/types"

// ServerConfig represents the configuration of the aggregator server.
type ServerConfig struct {
	Host string `mapstructure:"Host"`
	Port int    `mapstructure:"Port"`

	// IntervalToConsolidateState is the time the aggregator waits until
	// trying to consolidate a new state
	IntervalToConsolidateState types.Duration `mapstructure:"IntervalToConsolidateState"`
	// IntervalFrequencyToGetProofGenerationState is the time the aggregator waits until
	// trying to get proof generation status, in case prover client returns PENDING state
	IntervalFrequencyToGetProofGenerationState types.Duration `mapstructure:"IntervalFrequencyToGetProofGenerationState"`
}
