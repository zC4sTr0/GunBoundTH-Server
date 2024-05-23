package main

import (
	"net"

	"github.com/zC4sTr0/GunBoundTH-Server/broker"
)

func main() {
	serverOptions := []broker.ServerOption{
		{
			ServerName:        "GunBound Server 1",
			ServerDescription: "Open Area",
			ServerAddress:     "127.0.0.1",
			ServerPort:        12345,
			ServerUtilization: 0,
			ServerCapacity:    100,
			ServerEnabled:     true,
		},
	}

	worldSession := []net.Conn{}

	broker := broker.NewBrokerServer("localhost", 8080, serverOptions, worldSession)
	broker.Listen()
}
