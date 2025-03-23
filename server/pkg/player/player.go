package player

import (
	"github.com/google/uuid"
	"net"
)

type Player struct {
	Address *net.UDPAddr
	id      string
}

func New(addr *net.UDPAddr) Player {
	player := Player{
		Address: addr,
		id:      generateUUID(),
	}
	return player
}

func (p *Player) Id() string {
	return p.id
}

func generateUUID() string {
	return uuid.New().String()
}
