package filter

import (
	"github.com/Supernomad/nexus/nexusd/common"
)

// Filter interface to use for creating generic filters
type Filter interface {
	// Drop should determine whether or not the current filter causes a drop action.
	Drop(packet *common.Packet) bool
}
