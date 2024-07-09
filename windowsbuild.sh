#!/bin/bash

# Set the target OS and architecture
export GOOS=windows
export GOARCH=amd64

# Build the executable
go build -o file-archiver-go.exe file-archiver-go.go
