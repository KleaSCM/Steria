#!/bin/sh
set -e

cd "$(dirname "$0")/.."

PROJECT="steria"
OUTPUT="bin/${PROJECT}"

echo "=== Building ${PROJECT} ==="
go build -o "${OUTPUT}" .

echo ""
echo "=== Running go vet ==="
go vet ./...

echo ""
echo "=== Running all tests ==="
go test ./... -v -count=1

echo ""
echo "=== Build + Test complete ==="
