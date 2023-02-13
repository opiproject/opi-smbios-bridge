// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Dell Inc, or its subsidiaries.

package main

import (
	"context"
	"log"
	"net"

	// "reflect"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	pc "github.com/opiproject/opi-api/common/v1/gen/go"
)

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	pc.RegisterInventorySvcServer(server, &PluginInventory)

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func Test_InventoryGet(t *testing.T) {
	tests := []struct {
		name    string
		out     *pc.InventoryGetResponse
		errCode codes.Code
		errMsg  string
	}{
		{
			"valid request with valid response",
			&pc.InventoryGetResponse{
				Bios:      &pc.BIOSInfo{Vendor: "TBD", Version: "TBD", Date: "TBD"},
				System:    &pc.SystemInfo{Family: "TBD", Name: "TBD", Vendor: "TBD", SerialNumber: "TBD", Uuid: "TBD", Sku: "TBD", Version: "TBD"},
				Baseboard: &pc.BaseboardInfo{AssetTag: "TBD", SerialNumber: "TBD", Vendor: "TBD", Version: "TBD", Product: "TBD"},
				Chassis:   &pc.ChassisInfo{AssetTag: "TBD", SerialNumber: "TBD", Type: "TBD", TypeDescription: "TBD", Vendor: "TBD", Version: "TBD"},
				Processor: &pc.CPUInfo{TotalCores: 8, TotalThreads: 16},
				Memory:    &pc.MemoryInfo{TotalPhysicalBytes: 12, TotalUsableBytes: 55},
			},
			codes.OK,
			"",
		},
	}

	// start GRPC mockup server
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)
	client := pc.NewInventorySvcClient(conn)

	// run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := &pc.InventoryGetRequest{}
			response, err := client.InventoryGet(ctx, request)
			if response != nil {
				// if !reflect.DeepEqual(response, tt.out) {
				if response.Bios.Vendor == "" {
					t.Error("response: expected", tt.out, "received", response)
				}
			}

			if err != nil {
				if er, ok := status.FromError(err); ok {
					if er.Code() != tt.errCode {
						t.Error("error code: expected", codes.InvalidArgument, "received", er.Code())
					}
					if er.Message() != tt.errMsg {
						t.Error("error message: expected", tt.errMsg, "received", er.Message())
					}
				}
			}
		})
	}
}
