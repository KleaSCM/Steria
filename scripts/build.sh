#!/bin/sh
set -e

cd "$(dirname "$0")/.."

PROJECT="steria"
OUTPUT="bin/${PROJECT}"

echo "Building ${PROJECT}..."
go build -o "${OUTPUT}" .

echo "Running vet..."
go vet ./...

echo "Running tests..."
go test ./... -v -count=1 2>&1 | tail -20

echo "Build complete: ${OUTPUT}"
