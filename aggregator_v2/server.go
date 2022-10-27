package aggregatorv2

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/aggregator_v2/pb"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type Server struct {
	pb.UnimplementedAggregatorServiceServer

	cfg     *ServerConfig
	srv     *grpc.Server
	provers sync.Map
	ctx     context.Context
	exit    context.CancelFunc
}

func NewServer(cfg *ServerConfig) *Server {
	return &Server{
		cfg: cfg,
	}
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

	log.Infof("Server listening on port %d", s.cfg.Port)
	if err := s.srv.Serve(lis); err != nil {
		s.exit()
		log.Fatalf("failed to serve: %v", err)
	}
}

// Stop stops the server.
func (s *Server) Stop() {
	s.exit()
	s.srv.Stop()
}

var counter uint64

// Channel implements the bi-directional communication channel between Prover
// client and Aggregator server.
func (s *Server) Channel(stream pb.AggregatorService_ChannelServer) error {
	count := atomic.LoadUint64(&counter)
	atomic.AddUint64(&counter, 1)
	log.Debugf("establishing stream for channel %d", count)

	proverID, err := s.proverID(stream)
	if err != nil {
		return err
	}

	for {
		s.provers.Range(func(key, value interface{}) bool {
			log.Debugf("[channel %d] asking status for prover %s", count, key.(string))
			// log.Debugf("asking status for prover %s", proverID)
			msg, err := s.getStatus(value.(pb.AggregatorService_ChannelServer))
			// msg, err := s.getStatus(stream)
			if err != nil {
				log.Error(err)
				// return err
				return false
			}
			log.Debugf("[channel %d] prover id %s status is %s", count, proverID, msg.Status.String())
			time.Sleep(1 * time.Second)
			// return nil
			return true
		})
	}
}

func (s *Server) proverID(stream pb.AggregatorService_ChannelServer) (string, error) {
	var id string
	msg, err := s.getStatus(stream)
	if err != nil {
		return id, err
	}
	id = msg.ProverId
	if _, ok := s.provers.Load(id); !ok {
		// first message
		// store the prover stream for later communication
		s.provers.Store(id, stream)
	}
	return id, nil
}

func (s *Server) getStatus(stream pb.AggregatorService_ChannelServer) (*pb.GetStatusResponse, error) {
	req := &pb.AggregatorMessage{
		Request: &pb.AggregatorMessage_GetStatusRequest{
			GetStatusRequest: &pb.GetStatusRequest{},
		},
	}
	if err := stream.Send(req); err != nil {
		return nil, err
	}

	res, err := stream.Recv()
	if err != nil {
		return nil, err
	}
	if msg, ok := res.Response.(*pb.ProverMessage_GetStatusResponse); ok {
		return msg.GetStatusResponse, nil
	}
	return nil, errors.New("bad response") // FIXME(pg)
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
