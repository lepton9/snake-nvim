package main

import (
	"snake-nvim.lepton9/pkg/server"
)

func main() {
	server := server.UDPServer{
		Port: 5000,
		IP:   "0.0.0.0",
	}
	server.Start()
}
