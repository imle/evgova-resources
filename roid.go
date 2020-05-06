package resources

import (
	"encoding/binary"
	"image/color"

	"github.com/imle/resourcefork"
)

// Nova supports up to 16 asteroid types, each of which can have its own special properties. röid resources 128-143
// store the attributes for each asteroid type.

type RoidID IDType

type Roid struct {
	ID RoidID

	Strength    int16
	SpinRate    int16
	YieldType   int16
	YieldQty    int16
	PartCount   int16
	PartColor   color.Color
	FragType1   RoidID
	FragType2   RoidID
	FragCount   int16
	ExplodeType ExplodeType
	Mass        int16
}

func RoidFromResource(resource resourcefork.Resource) *Roid {
	return RoidFromBytes(RoidID(resource.ID), resource.Data)
}

func RoidFromBytes(id RoidID, b []byte) *Roid {
	t := &Roid{
		ID:        id,
		Strength:  int16(binary.BigEndian.Uint16(b[0:])),
		SpinRate:  int16(binary.BigEndian.Uint16(b[2:])),
		YieldType: int16(binary.BigEndian.Uint16(b[4:])),
		YieldQty:  int16(binary.BigEndian.Uint16(b[6:])),
		PartCount: int16(binary.BigEndian.Uint16(b[8:])),
		PartColor: color.RGBA{
			A: b[10],
			R: b[11],
			G: b[12],
			B: b[13],
		},
		FragType1:   RoidID(binary.BigEndian.Uint16(b[14:])),
		FragType2:   RoidID(binary.BigEndian.Uint16(b[16:])),
		FragCount:   int16(binary.BigEndian.Uint16(b[18:])),
		ExplodeType: ExplodeType(binary.BigEndian.Uint16(b[20:])),
		Mass:        int16(binary.BigEndian.Uint16(b[22:])),
	}

	return t
}
