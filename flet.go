package resources

import (
	"encoding/binary"

	"github.com/imle/resourcefork"
)

// A flet resource defines the parameters for a fleet, which is a collection of ships that can be made to
// appear randomly throughout the galaxy.

type FletID IDType

type FletFlags struct {
	FreightersHaveRandomCargo bool // 0x0001 Freighters (InherentAI <= 2) in this fleet will have random cargo when boarded.
}

type LinkSystID SystID

type Flet struct {
	ID FletID

	LeadShipType ShipID    // ID of the fleet's flagship's ship class.
	EscortType   [4]ShipID // IDs of the flagships escorts' ship classes. If you don't want to use four different escort types, you should still set the unused fields to a valid ship class ID. (you can set the min & max fields to 0 and just have the extra ships not appear).
	Min          [4]int16  // The minimum number of each type of escort to put in the fleet.
	Max          [4]int16  // The maximum number of each type of escort to put in the fleet.
	Govt         GovtID    // ID of the fleet's government, of -1 for none.

	// -1          Any system.
	// 128-2175    ID of a specific system.
	// 10000-10255 Any system belonging to this specific government.
	// 15000-15255 Any system belonging to an ally of this govt.
	// 20000-20255 Any system belonging to any but this govt.
	// 25000-25255 Any system belonging to an enemy of this govt.
	LinkSyst LinkSystID // Which systems the fleet can be created in.

	AppearOn ControlBitTest // A control bit test field that will cause a given fleet to appear only when the expression evaluates to true. If this field is left blank it will be ignored.

	Quote StrAID // Show a random string from the STR# resource with this ID when the fleet enters from hyperspace. Any occurrences of the character '#' in this string will be replaced with a random digit (0-9).

	Flags FletFlags
}

func FletFromResource(resource resourcefork.Resource) *Flet {
	return FletFromBytes(FletID(resource.ID), resource.Data)
}

func FletFromBytes(id FletID, b []byte) *Flet {
	flags := binary.BigEndian.Uint16(b[288:])

	t := &Flet{
		ID:           id,
		LeadShipType: ShipID(binary.BigEndian.Uint16(b[0:])),
		EscortType: [4]ShipID{
			ShipID(binary.BigEndian.Uint16(b[2:])),
			ShipID(binary.BigEndian.Uint16(b[4:])),
			ShipID(binary.BigEndian.Uint16(b[6:])),
			ShipID(binary.BigEndian.Uint16(b[8:])),
		},
		Min: [4]int16{
			int16(binary.BigEndian.Uint16(b[10:])),
			int16(binary.BigEndian.Uint16(b[12:])),
			int16(binary.BigEndian.Uint16(b[14:])),
			int16(binary.BigEndian.Uint16(b[16:])),
		},
		Max: [4]int16{
			int16(binary.BigEndian.Uint16(b[18:])),
			int16(binary.BigEndian.Uint16(b[20:])),
			int16(binary.BigEndian.Uint16(b[22:])),
			int16(binary.BigEndian.Uint16(b[24:])),
		},
		Govt:     GovtID(binary.BigEndian.Uint16(b[26:])),
		LinkSyst: LinkSystID(binary.BigEndian.Uint16(b[28:])),
		AppearOn: ControlBitTest(byteString(b[30:], 254)),
		Quote:    StrAID(binary.BigEndian.Uint16(b[286:])),
		Flags: FletFlags{
			FreightersHaveRandomCargo: flags&0x0001 == 0x0001,
		},
	}

	return t
}
