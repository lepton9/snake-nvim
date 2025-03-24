package server

import (
	"fmt"
	"net"
	"snake-nvim.lepton9/pkg/player"
	"strings"
	"sync"
	"time"
)

type UDPServer struct {
	port             int
	ip               string
	conn             *net.UDPConn
	connectedPlayers map[uint32]player.Player
	timeoutDuration  time.Duration
	mu               sync.Mutex
}

func Init(ip string, port int) *UDPServer {
	server := UDPServer{
		port:             port,
		ip:               ip,
		conn:             nil,
		connectedPlayers: make(map[uint32]player.Player),
		timeoutDuration:  60 * time.Second,
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

	go s.checkTimeouts()

	fmt.Printf("UDP server listening on %s:%d\n", s.ip, s.port)
	s.run()
}

func (s *UDPServer) run() {
	buffer := make([]byte, 1024)
	for {
		n, remoteAddr, err := s.conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}

		fmt.Printf("%s > %d bytes: %s\n", remoteAddr, n, string(buffer[:n]))

		// TODO: make a message format
		// parse message
		message := string(buffer[:n])
		parts := strings.SplitN(message, ",", 2)
		if len(parts) != 2 {
			fmt.Println("Invalid message format")
			continue
		}

		playerIDStr := parts[0]
		data := parts[1]

		// Convert playerIDStr to uint32
		var playerID uint32
		if playerIDStr != "" {
			_, err := fmt.Sscanf(playerIDStr, "%d", &playerID)
			if err != nil {
				fmt.Println("Invalid player ID format")
				continue
			}
		}

		// New connection
		if playerIDStr == "" && !s.IsConnectedAddr(remoteAddr) {
			newPlayer := s.Connect(remoteAddr)
			fmt.Printf("New connection: %s, ID: %d\n", remoteAddr, newPlayer.Id())
			s.Send(remoteAddr, fmt.Sprintf("%d", newPlayer.Id()))
		} else { // Old connection
			s.mu.Lock()
			player, exists := s.connectedPlayers[playerID]
			if exists {
				player.UpdateLastSeen()
				s.connectedPlayers[playerID] = player
			}
			s.mu.Unlock()

			if !exists || player.Address.String() != remoteAddr.String() {
				fmt.Println("Invalid player ID or address")
				continue
			}
			s.Send(remoteAddr, "Success: "+data)
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

func (s *UDPServer) GetPlayer(id uint32) *player.Player {
	s.mu.Lock()
	defer s.mu.Unlock()
	player, exists := s.connectedPlayers[id]
	if !exists {
		return nil
	}
	return &player
}

func (s *UDPServer) Connect(addr *net.UDPAddr) *player.Player {
	s.mu.Lock()
	defer s.mu.Unlock()
	newPlayer := player.New(addr)
	s.connectedPlayers[newPlayer.Id()] = newPlayer
	return &newPlayer
}

func (s *UDPServer) checkTimeouts() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		for id, player := range s.connectedPlayers {
			if time.Since(player.LastSeen) > s.timeoutDuration {
				fmt.Printf("Player %d timed out\n", id)
				s.Send(player.Address, "Connection timed out..")
				delete(s.connectedPlayers, id)
			}
		}
		s.mu.Unlock()
	}
}
