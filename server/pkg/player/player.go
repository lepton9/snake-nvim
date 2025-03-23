package player

import (
	"github.com/google/uuid"
	"net"
	"time"
)

type Player struct {
	Address  *net.UDPAddr
	id       string
	LastSeen time.Time
}

func New(addr *net.UDPAddr) Player {
	player := Player{
		Address:  addr,
		id:       generateUUID(),
		LastSeen: time.Now(),
	}
	return player
}

func (p *Player) Id() string {
	return p.id
}

func generateUUID() string {
	return uuid.New().String()
}

func (p *Player) UpdateLastSeen() {
	p.LastSeen = time.Now()
}
