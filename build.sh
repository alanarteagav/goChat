# Bash script for building goChat
export CGO_ENABLES="0"
export GOPATH=$(pwd)/
export GOBIN=$(pwd)/bin
go build ./...
go install ./...
gcc -o bin/goClient ./src/github.com/alanarteagav/goClient/goClient.c -Wall `pkg-config --cflags --libs gtk+-3.0` -export-dynamic
