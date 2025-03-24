package player

import (
	"net"
	"sync/atomic"
	"time"
)

type Player struct {
	Address  *net.UDPAddr
	id       uint32
	LastSeen time.Time
}

var idCounter uint32

func New(addr *net.UDPAddr) Player {
	player := Player{
		Address:  addr,
		id:       generateID(),
		LastSeen: time.Now(),
	}
	return player
}

func (p *Player) Id() uint32 {
	return p.id
}

func generateID() uint32 {
	return atomic.AddUint32(&idCounter, 1)
}

func (p *Player) UpdateLastSeen() {
	p.LastSeen = time.Now()
}
