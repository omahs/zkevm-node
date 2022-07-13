package jsonrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/didip/tollbooth/v6"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	// APIEth represents the eth API prefix.
	APIEth = "eth"
	// APINet represents the net API prefix.
	APINet = "net"
	// APIDebug represents the debug API prefix.
	APIDebug = "debug"
	// APIZKEVM represents the zkevm API prefix.
	APIZKEVM = "zkevm"
	// APITxPool represents the txpool API prefix.
	APITxPool = "txpool"
	// APIWeb3 represents the web3 API prefix.
	APIWeb3 = "web3"
)

// Server is an API backend to handle RPC requests
type Server struct {
	config  Config
	handler *Handler
	srv     *http.Server
}

var (
	requests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "jsonrpc_requests",
			Help: "JSONRPC number of requests received",
		},
		[]string{"request_label"},
	)
	start           = 0.1
	width           = 0.1
	count           = 10
	requestDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "jsonrpc_request_duration",
			Help:    "JSONRPC Histogram for the runtime of requests",
			Buckets: prometheus.LinearBuckets(start, width, count),
		},
	)
)

// NewServer returns the JsonRPC server
func NewServer(cfg Config, p jsonRPCTxPool, s stateInterface,
	gpe gasPriceEstimator, storage storageInterface, apis map[string]bool) *Server {
	handler := newJSONRpcHandler()

	if _, ok := apis[APIEth]; ok {
		ethEndpoints := &Eth{cfg: cfg, pool: p, state: s, gpe: gpe, storage: storage}
		handler.registerService(APIEth, ethEndpoints)
	}

	if _, ok := apis[APINet]; ok {
		netEndpoints := &Net{cfg: cfg}
		handler.registerService(APINet, netEndpoints)
	}

	if _, ok := apis[APIZKEVM]; ok {
		hezEndpoints := &ZKEVM{state: s, config: cfg}
		handler.registerService(APIZKEVM, hezEndpoints)
	}

	if _, ok := apis[APITxPool]; ok {
		txPoolEndpoints := &TxPool{}
		handler.registerService(APITxPool, txPoolEndpoints)
	}

	if _, ok := apis[APIDebug]; ok {
		debugEndpoints := &Debug{state: s}
		handler.registerService(APIDebug, debugEndpoints)
	}

	if _, ok := apis[APIWeb3]; ok {
		web3Endpoints := &Web3{}
		handler.registerService(APIWeb3, web3Endpoints)
	}

	srv := &Server{
		config:  cfg,
		handler: handler,
	}
	return srv
}

// Start initializes the JSON RPC server to listen for request
func (s *Server) Start() error {
	if s.srv != nil {
		return fmt.Errorf("server already started")
	}

	address := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	log.Infof("http server started: %s", address)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Errorf("failed to create tcp listener: %v", err)
		return err
	}

	mux := http.NewServeMux()

	lmt := tollbooth.NewLimiter(s.config.MaxRequestsPerIPAndSecond, nil)
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/", tollbooth.LimitFuncHandler(lmt, s.handle))

	s.srv = &http.Server{
		Handler: mux,
	}
	if err := s.srv.Serve(lis); err != nil {
		if err == http.ErrServerClosed {
			log.Infof("http server stopped")
			return nil
		}
		log.Errorf("closed http connection: %v", err)
		return err
	}
	return nil
}

// Stop shutdown the rpc server
func (s *Server) Stop() error {
	if s.srv == nil {
		return nil
	}

	if err := s.srv.Shutdown(context.Background()); err != nil {
		return err
	}

	if err := s.srv.Close(); err != nil {
		return err
	}

	s.srv = nil

	return nil
}

func (s *Server) handle(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	requestLabel := "undefined"
	defer func() {
		requests.WithLabelValues("Total").Inc()
		requests.WithLabelValues(requestLabel).Inc()
	}()

	if (*req).Method == "OPTIONS" {
		requestLabel = "foo"
		return
	}

	if req.Method == "GET" {
		requestLabel = "foo"
		_, err := w.Write([]byte("zkEVM JSON RPC Server"))
		if err != nil {
			log.Error(err)
		}
		return
	}

	if req.Method != "POST" {
		requestLabel = "invalid"
		err := errors.New("method " + req.Method + " not allowed")
		handleError(w, err)
		return
	}

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		requestLabel = "invalid"
		handleError(w, err)
		return
	}

	single, err := s.isSingleRequest(data)
	if err != nil {
		requestLabel = "invalid"
		handleError(w, err)
		return
	}

	timer := prometheus.NewTimer(requestDuration)
	if single {
		requestLabel = "single"
		s.handleSingleRequest(w, data)
	} else {
		requestLabel = "batch"
		s.handleBatchRequest(w, data)
	}
	timer.ObserveDuration()
}

func (s *Server) isSingleRequest(data []byte) (bool, rpcError) {
	x := bytes.TrimLeft(data, " \t\r\n")

	if len(x) == 0 {
		return false, newRPCError(invalidRequestErrorCode, "Invalid json request")
	}

	return x[0] == '{', nil
}

func (s *Server) handleSingleRequest(w http.ResponseWriter, data []byte) {
	request, err := s.parseRequest(data)
	if err != nil {
		handleError(w, err)
		return
	}

	response := s.handler.Handle(request)

	respBytes, err := json.Marshal(response)
	if err != nil {
		handleError(w, err)
		return
	}

	_, err = w.Write(respBytes)
	if err != nil {
		handleError(w, err)
		return
	}
}

func (s *Server) handleBatchRequest(w http.ResponseWriter, data []byte) {
	requests, err := s.parseRequests(data)
	if err != nil {
		handleError(w, err)
		return
	}

	responses := make([]Response, 0, len(requests))

	for _, request := range requests {
		response := s.handler.Handle(request)
		responses = append(responses, response)
	}

	respBytes, _ := json.Marshal(responses)
	_, err = w.Write(respBytes)
	if err != nil {
		log.Error(err)
	}
}

func (s *Server) parseRequest(data []byte) (Request, error) {
	var req Request

	if err := json.Unmarshal(data, &req); err != nil {
		return Request{}, newRPCError(invalidRequestErrorCode, "Invalid json request")
	}

	return req, nil
}

func (s *Server) parseRequests(data []byte) ([]Request, error) {
	var requests []Request

	if err := json.Unmarshal(data, &requests); err != nil {
		return nil, newRPCError(invalidRequestErrorCode, "Invalid json request")
	}

	return requests, nil
}

func handleError(w http.ResponseWriter, err error) {
	log.Error(err)
	_, err = w.Write([]byte(err.Error()))
	if err != nil {
		log.Error(err)
	}
}

func rpcErrorResponse(code int, errorMessage string, err error) (interface{}, rpcError) {
	if err != nil {
		log.Errorf("%v:%v", errorMessage, err.Error())
	} else {
		log.Error(errorMessage)
	}
	return nil, newRPCError(code, errorMessage)
}
