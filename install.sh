#!/bin/bash
go mod download
go build -o pstated -ldflags="-s -w" cmd/server/main.go
go build -o pdctl -ldflags="-s -w" cmd/client/main.go
sudo cp pstated /usr/bin/ -v
sudo cp pdctl /usr/bin/ -v
sudo cp auto-pstate.service /etc/systemd/system/ -v
sudo mkdir /run/pstated
sudo systemctl enable --now auto-pstate
