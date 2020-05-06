package resources

import (
	"encoding/binary"

	"github.com/imle/resourcefork"
)

// The spïn resource
// Spin resources contain sprite info for simple graphical objects. Whenever Nova needs to load a set
// of sprites for a particular object, it looks at that object's spin resource, which in turn tells the game
// how to load the object's sprites. Nova sprites are stored as paired sprite and mask PICT resources,
// or as rleD/rle8 resources. The sprites in each PICT are arranged in a grid, which can be of any size.
// The spin resource tells Nova what shape and size the sprites' grid is.

// Spin resources have certain reserved ID numbers, which correspond to different types of objects:
// 400-463     Explosions.
// 500         Cargo boxes.
// 501-504     Mini-asteroids for mining. 600-605 Main menu buttons.
// 606         Main screen logo.
// 607         Main screen rollover images.
// 608-610     Main screen sliding buttons. 650 Target cursor.
// 700         Starfield.
// 800-815     Asteroids.
// 1000-1255   Stellar objects. 3000-3255 Weapons.

const (
	SpinIDCargoBoxes               SpinID = 500
	SpinIDMainScreenLogo           SpinID = 606
	SpinIDMainScreenRolloverImages SpinID = 607
	SpinIDStarField                SpinID = 700
)

// It is important to note that the ID numbers of the PICT/rleD/rle8 resources are non-critical, as Nova
// looks at the spin resources to find the sprites, and not at the actual PICT/rleD/rle8 ID numbers themselves.

type SpinID IDType
type GraphicID IDType

func (g GraphicID) PictID() PictID {
	return PictID(g)
}

func (g GraphicID) RleDID() RleDID {
	return RleDID(g)
}

type Spin struct {
	ID        SpinID
	SpritesID GraphicID // ID number of the sprites' PICT resource (or the ID of the rleD/rle8 resource).
	MasksID   PictID    // ID number of the masks' PICT resource.
	xSize     int16     // Horizontal size of each sprite.
	ySize     int16     // Vertical size of each sprite.
	xTiles    int16     // Horizontal grid dimension.
	yTiles    int16     // Vertical grid dimension.
}

func SpinFromResource(resource resourcefork.Resource) *Spin {
	return SpinFromBytes(SpinID(resource.ID), resource.Data)
}

func SpinFromBytes(id SpinID, b []byte) *Spin {
	t := &Spin{
		ID:        id,
		SpritesID: GraphicID(binary.BigEndian.Uint16(b[0:])),
		MasksID:   PictID(binary.BigEndian.Uint16(b[2:])),
		xSize:     int16(binary.BigEndian.Uint16(b[4:])),
		ySize:     int16(binary.BigEndian.Uint16(b[6:])),
		xTiles:    int16(binary.BigEndian.Uint16(b[8:])),
		yTiles:    int16(binary.BigEndian.Uint16(b[10:])),
	}

	return t
}
