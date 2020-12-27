#!/bin/sh
set -e

if ! command -v protoc-gen-go &> /dev/null; then
  echo "installing protoc-gen-go"
  go get google.golang.org/protobuf/cmd/protoc-gen-go
fi

if ! command -v protoc-gen-go-grpc &> /dev/null; then
  echo "installing protoc-gen-go-grpc"
  go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
fi

SCRIPT_DIR="$(dirname $0)"
OUT_DIR="$SCRIPT_DIR/protos/protos"

rm -r "$OUT_DIR" || true
mkdir -p "$OUT_DIR" || true

export PATH="$PATH:$(go env GOPATH)/bin"

protoc \
  -I protos/ \
  protos/*.proto \
  --go_out="$OUT_DIR" \
  --go-grpc_out="$OUT_DIR"
