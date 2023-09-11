// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022-2023 Dell Inc, or its subsidiaries.

// main is the main package of the application
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	pc "github.com/opiproject/opi-api/common/v1/gen/go"
	"github.com/opiproject/opi-smbios-bridge/pkg/inventory"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

var (
	grpcPort = flag.Int("grpc_port", 50051, "The gRPC server port")
	httpPort = flag.Int("http_port", 8082, "The HTTP server port")
)

func main() {
	flag.Parse()

	go runGatewayServer()
	runGrpcServer()
}

func runGrpcServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pc.RegisterInventorySvcServer(s, &inventory.Server{})
	reflection.Register(s)

	log.Printf("gRPC Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func runGatewayServer() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pc.RegisterInventorySvcHandlerFromEndpoint(ctx, mux, fmt.Sprintf(":%d", *grpcPort), opts)
	if err != nil {
		log.Panic("cannot register handler server")
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	log.Printf("HTTP Server listening at %v", *httpPort)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", *httpPort),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Panic("cannot start HTTP gateway server")
	}
}
