#!/bin/sh
name=$(basename $(pwd))
go mod tidy
CGO_ENABLED=0 go build -o ${name}.elf -ldflags '-s -w -extldflags "-static"' -trimpath .
