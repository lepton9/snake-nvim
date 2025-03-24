package server

import (
	"fmt"
	"net"
	"snake-nvim.lepton9/pkg/packet"
	"snake-nvim.lepton9/pkg/player"
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

		packetDecoded, err := packet.DecodePacket(buffer[:n])
		if err != nil {
			fmt.Println("Invalid packet format: ", err.Error())
			continue
		}

		s.handlePacket(&packetDecoded, remoteAddr)
	}
}

func (s *UDPServer) handlePacket(p *packet.Packet, addr *net.UDPAddr) error {
	if p.Type == packet.JOIN && !s.IsConnectedAddr(addr) {
		newPlayer := s.Connect(addr)
		newPlayer.Name = string(p.Data)
		fmt.Printf("New connection: %s, ID: %d, Name: %s\n", addr, newPlayer.Id(), newPlayer.Name)
		s.Send(addr, fmt.Sprintf("%d", newPlayer.Id()))
	} else {
		player := s.GetPlayer(p.PlayerID)
		if player == nil {
			return fmt.Errorf("No player connected with id: %d", p.PlayerID)
		}
		switch p.Type {
		case packet.LEAVE:
			leave := packet.MakePacket(packet.LEAVE, []byte("Disconnecting"))
			s.Send(addr, string(packet.EncodePacket(leave)))
			s.DisconnectPlayer(player.Id())
			break
		case packet.MOVE:
			s.Send(addr, "Success: "+string(p.Data))
			break
		case packet.PING:
			s.Send(addr, "Success: "+string(p.Data))
			break
		default:
			return fmt.Errorf("unknown packet type: %d", p.Type)
		}
	}
	return nil
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
	if exists {
		player.UpdateLastSeen()
		s.connectedPlayers[id] = player
		return &player
	}
	return nil
}

func (s *UDPServer) DisconnectPlayer(id uint32) {
	delete(s.connectedPlayers, id)
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
				s.DisconnectPlayer(id)
			}
		}
		s.mu.Unlock()
	}
}
