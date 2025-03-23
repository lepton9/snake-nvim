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
	conn             *net.UDPConn
	connectedPlayers map[string]player.Player
	mu               sync.Mutex
}

func Init(ip string, port int) *UDPServer {
	server := UDPServer{
		port:             port,
		ip:               ip,
		conn:             nil,
		connectedPlayers: make(map[string]player.Player),
	}
	return &server
}

func (s *UDPServer) Start() {
	addr := net.UDPAddr{
		Port: s.port,
		IP:   net.ParseIP(s.ip),
	}

	var err error
	s.conn, err = net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer s.conn.Close()

	fmt.Printf("UDP server listening on %s:%d\n", s.ip, s.port)

	buffer := make([]byte, 1024)

	for {
		n, remoteAddr, err := s.conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}

		fmt.Printf("Received %d bytes from %s: %s\n", n, remoteAddr, string(buffer[:n]))

		if !s.IsConnectedAddr(remoteAddr) {
			newPlayer := s.Connect(remoteAddr)
			fmt.Printf("New connection: %s, ID: %s\n", remoteAddr, newPlayer.Id())
			s.Send(remoteAddr, newPlayer.Id())
		} else {
			s.Send(remoteAddr, string(buffer[:n]))
		}
	}
}

func (s *UDPServer) Send(addr *net.UDPAddr, msg string) bool {
	_, err := s.conn.WriteToUDP([]byte(msg), addr)
	if err != nil {
		fmt.Println("Error writing to UDP:", err)
		return false
	}
	return true
}

func (s *UDPServer) IsConnectedAddr(addr *net.UDPAddr) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, player := range s.connectedPlayers {
		if player.Address.String() == addr.String() {
			return true
		}
	}
	return false
}

func (s *UDPServer) Connect(addr *net.UDPAddr) *player.Player {
	s.mu.Lock()
	defer s.mu.Unlock()
	newPlayer := player.New(addr)
	s.connectedPlayers[newPlayer.Id()] = newPlayer
	return &newPlayer
}
