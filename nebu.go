package resources

import (
	"encoding/binary"

	"github.com/imle/resourcefork"
)

// Nëbu resources contain info on the nebulae (or other space phenomena) which are displayed in the background of the
// star map. These images don't actually have any effect on events in the game, they're just there to look pretty.
// You can, however, combine nëbu background images with custom asteroid or interference data in the sÿst resources
// for cool localized effects, and you can also use the nëbu resource's OnExplore field to enable events when the
// player explores a certain nebula.
// The PICT resources associated with the 32 available nëbu resources are numbered 9500-9724. Each nebula can have up
// to seven PICT resources associated with it; the first nebula's PICTs are 9500-96, the second's are 9507-13, etc.
// The engine will pick the best nebula image to display for a given map scale (the map scales used in Nova are 42.1%,
// 56.2%, 75.0%, 100.0%, 133.3%, 177.7%, and 237.0%). If nebula image of the proper size for the current map scale
// isn't available, Nova will pick the closest one and stretch it as necessary. PICTs for a given nebula should be
// sorted by size in ascending order.

type NebuID IDType

type Nebu struct {
	ID NebuID

	XPos      int16
	YPos      int16
	XSize     int16
	YSize     int16
	ActiveOn  ControlBitTest
	OnExplore ControlBitFunction
}

func NebuFromResource(resource resourcefork.Resource) *Nebu {
	return NebuFromBytes(NebuID(resource.ID), resource.Data)
}

func NebuFromBytes(id NebuID, b []byte) *Nebu {
	t := &Nebu{
		ID:        id,
		XPos:      int16(binary.BigEndian.Uint16(b[0:])),
		YPos:      int16(binary.BigEndian.Uint16(b[2:])),
		XSize:     int16(binary.BigEndian.Uint16(b[4:])),
		YSize:     int16(binary.BigEndian.Uint16(b[6:])),
		ActiveOn:  ControlBitTest(byteString(b[5:], 254)),
		OnExplore: ControlBitFunction(byteString(b[263:], 255)),
	}

	return t
}
