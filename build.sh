#!/bin/sh

set -e

go build cmd/server/server.go
go build cmd/client/client.go
