# protoc-gen-mapper

A simple Static Code generator for converting Proto structures into a Flat map. Plugin this tool into your regular proto generaors for golang
All the scalar types are converted to string types by default

## Build

```bash
go build .
```

## Installation

Installs to your go root. If you install this., no need to specify plugin option for protoc

```bash
go install
```

## Run the plugin

Multiple options to run the plugin.
Installed

```bash
protoc --mapper_out="parent=Product:." --go_out=. product.proto
```

From binary

```bash

protoc --plugin proto-gen-mapper  --mapper_out="parent=Product:." product.proto
```

Features
- multi package
- nested structure
- multi supported types: string, int32, int64, float, double, bool

To do
- repeated
- map
- enum