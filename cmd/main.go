// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Dell Inc, or its subsidiaries.

// main is the main package of the application
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"plugin"

	pc "github.com/opiproject/opi-api/common/v1/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()
	// Load the plugin
	plug, err := plugin.Open("/opi-smbios-bridge.so")
	if err != nil {
		log.Fatal(err)
	}
	// 2. Look for an exported symbol such as a function or variable
	inventorySymbol, err := plug.Lookup("PluginInventory")
	if err != nil {
		log.Fatal(err)
	}
	// 3. Attempt to cast the symbol to the Shipper
	var inventory pc.InventorySvcServer
	inventory, ok := inventorySymbol.(pc.InventorySvcServer)
	if !ok {
		log.Fatal("Invalid inventory type")
	}
	log.Printf("plugin serevr is %v", inventory)
	// 4. If everything is ok from the previous assertions, then we can proceed
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pc.RegisterInventorySvcServer(s, inventory)
	reflection.Register(s)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
