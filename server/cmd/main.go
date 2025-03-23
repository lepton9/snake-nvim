package main

import (
	"snake-nvim.lepton9/pkg/server"
)

func main() {
	server := server.Init("0.0.0.0", 5000)
	server.Start()
}
