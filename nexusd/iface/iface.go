package iface

import (
	"github.com/Supernomad/nexus/nexusd/common"
)

const (
	// UDPSocket iface
	UDPSocket int = 0
	// AFPacket iface
	AFPacket int = 1
	// NetMap iface
	NetMap int = 2
)

// Iface is a generic multi-queue networking interface
type Iface interface {
	Read(queue int) (*common.Packet, error)
	Write(queue int, payload *common.Packet) error
	Open() error
	Close() error
	GetFDs() []int
}

// New Iface object
func New(kind int, log *common.Logger, cfg *common.Config) Iface {
	switch kind {
	case UDPSocket:
		return newSocket(log, cfg)
		/*
			case AFPacket:
				return newAFPacket(log, cfg)
			case NetMap:
				return newNetMap(log, cfg)
		*/
	}
	return nil
}
