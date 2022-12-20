# OPI inventory gRPC to SMBIOS DMI bridge

[![Linters](https://github.com/opiproject/opi-smbios-bridge/actions/workflows/linters.yml/badge.svg)](https://github.com/opiproject/opi-smbios-bridge/actions/workflows/linters.yml)
[![tests](https://github.com/opiproject/opi-smbios-bridge/actions/workflows/go.yml/badge.svg)](https://github.com/opiproject/opi-smbios-bridge/actions/workflows/go.yml)
[![Docker](https://github.com/opiproject/opi-smbios-bridge/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/opiproject/opi-smbios-bridge/actions/workflows/docker-publish.yml)
[![License](https://img.shields.io/github/license/opiproject/opi-smbios-bridge?style=flat-square&color=blue&label=License)](https://github.com/opiproject/opi-smbios-bridge/blob/master/LICENSE)
[![codecov](https://codecov.io/gh/opiproject/opi-smbios-bridge/branch/main/graph/badge.svg)](https://codecov.io/gh/opiproject/opi-smbios-bridge)
[![Go Report Card](https://goreportcard.com/badge/github.com/opiproject/opi-smbios-bridge)](https://goreportcard.com/report/github.com/opiproject/opi-smbios-bridge)
[![Last Release](https://img.shields.io/github/v/release/opiproject/opi-smbios-bridge?label=Latest&style=flat-square&logo=go)](https://github.com/opiproject/opi-smbios-bridge/releases)

This is a [SMBIOS](https://www.dmtf.org/standards/smbios) plugin to OPI inventory gRPC APIs based on dmidecode and ghw go library.

## I Want To Contribute

This project welcomes contributions and suggestions.  We are happy to have the Community involved via submission of **Issues and Pull Requests** (with substantive content or even just fixes). We are hoping for the documents, test framework, etc. to become a community process with active engagement.  PRs can be reviewed by by any number of people, and a maintainer may accept.

See [CONTRIBUTING](https://github.com/opiproject/opi/blob/main/CONTRIBUTING.md) and [GitHub Basic Process](https://github.com/opiproject/opi/blob/main/doc-github-rules.md) for more details.

## Getting started

```bash
go build -v -buildmode=plugin -o /opi-smbios-bridge.so ./...
```

 in main app:

```go
package main
import (
    "plugin"
    pc "github.com/opiproject/opi-api/common/v1/gen/go"
)
func main() {
    plug, err := plugin.Open("/opi-smbios-bridge.so")
    inventorySymbol, err := plug.Lookup("PluginInventory")
    var inventory pc.InventorySvcServer
    inventory, ok := inventorySymbol.(pc.InventorySvcServer)
    s := grpc.NewServer()
    pc.RegisterInventorySvcServer(s, inventory)
    reflection.Register(s)
}
```

## Using docker

on DPU/IPU (i.e. with IP=10.10.10.1) run

```bash
$ docker run --rm -it -v /var/tmp/:/var/tmp/ -p 50051:50051 ghcr.io/opiproject/opi-smbios-bridge:main
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
  "uuid": "63ec430e-7620-479a-8ec5-133de3972679",
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
 "processor": {},
 "memory": {
  "totalPhysicalBytes": "17179869184",
  "totalUsableBytes": "16733876224"
 }
}
Rpc succeeded with OK status
```
