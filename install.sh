#!/bin/bash
go mod download
go build -o pstated -ldflags="-s -w" cmd/server/main.go
go build -o pdctl -ldflags="-s -w" cmd/client/main.go
sudo cp pstated /usr/bin/ -v
sudo cp pdctl /usr/bin/ -v
sudo cp auto-pstate.service /etc/systemd/system/ -v

# Creation of pstated daemon
sudo mkdir /run/pstated
sudo touch /run/pstated/pstated.sock
sudo chmod 0660 /run/pstated/pstated.sock
sudo systemctl enable --now auto-pstate
