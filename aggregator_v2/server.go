package aggregatorv2

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/aggregator_v2/pb"
	"github.com/0xPolygonHermez/zkevm-node/aggregator_v2/prover"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type Server struct {
	pb.UnimplementedAggregatorServiceServer

	cfg  *ServerConfig
	srv  *grpc.Server
	ctx  context.Context
	exit context.CancelFunc
}

func NewServer(cfg *ServerConfig) *Server {
	return &Server{cfg: cfg}
}

// Start sets up the server to process requests.
func (s *Server) Start(ctx context.Context) {
	var cancel context.CancelFunc
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel = context.WithCancel(ctx)
	s.ctx = ctx
	s.exit = cancel
	address := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s.srv = grpc.NewServer()
	pb.RegisterAggregatorServiceServer(s.srv, s)

	healthService := newHealthChecker()
	grpc_health_v1.RegisterHealthServer(s.srv, healthService)

	go func() {
		log.Infof("Server listening on port %d", s.cfg.Port)
		if err := s.srv.Serve(lis); err != nil {
			s.exit()
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}

// Stop stops the server.
func (s *Server) Stop() {
	s.exit()
	s.srv.Stop()
}

// Channel implements the bi-directional communication channel between the
// Prover client and the Aggregator server.
func (s *Server) Channel(stream pb.AggregatorService_ChannelServer) error {
	prover, err := prover.New(stream)
	if err != nil {
		return err
	}
	log.Debugf("establishing stream connection for prover %s", prover.ID())

	ctx := stream.Context()

	go s.handle(ctx, prover)

	// keep this scope alive, the stream gets closed if we exit from here.
	for {
		select {
		case <-s.ctx.Done():
			// server disconnecting
			return nil
		case <-ctx.Done():
			// client disconnected
			// TODO(pg): reconnect?
			return nil
		}
	}
}

func (s *Server) handle(ctx context.Context, prover proverInterface) {
	ticker := time.NewTicker(s.cfg.IntervalToConsolidateState.Duration)
	for {
		select {
		case <-s.ctx.Done():
			// server disconnected
			return
		case <-ctx.Done():
			// client disconnected
			return
		case <-ticker.C:
			if prover.IsIdle() {
				log.Debugf("prover id %s status is IDLE", prover.ID())
				// TODO(pg): aggregate or batch prove
				if err := prover.Aggregate(ctx); err != nil {
					// TODO(pg): inspect error
					// if error == nothing to aggregate --> batch prove here
				}

			} else {
				log.Debugf("prover id %s status is not IDLE", prover.ID())
			}
		}
	}
}

// HealthChecker will provide an implementation of the HealthCheck interface.
type healthChecker struct{}

// NewHealthChecker returns a health checker according to standard package
// grpc.health.v1.
func newHealthChecker() *healthChecker {
	return &healthChecker{}
}

// HealthCheck interface implementation.

// Check returns the current status of the server for unary gRPC health requests,
// for now if the server is up and able to respond we will always return SERVING.
func (s *healthChecker) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	log.Info("Serving the Check request for health check")
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

// Watch returns the current status of the server for stream gRPC health requests,
// for now if the server is up and able to respond we will always return SERVING.
func (s *healthChecker) Watch(req *grpc_health_v1.HealthCheckRequest, server grpc_health_v1.Health_WatchServer) error {
	log.Info("Serving the Watch request for health check")
	return server.Send(&grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})
}

func (s *Server) waitTick(ctx context.Context, ticker *time.Ticker) {
	select {
	case <-ticker.C:
		// nothing
	case <-s.ctx.Done():
		return
	case <-ctx.Done():
		return
	}
}
