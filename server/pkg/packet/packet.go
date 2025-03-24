package packet

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type PacketType uint8

const (
	JOIN PacketType = iota
	MOVE
	LEAVE
	UPDATE
)

type Packet struct {
	Len      int
	PlayerID uint32     // 4 bytes
	Type     PacketType // 1 byte
	Data     []byte     // Variable length
}

func EncodePacket(pkt Packet) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, pkt.PlayerID)
	binary.Write(buf, binary.LittleEndian, pkt.Type)
	buf.Write(pkt.Data)
	return buf.Bytes()
}

func DecodePacket(data []byte) (Packet, error) {
	var pkt Packet
	pkt.Len = len(data)
	if pkt.Len < 5 {
		return pkt, fmt.Errorf("packet too short")
	}

	pkt.PlayerID = binary.LittleEndian.Uint32(data[:4])
	pkt.Type = PacketType(data[4])
	pkt.Data = data[5:]
	return pkt, nil
}
