#!/bin/sh

export RELEASES=${PWD}/RELEASES

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-extldflags=-static" -o "$RELEASES/linux/"
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-extldflags=-static" -o "$RELEASES/macos/"
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-extldflags=-static" -o "$RELEASES/windows/"