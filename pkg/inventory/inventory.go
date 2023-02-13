// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Dell Inc, or its subsidiaries.

// main is the main package of the application
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jaypipes/ghw"
	pc "github.com/opiproject/opi-api/common/v1/gen/go"
)

type server struct {
	pc.UnimplementedInventorySvcServer
}

// PluginInventory is the server that we export to load dynamically at runtime
var PluginInventory server

func (s *server) InventoryGet(ctx context.Context, in *pc.InventoryGetRequest) (*pc.InventoryGetResponse, error) {
	log.Printf("InventoryGet: Received from client: %v", in)

	cpu, err := ghw.CPU()
	if err != nil {
		fmt.Printf("Error getting CPU info: %v", err)
		// return nil, status.Errorf(codes.InvalidArgument, msg)
	}
	fmt.Printf("%v\n", cpu)

	memory, err := ghw.Memory()
	if err != nil {
		fmt.Printf("Error getting memory info: %v", err)
	}
	fmt.Println(memory.String())

	chassis, err := ghw.Chassis()
	if err != nil {
		fmt.Printf("Error getting chassis info: %v", err)
	}
	fmt.Printf("%v\n", chassis)

	bios, err := ghw.BIOS()
	if err != nil {
		fmt.Printf("Error getting BIOS info: %v", err)
	}
	fmt.Printf("%v\n", bios)

	baseboard, err := ghw.Baseboard()
	if err != nil {
		fmt.Printf("Error getting baseboard info: %v", err)
	}
	fmt.Printf("%v\n", baseboard)

	product, err := ghw.Product()
	if err != nil {
		fmt.Printf("Error getting product info: %v", err)
	}
	fmt.Printf("%v\n", product)

	return &pc.InventoryGetResponse{
		Bios:      &pc.BIOSInfo{Vendor: bios.Vendor, Version: bios.Version, Date: bios.Date},
		System:    &pc.SystemInfo{Family: product.Family, Name: product.Name, Vendor: product.Vendor, SerialNumber: product.SerialNumber, Uuid: product.UUID, Sku: product.SKU, Version: product.Version},
		Baseboard: &pc.BaseboardInfo{AssetTag: baseboard.AssetTag, SerialNumber: baseboard.SerialNumber, Vendor: baseboard.Vendor, Version: baseboard.Version, Product: baseboard.Product},
		Chassis:   &pc.ChassisInfo{AssetTag: chassis.AssetTag, SerialNumber: chassis.SerialNumber, Type: chassis.Type, TypeDescription: chassis.TypeDescription, Vendor: chassis.Vendor, Version: chassis.Version},
		Processor: &pc.CPUInfo{TotalCores: int32(cpu.TotalCores), TotalThreads: int32(cpu.TotalThreads)},
		Memory:    &pc.MemoryInfo{TotalPhysicalBytes: memory.TotalPhysicalBytes, TotalUsableBytes: memory.TotalUsableBytes},
	}, nil
}
