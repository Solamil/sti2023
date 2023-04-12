#!/bin/sh
export GOROOT=/usr/local/go
export PATH=$GOROOT/bin:"$PATH"
web_dir="web/"
go test -v -coverprofile ./...
go tool cover -html=cover.out -o ${web_dir}cover.html
kill -2 "$(pgrep main -u "$(whoami)")" 2>&1
setsid -f go run ./cmd/sti2023/main.go >/dev/null 2>&1
