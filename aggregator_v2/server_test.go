package aggregatorv2

import (
	"context"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/aggregator_v2/mocks"
	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
)

func init() {
	log.Init(log.Config{
		Level:   "debug",
		Outputs: []string{"stdout"},
	})
}

func newTestServer(cfg *ServerConfig) *Server {
	return NewServer(cfg)
}

func defaultTestServerConfig() *ServerConfig {
	return &ServerConfig{
		Host:                       "0.0.0.0",
		Port:                       8888,
		IntervalToConsolidateState: types.NewDuration(time.Millisecond),
		IntervalFrequencyToGetProofGenerationState: types.NewDuration(5 * time.Second),
	}
}

func TestHandle(t *testing.T) {
	testCases := []struct {
		name  string
		setup func(*mocks.ProverMock)
	}{
		{
			name: "idle",
			setup: func(p *mocks.ProverMock) {
				p.On("ID").Return("id0")
				p.On("IsIdle").Return(true)
			},
		},
		{
			name: "not idle",
			setup: func(p *mocks.ProverMock) {
				p.On("ID").Return("id0")
				p.On("IsIdle").Return(false)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			s := newTestServer(defaultTestServerConfig())
			p := mocks.NewProverMock(t)
			tc.setup(p)
			s.Start(ctx)
			defer s.Stop()
			time.Sleep(time.Millisecond)

			go func() {
				s.handle(ctx, p)
			}()

			time.Sleep(50 * time.Millisecond)
		})
	}
}
