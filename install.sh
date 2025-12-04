#!/bin/bash
go mod download
go build -o pstated -ldflags="-s -w" cmd/server/main.go
go build -o pdctl -ldflags="-s -w" cmd/client/main.go
cp pstated /usr/bin/ -v
cp pdctl /usr/bin/ -v
cp auto-pstate.service /etc/systemd/system/ -v
mkdir /run/pstated
systemctl enable --now auto-pstate
rm -rf pstated pdctl
