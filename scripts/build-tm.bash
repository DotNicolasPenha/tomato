#!/usr/bin/env bash

set -e

APP_NAME="tm"
ENTRY="./main.go"
OUT_DIR="./bin"

mkdir -p "$OUT_DIR"

echo "==> Linux"
GOOS=linux   GOARCH=amd64 go build -o "$OUT_DIR/${APP_NAME}-linux-amd64"   "$ENTRY"
GOOS=linux   GOARCH=arm64 go build -o "$OUT_DIR/${APP_NAME}-linux-arm64"   "$ENTRY"

echo "==> macOS"
GOOS=darwin  GOARCH=amd64 go build -o "$OUT_DIR/${APP_NAME}-darwin-amd64"  "$ENTRY"
GOOS=darwin  GOARCH=arm64 go build -o "$OUT_DIR/${APP_NAME}-darwin-arm64"  "$ENTRY"

echo "==> Windows"
GOOS=windows GOARCH=amd64 go build -o "$OUT_DIR/${APP_NAME}-windows-amd64.exe" "$ENTRY"

echo "✔ Builds generate in $OUT_DIR"