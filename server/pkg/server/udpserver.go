package server

import (
	"fmt"
	"net"
	"snake-nvim.lepton9/pkg/player"
	"sync"
)

type UDPServer struct {
	port             int
	ip               string
	connectedPlayers map[string]player.Player
	mu               sync.Mutex
}

func Init(ip string, port int) *UDPServer {
	server := UDPServer{
		port:             port,
		ip:               ip,
		connectedPlayers: make(map[string]player.Player),
	}
	return &server
}

func (s *UDPServer) Start() {
	addr := net.UDPAddr{
		Port: s.port,
		IP:   net.ParseIP(s.ip),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer conn.Close()

	fmt.Printf("UDP server listening on %s:%d\n", s.ip, s.port)

	buffer := make([]byte, 1024)

	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}

		s.mu.Lock()
		addr_str := remoteAddr.String()
		if _, exists := s.connectedPlayers[addr_str]; !exists {
			fmt.Printf("New connection: %s\n", addr_str)
			s.connectedPlayers[addr_str] = player.Player{Address: addr_str}
		}
		s.mu.Unlock()

		fmt.Printf("Received %d bytes from %s: %s\n", n, remoteAddr, string(buffer[:n]))

		_, err = conn.WriteToUDP(buffer[:n], remoteAddr)
		if err != nil {
			fmt.Println("Error writing to UDP:", err)
		}
	}
}

func (s *UDPServer) GetConnectedPlayers() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	players := make([]string, 0, len(s.connectedPlayers))
	for _, player := range s.connectedPlayers {
		players = append(players, player.Address)
	}
	return players
}
