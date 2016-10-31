package common

import (
	"syscall"
)

// Packet object to represent a packet traversing the system
type Packet struct {
	Data               []byte
	Length             int
	DestinationAddress syscall.Sockaddr
	SourceAddress      syscall.Sockaddr
}

// NewPacket object
func NewPacket(buf []byte, src syscall.Sockaddr, size int) *Packet {
	return &Packet{
		Data:               buf,
		Length:             size,
		DestinationAddress: src,
		SourceAddress:      src,
	}
}
