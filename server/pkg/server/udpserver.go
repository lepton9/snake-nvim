package server

import (
	"fmt"
	"net"
)

type UDPServer struct {
	Port int
	IP   string
}

func (s *UDPServer) Start() {
	addr := net.UDPAddr{
		Port: s.Port,
		IP:   net.ParseIP(s.IP),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer conn.Close()

	fmt.Printf("UDP server listening on %s:%d\n", s.IP, s.Port)

	buffer := make([]byte, 1024)

	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}

		fmt.Printf("Received %d bytes from %s: %s\n", n, remoteAddr, string(buffer[:n]))

		_, err = conn.WriteToUDP(buffer[:n], remoteAddr)
		if err != nil {
			fmt.Println("Error writing to UDP:", err)
		}
	}
}
