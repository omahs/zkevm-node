package main

import (
	aggregatorv2 "github.com/0xPolygonHermez/zkevm-node/aggregator_v2"
)

func main() {
	cfg := aggregatorv2.ServerConfig{
		Host: "0.0.0.0",
		Port: 8888,
	}

	srv := aggregatorv2.NewServer(&cfg)

	srv.Start()
}
