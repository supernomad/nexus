package common

import (
	"syscall"
)

type Packet struct {
	Data               []byte
	Length             int
	DestinationAddress syscall.Sockaddr
	SourceAddress      syscall.Sockaddr
}

func NewPacket(buf []byte, src syscall.Sockaddr, size int) *Packet {
	return &Packet{
		Data:               buf,
		Length:             size,
		DestinationAddress: src,
		SourceAddress:      src,
	}
}
