# OPI inventory gRPC to SMBIOS DMI bridge

[![Linters](https://github.com/opiproject/opi-smbios-bridge/actions/workflows/linters.yml/badge.svg)](https://github.com/opiproject/opi-smbios-bridge/actions/workflows/linters.yml)
[![tests](https://github.com/opiproject/opi-smbios-bridge/actions/workflows/go.yml/badge.svg)](https://github.com/opiproject/opi-smbios-bridge/actions/workflows/go.yml)
[![Docker](https://github.com/opiproject/opi-smbios-bridge/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/opiproject/opi-smbios-bridge/actions/workflows/docker-publish.yml)
[![License](https://img.shields.io/github/license/opiproject/opi-smbios-bridge?style=flat-square&color=blue&label=License)](https://github.com/opiproject/opi-smbios-bridge/blob/master/LICENSE)
[![codecov](https://codecov.io/gh/opiproject/opi-smbios-bridge/branch/main/graph/badge.svg)](https://codecov.io/gh/opiproject/opi-smbios-bridge)
[![Go Report Card](https://goreportcard.com/badge/github.com/opiproject/opi-smbios-bridge)](https://goreportcard.com/report/github.com/opiproject/opi-smbios-bridge)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/opiproject/opi-smbios-bridge)
[![Pulls](https://img.shields.io/docker/pulls/opiproject/opi-smbios-bridge.svg?logo=docker&style=flat&label=Pulls)](https://hub.docker.com/r/opiproject/opi-smbios-bridge)
[![Last Release](https://img.shields.io/github/v/release/opiproject/opi-smbios-bridge?label=Latest&style=flat-square&logo=go)](https://github.com/opiproject/opi-smbios-bridge/releases)
[![GitHub stars](https://img.shields.io/github/stars/opiproject/opi-smbios-bridge.svg?style=flat-square&label=github%20stars)](https://github.com/opiproject/opi-smbios-bridge)
[![GitHub Contributors](https://img.shields.io/github/contributors/opiproject/opi-smbios-bridge.svg?style=flat-square)](https://github.com/opiproject/opi-smbios-bridge/graphs/contributors)

This is a [SMBIOS](https://www.dmtf.org/standards/smbios) plugin to OPI inventory gRPC APIs based on [dmidecode](https://linux.die.net/man/8/dmidecode) and [ghw](https://github.com/jaypipes/ghw) go library implementing [protobuf](https://github.com/opiproject/opi-api/blob/main/common/v1/inventory.proto).

## I Want To Contribute

This project welcomes contributions and suggestions.  We are happy to have the Community involved via submission of **Issues and Pull Requests** (with substantive content or even just fixes). We are hoping for the documents, test framework, etc. to become a community process with active engagement.  PRs can be reviewed by by any number of people, and a maintainer may accept.

See [CONTRIBUTING](https://github.com/opiproject/opi/blob/main/CONTRIBUTING.md) and [GitHub Basic Process](https://github.com/opiproject/opi/blob/main/doc-github-rules.md) for more details.

## Getting started

build like this:

```bash
go build -v -o /opi-smbios-bridge ./cmd/...
```

import like this:

```go
import "github.com/opiproject/opi-smbios-bridge/pkg/inventory"
```

## Using docker

on DPU/IPU (i.e. with IP=10.10.10.1) run

```bash
$ docker run --rm -it -p 50051:50051 ghcr.io/opiproject/opi-smbios-bridge:main
2022/11/29 00:03:55 plugin serevr is &{{}}
2022/11/29 00:03:55 server listening at [::]:50051
```

on X86 management VM run

```bash
$ docker run --network=host --rm -it namely/grpc-cli ls 10.10.10.10:50051 opi_api.inventory.v1.InventorySvc -l
filename: inventory.proto
package: opi_api.inventory.v1;
service InventorySvc {
  rpc InventoryGet(opi_api.inventory.v1.InventoryGetRequest) returns (opi_api.inventory.v1.InventoryGetResponse) {}
}

# NVIDIA example:

$ docker run --network=host --rm -it namely/grpc-cli call --json_input --json_output 10.10.10.10:50051 InventoryGet "{}"
connecting to 10.10.10.10:50051
{
  "bios": {
    "vendor": "https://www.mellanox.com",
    "version": "BlueField:3.7.0-20-g98daf29",
    "date": "Jun 26 2021"
  },
  "system": {
    "family": "BlueField",
    "name": "BlueField SoC",
    "vendor": "https://www.mellanox.com",
    "serialNumber": "Unspecified System Serial Number",
    "uuid": "2e3bc1d1-e205-4830-a817-968ed1978bac",
    "sku": "Unspecified System SKU",
    "version": "1.0.0"
  },
  "baseboard": {
    "assetTag": "Unspecified Asset Tag",
    "serialNumber": "Unspecified Base Board Serial Number",
    "vendor": "https://www.mellanox.com",
    "version": "1.0.0",
    "product": "BlueField SoC"
  },
  "chassis": {
    "assetTag": "Unspecified Chassis Board Asset Tag",
    "serialNumber": "Unspecified Chassis Board Serial Number",
    "type": "1",
    "typeDescription": "Other",
    "vendor": "https://www.mellanox.com",
    "version": "1.0.0"
  },
  "processor": {
    "totalCores": 8,
    "totalThreads": 8
  },
  "memory": {
    "totalPhysicalBytes": "17179869184",
    "totalUsableBytes": "16733876224"
  },
  "pci": [
    {
      "driver": "pcieport",
      "address": "0000:00:00.0",
      "vendor": "Mellanox Technologies",
      "product": "MT42822 BlueField-2 SoC Crypto enabled",
      "revision": "0x00",
      "subsystem": "unknown",
      "class": "Bridge",
      "subclass": "PCI bridge"
    },
    {
      "driver": "pcieport",
      "address": "0000:01:00.0",
      "vendor": "Mellanox Technologies",
      "product": "MT42822 Family [BlueField-2 SoC PCIe Bridge]",
      "revision": "0x00",
      "subsystem": "unknown",
      "class": "Bridge",
      "subclass": "PCI bridge"
    },
    {
      "driver": "pcieport",
      "address": "0000:02:00.0",
      "vendor": "Mellanox Technologies",
      "product": "MT42822 Family [BlueField-2 SoC PCIe Bridge]",
      "revision": "0x00",
      "subsystem": "unknown",
      "class": "Bridge",
      "subclass": "PCI bridge"
    },
    {
      "driver": "mlx5_core",
      "address": "0000:03:00.0",
      "vendor": "Mellanox Technologies",
      "product": "MT42822 BlueField-2 integrated ConnectX-6 Dx network controller",
      "revision": "0x00",
      "subsystem": "unknown",
      "class": "Network controller",
      "subclass": "Ethernet controller"
    },
    {
      "driver": "mlx5_core",
      "address": "0000:03:00.1",
      "vendor": "Mellanox Technologies",
      "product": "MT42822 BlueField-2 integrated ConnectX-6 Dx network controller",
      "revision": "0x00",
      "subsystem": "unknown",
      "class": "Network controller",
      "subclass": "Ethernet controller"
    }
  ]
}
Rpc succeeded with OK status

# MARVELL example:

docker run --network=host --rm -it namely/grpc-cli call --json_input --json_output 11.11.11.11:50051 InventoryGet "{}"
connecting to 11.11.11.11:50051
{
 "bios": {
  "vendor": "U-Boot",
  "version": "2020.10-6.0.0",
  "date": "08/25/2022"
 },
 "system": {
  "name": "cn10k",
  "vendor": "Marvell",
  "serialNumber": "CN106-A1-CRB-R1-143",
  "uuid": "30314e43-2d36-3141-2d43-52422d52312d"
 },
 "baseboard": {
  "vendor": "Marvell",
  "product": "cn10k"
 },
 "chassis": {
  "type": "23",
  "typeDescription": "Rack mount chassis",
  "vendor": "Marvell"
 },
 "processor": {},
 "memory": {
  "totalPhysicalBytes": "51099402240",
  "totalUsableBytes": "51099402240"
 }
}
Rpc succeeded with OK status
```
