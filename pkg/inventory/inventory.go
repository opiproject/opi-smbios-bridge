// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022-2023 Dell Inc, or its subsidiaries.

// Package inventory is the main package of the application
package inventory

import (
	"context"
	"fmt"
	"log"

	"github.com/jaypipes/ghw"
	pc "github.com/opiproject/opi-api/common/v1/gen/go"
)

// Server contains inventory related OPI services
type Server struct {
	pc.UnimplementedInventorySvcServer
}

// GetInventory returns inventory information
func (s *Server) GetInventory(_ context.Context, in *pc.GetInventoryRequest) (*pc.Inventory, error) {
	log.Printf("GetInventory: Received from client: %v", in)

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

	pci, err := ghw.PCI()
	if err != nil {
		fmt.Printf("Error getting pci info: %v", err)
	}
	Blobarray := make([]*pc.PCIeDeviceInfo, len(pci.Devices))
	for i, r := range pci.Devices {
		fmt.Printf("PCI=%v\n", r)
		Blobarray[i] = &pc.PCIeDeviceInfo{Driver: r.Driver, Address: r.Address, Vendor: r.Vendor.Name, Product: r.Product.Name, Revision: r.Revision, Subsystem: r.Subsystem.Name, Class: r.Class.Name, Subclass: r.Subclass.Name}
	}

	return &pc.Inventory{
		Bios:      &pc.BIOSInfo{Vendor: bios.Vendor, Version: bios.Version, Date: bios.Date},
		System:    &pc.SystemInfo{Family: product.Family, Name: product.Name, Vendor: product.Vendor, SerialNumber: product.SerialNumber, Uuid: product.UUID, Sku: product.SKU, Version: product.Version},
		Baseboard: &pc.BaseboardInfo{AssetTag: baseboard.AssetTag, SerialNumber: baseboard.SerialNumber, Vendor: baseboard.Vendor, Version: baseboard.Version, Product: baseboard.Product},
		Chassis:   &pc.ChassisInfo{AssetTag: chassis.AssetTag, SerialNumber: chassis.SerialNumber, Type: chassis.Type, TypeDescription: chassis.TypeDescription, Vendor: chassis.Vendor, Version: chassis.Version},
		Processor: &pc.CPUInfo{TotalCores: int32(cpu.TotalCores), TotalThreads: int32(cpu.TotalThreads)},
		Memory:    &pc.MemoryInfo{TotalPhysicalBytes: memory.TotalPhysicalBytes, TotalUsableBytes: memory.TotalUsableBytes},
		Pci:       Blobarray,
	}, nil
}
