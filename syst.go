package resources

import (
	"encoding/binary"
	"image/color"

	"github.com/imle/resourcefork"
)

// Syst resources store information on star systems, in which all combat, trading, and spaceflight take place. Each
// system can be linked to up to 16 other systems, and the player can make hyperspace jumps back and forth between them.

type SystID IDType

type AsteroidFlags struct {
	SmallMetal    bool // 0x0001 Small metal     (röid ID 128)
	MediumMetal   bool // 0x0002 Medium metal    (röid ID 129)
	LargeMetal    bool // 0x0004 Large metal     (röid ID 130)
	HugeMetal     bool // 0x0008 Huge metal      (röid ID 131)
	SmallIce      bool // 0x0010 Small ice       (röid ID 132)
	MediumIce     bool // 0x0020 Medium ice      (röid ID 133)
	LargeIce      bool // 0x0040 Large ice       (röid ID 134)
	HugeIce       bool // 0x0080 Huge ice        (röid ID 135)
	SmallDust     bool // 0x0100 Small Dust      (röid ID 136)
	MediumDust    bool // 0x0200 Medium Dust     (röid ID 137)
	LargeDust     bool // 0x0400 Large Dust      (röid ID 138)
	HugeDust      bool // 0x0800 Huge Dust       (röid ID 139)
	SmallCrystal  bool // 0x1000 Small Crystal   (röid ID 140)
	MediumCrystal bool // 0x2000 Medium Crystal  (röid ID 141)
	LargeCrystal  bool // 0x4000 Large Crystal   (röid ID 142)
	HugeCrystal   bool // 0x8000 Huge Crystal    (röid ID 143)
}

type Syst struct {
	ID SystID

	xPos              int16
	yPos              int16
	Connection        [16]SystID
	NavDef            [16]SpobID
	DudeTypes         [8]DudeID
	Prob              [8]int16
	AvgShips          int16
	Govt              GovtID
	Message           StrAID
	Asteroids         int16
	Interference      int16
	Visibility        ControlBitTest
	BackgroundColor   color.NRGBA
	Murk              int16
	AstTypes          AsteroidFlags // Lookup this info?
	ReinforceFleet    FletID
	ReinforceTime     FrameCount // FramesPerSecond
	ReinforceInterval int16
	Persons           [8]PersID
}

func SystFromResource(resource resourcefork.Resource) *Syst {
	return SystFromBytes(SystID(resource.ID), resource.Data)
}

func SystFromBytes(id SystID, b []byte) *Syst {
	flags := binary.BigEndian.Uint16(b[148:])

	t := &Syst{
		ID:   id,
		xPos: int16(binary.BigEndian.Uint16(b[0:])),
		yPos: int16(binary.BigEndian.Uint16(b[2:])),
		Connection: [16]SystID{
			SystID(binary.BigEndian.Uint16(b[4:])),
			SystID(binary.BigEndian.Uint16(b[6:])),
			SystID(binary.BigEndian.Uint16(b[8:])),
			SystID(binary.BigEndian.Uint16(b[10:])),
			SystID(binary.BigEndian.Uint16(b[12:])),
			SystID(binary.BigEndian.Uint16(b[14:])),
			SystID(binary.BigEndian.Uint16(b[16:])),
			SystID(binary.BigEndian.Uint16(b[18:])),
			SystID(binary.BigEndian.Uint16(b[20:])),
			SystID(binary.BigEndian.Uint16(b[22:])),
			SystID(binary.BigEndian.Uint16(b[24:])),
			SystID(binary.BigEndian.Uint16(b[26:])),
			SystID(binary.BigEndian.Uint16(b[28:])),
			SystID(binary.BigEndian.Uint16(b[30:])),
			SystID(binary.BigEndian.Uint16(b[32:])),
			SystID(binary.BigEndian.Uint16(b[34:])),
		},
		NavDef: [16]SpobID{
			SpobID(binary.BigEndian.Uint16(b[36:])),
			SpobID(binary.BigEndian.Uint16(b[38:])),
			SpobID(binary.BigEndian.Uint16(b[40:])),
			SpobID(binary.BigEndian.Uint16(b[42:])),
			SpobID(binary.BigEndian.Uint16(b[44:])),
			SpobID(binary.BigEndian.Uint16(b[46:])),
			SpobID(binary.BigEndian.Uint16(b[48:])),
			SpobID(binary.BigEndian.Uint16(b[50:])),
			SpobID(binary.BigEndian.Uint16(b[52:])),
			SpobID(binary.BigEndian.Uint16(b[54:])),
			SpobID(binary.BigEndian.Uint16(b[56:])),
			SpobID(binary.BigEndian.Uint16(b[58:])),
			SpobID(binary.BigEndian.Uint16(b[60:])),
			SpobID(binary.BigEndian.Uint16(b[62:])),
			SpobID(binary.BigEndian.Uint16(b[64:])),
			SpobID(binary.BigEndian.Uint16(b[66:])),
		},
		DudeTypes: [8]DudeID{
			DudeID(binary.BigEndian.Uint16(b[68:])),
			DudeID(binary.BigEndian.Uint16(b[70:])),
			DudeID(binary.BigEndian.Uint16(b[72:])),
			DudeID(binary.BigEndian.Uint16(b[74:])),
			DudeID(binary.BigEndian.Uint16(b[76:])),
			DudeID(binary.BigEndian.Uint16(b[78:])),
			DudeID(binary.BigEndian.Uint16(b[80:])),
			DudeID(binary.BigEndian.Uint16(b[82:])),
		},
		Prob: [8]int16{
			int16(binary.BigEndian.Uint16(b[84:])),
			int16(binary.BigEndian.Uint16(b[86:])),
			int16(binary.BigEndian.Uint16(b[88:])),
			int16(binary.BigEndian.Uint16(b[90:])),
			int16(binary.BigEndian.Uint16(b[92:])),
			int16(binary.BigEndian.Uint16(b[94:])),
			int16(binary.BigEndian.Uint16(b[96:])),
			int16(binary.BigEndian.Uint16(b[98:])),
		},
		AvgShips:     int16(binary.BigEndian.Uint16(b[100:])),
		Govt:         GovtID(binary.BigEndian.Uint16(b[102:])),
		Message:      StrAID(binary.BigEndian.Uint16(b[104:])),
		Asteroids:    int16(binary.BigEndian.Uint16(b[106:])),
		Interference: int16(binary.BigEndian.Uint16(b[108:])),
		Persons: [8]PersID{
			PersID(binary.BigEndian.Uint16(b[110:])),
			PersID(binary.BigEndian.Uint16(b[112:])),
			PersID(binary.BigEndian.Uint16(b[114:])),
			PersID(binary.BigEndian.Uint16(b[116:])),
			PersID(binary.BigEndian.Uint16(b[118:])),
			PersID(binary.BigEndian.Uint16(b[120:])),
			PersID(binary.BigEndian.Uint16(b[122:])),
			PersID(binary.BigEndian.Uint16(b[124:])),
		},
		BackgroundColor: color.NRGBA{
			A: b[142],
			R: b[143],
			G: b[144],
			B: b[145],
		},
		Murk: int16(binary.BigEndian.Uint16(b[146:])),
		AstTypes: AsteroidFlags{
			SmallMetal:    flags&0x0001 == 0x0001,
			MediumMetal:   flags&0x0002 == 0x0002,
			LargeMetal:    flags&0x0004 == 0x0004,
			HugeMetal:     flags&0x0008 == 0x0008,
			SmallIce:      flags&0x0010 == 0x0010,
			MediumIce:     flags&0x0020 == 0x0020,
			LargeIce:      flags&0x0040 == 0x0040,
			HugeIce:       flags&0x0080 == 0x0080,
			SmallDust:     flags&0x0100 == 0x0100,
			MediumDust:    flags&0x0200 == 0x0200,
			LargeDust:     flags&0x0400 == 0x0400,
			HugeDust:      flags&0x0800 == 0x0800,
			SmallCrystal:  flags&0x1000 == 0x1000,
			MediumCrystal: flags&0x2000 == 0x2000,
			LargeCrystal:  flags&0x4000 == 0x4000,
			HugeCrystal:   flags&0x8000 == 0x8000,
		},
		Visibility:        ControlBitTest(byteString(b[150:], 254)),
		ReinforceFleet:    FletID(binary.BigEndian.Uint16(b[406:])),
		ReinforceTime:     FrameCount(binary.BigEndian.Uint16(b[408:])),
		ReinforceInterval: int16(binary.BigEndian.Uint16(b[410:])),
	}

	return t
}
