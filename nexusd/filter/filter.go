package filter

import (
	"github.com/Supernomad/nexus/nexusd/common"
)

type Filter interface {
	Drop(packet *common.Packet) bool
}
