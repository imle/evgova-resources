package resources

import (
	"encoding/binary"

	"github.com/imle/resourcefork"
)

// A dude resource can be thought of as a container for ships that share certain characteristics.
// Any ship of a given dude class will have that dude class's AI type and governmental affiliation,
// and will yield the same types of booty when boarded. In a dude resource, up to 16 different ship
// classes can be pointed to, with a probability set for each ship class. The result of all this is
// that, in other parts of Nova's data file, you can point to a dude class and know that Nova will
// create a ship of the proper AI type and governmental alignment, and will pick the new ship's type
// based on the probabilities you set in the dude resource.

type DudeID IDType

type DudeFlags struct {
	CarriesFood            bool // 0x0001 Carries food when plundered.
	CarriesIndustrialGoods bool // 0x0002 Carries industrial goods.
	CarriesMedicalSupplies bool // 0x0004 Carries medical supplies.
	CarriesLuxuryGoods     bool // 0x0008 Carries luxury goods.
	CarriesMetal           bool // 0x0010 Carries metal.
	CarriesEquipment       bool // 0x0020 Carries equipment.
	CarriesMoney           bool // 0x0040 Carries money (amount depends on the ship's purchase price).
	NoHitBoxForPlayer      bool // 0x0100 Ships of this dude type can't be hit by the player and their shots can't hit the player (useful for things like AuxShip mission escorts, etc.).
}

type DudeInfoTypes struct {
	GoodsPrices    bool    // 0x1000 Good prices.
	DisasterInfo   bool    // 0x2000 Disaster info.
	SpecificAdvice *StrAID // 0x4xxx Specific advice (the lower 12 bits of this value are added to 7500 to get the ID of the STR# resource from which to get the quote).
	GenericHail    bool    // 0x8000 Generic govt hail messages.
}

// You can set different combinations of booty to be had from ships of a certain dude class by ORing different bits
// into the dude's Booty field. If a dude class has a booty flag of 0x0000, then you can't get anything from the ship,
// and you're told that you were "repelled while attempting to board" it. The different booty flags are documented above

type Dude struct {
	ID DudeID

	AIType      AIType     // Which type of AI to use for ships of this dude class (see below). If you set this to 0, each ship will use its own inherent AI type.
	Govt        GovtID     // The ID number of the dude class's government, or -1 for independent.
	Booty       FlagMask32 // Flags that define what you'll get when you board a ship of this dude class.
	Flags       DudeFlags
	InfoTypes   DudeInfoTypes // What kind of info to display when hailed.
	ShipType    [16]ShipID    // These fields contain the ID numbers of up to 16 different ship classes. Set to 0 or -1 if unused.
	Probability [16]int16     // These fields set the probability that a ship of this dude class will be of a certain ship type.
}

func DudeFromResource(resource resourcefork.Resource) *Dude {
	return DudeFromBytes(DudeID(resource.ID), resource.Data)
}

func DudeFromBytes(id DudeID, b []byte) *Dude {
	flags := binary.BigEndian.Uint16(b[4:])
	flags2 := binary.BigEndian.Uint16(b[4:])

	t := &Dude{
		ID:     id,
		AIType: AIType(binary.BigEndian.Uint16(b[0:])),
		Govt:   GovtID(binary.BigEndian.Uint16(b[2:])),
		Flags: DudeFlags{
			CarriesFood:            flags&0x0001 == 0x0001,
			CarriesIndustrialGoods: flags&0x0002 == 0x0002,
			CarriesMedicalSupplies: flags&0x0004 == 0x0004,
			CarriesLuxuryGoods:     flags&0x0008 == 0x0008,
			CarriesMetal:           flags&0x0010 == 0x0010,
			CarriesEquipment:       flags&0x0020 == 0x0020,
			CarriesMoney:           flags&0x0040 == 0x0040,
			NoHitBoxForPlayer:      flags&0x0100 == 0x0100,
		},
		InfoTypes: DudeInfoTypes{
			GoodsPrices:  flags2&0x1000 == 0x1000,
			DisasterInfo: flags2&0x2000 == 0x2000,
			GenericHail:  flags2&0x8000 == 0x8000,
		},
		ShipType: [16]ShipID{
			ShipID(binary.BigEndian.Uint16(b[8:])),
			ShipID(binary.BigEndian.Uint16(b[10:])),
			ShipID(binary.BigEndian.Uint16(b[12:])),
			ShipID(binary.BigEndian.Uint16(b[14:])),
			ShipID(binary.BigEndian.Uint16(b[16:])),
			ShipID(binary.BigEndian.Uint16(b[18:])),
			ShipID(binary.BigEndian.Uint16(b[20:])),
			ShipID(binary.BigEndian.Uint16(b[22:])),
			ShipID(binary.BigEndian.Uint16(b[24:])),
			ShipID(binary.BigEndian.Uint16(b[26:])),
			ShipID(binary.BigEndian.Uint16(b[28:])),
			ShipID(binary.BigEndian.Uint16(b[30:])),
			ShipID(binary.BigEndian.Uint16(b[32:])),
			ShipID(binary.BigEndian.Uint16(b[34:])),
			ShipID(binary.BigEndian.Uint16(b[36:])),
			ShipID(binary.BigEndian.Uint16(b[38:])),
		},
		Probability: [16]int16{
			int16(binary.BigEndian.Uint16(b[40:])),
			int16(binary.BigEndian.Uint16(b[42:])),
			int16(binary.BigEndian.Uint16(b[44:])),
			int16(binary.BigEndian.Uint16(b[46:])),
			int16(binary.BigEndian.Uint16(b[48:])),
			int16(binary.BigEndian.Uint16(b[50:])),
			int16(binary.BigEndian.Uint16(b[52:])),
			int16(binary.BigEndian.Uint16(b[54:])),
			int16(binary.BigEndian.Uint16(b[56:])),
			int16(binary.BigEndian.Uint16(b[58:])),
			int16(binary.BigEndian.Uint16(b[60:])),
			int16(binary.BigEndian.Uint16(b[62:])),
			int16(binary.BigEndian.Uint16(b[64:])),
			int16(binary.BigEndian.Uint16(b[66:])),
			int16(binary.BigEndian.Uint16(b[68:])),
			int16(binary.BigEndian.Uint16(b[70:])),
		},
	}

	// See DudeInfoTypes.SpecificAdvice comment
	if flags2&0x4000 == 0x4000 {
		a := StrAID(flags2&0x0FFF) + 7500
		t.InfoTypes.SpecificAdvice = &a
	}

	return t
}
