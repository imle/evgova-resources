package resources

import (
	"encoding/binary"

	"github.com/imle/resourcefork"
)

// Junk resources store info on specialized commodities that can be bought and sold at a few locations.

type JunkID IDType

type JunkFlags struct {
	Tribbles   bool // 0x0001 Tribbles flag - When in your cargo bay, the commodity multiplies like tribbles.
	Perishable bool // 0x0002 Perishable - When in your cargo bay, the commodity gradually decays away.
}

type Junk struct {
	ID JunkID

	SoldAt    [8]SpobID      // ID number of the stellar object where the commodity is sold. Set to 0 or -1 if unused.
	BoughtAt  [8]SpobID      // ID number of the stellar object where the commodity is purchased. Set to 0 or -1 if unused.
	BasePrice int16          // The average price of the commodity (works much like the base prices for "regular" commodities).
	Flags     JunkFlags      // Misc flag bits.
	ScanMask  FlagMask16     // Tribbles flag - When in your cargo bay, the commodity multiplies like tribbles.
	LCName    string         // The lower-case string to display in the player-info dialog box, among other places, e.g. "machine parts".
	Abbrev    string         // The short string that is displayed in the player's status bar when the player is carrying jünk of this type, e.g. "Parts".
	BuyOn     ControlBitTest // This jünk will only be available to be bought when this expression evaluates true. Leave blank if unused.
	SellOn    ControlBitTest // This jünk will only be able to be sold when this expression evaluates true. Leave blank if unused.
}

func JunkFromResource(resource resourcefork.Resource) *Junk {
	return JunkFromBytes(JunkID(resource.ID), resource.Data)
}

func JunkFromBytes(id JunkID, b []byte) *Junk {
	flags1 := binary.BigEndian.Uint16(b[34:])

	t := &Junk{
		ID: id,
		SoldAt: [8]SpobID{
			SpobID(binary.BigEndian.Uint16(b[0:])),
			SpobID(binary.BigEndian.Uint16(b[2:])),
			SpobID(binary.BigEndian.Uint16(b[4:])),
			SpobID(binary.BigEndian.Uint16(b[6:])),
			SpobID(binary.BigEndian.Uint16(b[8:])),
			SpobID(binary.BigEndian.Uint16(b[10:])),
			SpobID(binary.BigEndian.Uint16(b[12:])),
			SpobID(binary.BigEndian.Uint16(b[14:])),
		},
		BoughtAt: [8]SpobID{
			SpobID(binary.BigEndian.Uint16(b[16:])),
			SpobID(binary.BigEndian.Uint16(b[18:])),
			SpobID(binary.BigEndian.Uint16(b[20:])),
			SpobID(binary.BigEndian.Uint16(b[22:])),
			SpobID(binary.BigEndian.Uint16(b[24:])),
			SpobID(binary.BigEndian.Uint16(b[26:])),
			SpobID(binary.BigEndian.Uint16(b[28:])),
			SpobID(binary.BigEndian.Uint16(b[30:])),
		},
		BasePrice: int16(binary.BigEndian.Uint16(b[32:])),
		Flags: JunkFlags{
			Tribbles:   flags1&0x0001 == 0x0001,
			Perishable: flags1&0x0002 == 0x0002,
		},
		ScanMask: FlagMask16(binary.BigEndian.Uint16(b[36:])),
		LCName:   byteString(b[38:], 62),
		Abbrev:   byteString(b[102:], 63),
		BuyOn:    ControlBitTest(byteString(b[421:], 254)),
		SellOn:   ControlBitTest(byteString(b[167:], 254)),
	}

	return t
}
