package resources

import (
	"encoding/binary"

	"github.com/imle/resourcefork"
)

// Oops resources contain info on planetary disasters. Actually, the term 'disasters' is a misnomer, as these
// occurrences simply affect the price of a single commodity at a planet or station, for good or bad. Nova uses
// the name of the resource in the commodity exchange dialog box to indicate that a disaster is currently going
// on at a planet.

type OopsID IDType

type Oops struct {
	ID OopsID

	Stellar    SpobID
	Commodity  CommodityType
	PriceDelta int16
	Duration   int16
	Freq       int16
	ActivateOn ControlBitTest
}

func OopsFromResource(resource resourcefork.Resource) *Oops {
	return OopsFromBytes(OopsID(resource.ID), resource.Data)
}

func OopsFromBytes(id OopsID, b []byte) *Oops {
	t := &Oops{
		ID:         id,
		Stellar:    SpobID(binary.BigEndian.Uint16(b[0:])),
		Commodity:  CommodityType(int16(binary.BigEndian.Uint16(b[2:]))),
		PriceDelta: int16(binary.BigEndian.Uint16(b[4:])),
		Duration:   int16(binary.BigEndian.Uint16(b[6:])),
		Freq:       int16(binary.BigEndian.Uint16(b[8:])),
		ActivateOn: ControlBitTest(byteString(b[10:], 254)),
	}

	return t
}
