package resources

import (
	"encoding/binary"
	"time"

	"github.com/imle/resourcefork"
)

// Spob resources describe stellar objects, such as planets and space stations. (spob stands for space object)
// Each spob resource represents one stellar object, whose name is the name as the name of the resource.

type SpobID IDType

type PlanetGraphic int16

func (p PlanetGraphic) RleDID() RleDID {
	return RleDID(p) + RleDIDOffsetSpob
}

type FleetInfo struct {
	value int16
}

func (f FleetInfo) Get() (ships int16, waves int16) {
	if f.value >= 1000 {
		return f.value % 10, f.value/10 - 100
	} else {
		return f.value, 1
	}
}

type SpobFlags struct {
	CanLand               bool // 0x00000001 Can land/dock here.
	HasCommodityExchange  bool // 0x00000002 Has commodity exchange.
	HasOutfitter          bool // 0x00000004 Can outfit ship here.
	HasShipyard           bool // 0x00000008 Can buy ships here.
	IsStation             bool // 0x00000010 Stellar is a station instead of a planet.
	Uninhabited           bool // 0x00000020 Stellar is uninhabited (no traffic control or refuelling).
	HasBar                bool // 0x00000040 Has bar.
	DestroyToLand         bool // 0x00000080 Can only land here if stellar is destroyed first.
	NoTradeFood           bool // 0x00000000 Won't trade in food.
	LowPriceFood          bool // 0x10000000 Low food prices.
	MediumPriceFood       bool // 0x20000000 Medium food prices.
	HighPriceFood         bool // 0x40000000 High food prices.
	NoTradeIndustrial     bool // 0x00000000 Won't trade in industrial goods.
	LowPriceIndustrial    bool // 0x01000000 Low industrial prices.
	MediumPriceIndustrial bool // 0x02000000 Medium industrial prices.
	HighPriceIndustrial   bool // 0x04000000 High industrial prices.
	NoTradeMedical        bool // 0x00000000 Won't trade in medical supplies.
	LowPriceMedical       bool // 0x00100000 Low medical prices.
	MediumPriceMedical    bool // 0x00200000 Medium medical prices.
	HighPriceMedical      bool // 0x00400000 High medical prices.
	NoTradeLuxury         bool // 0x00000000 Won't trade in luxury goods.
	LowPriceLuxury        bool // 0x00010000 Low luxury prices.
	MediumPriceLuxury     bool // 0x00020000 Medium luxury prices.
	HighPriceLuxury       bool // 0x00040000 High luxury prices.
	NoTradeMetal          bool // 0x00000000 Won't trade in metal.
	LowPriceMetal         bool // 0x00001000 Low metal prices.
	MediumPriceMetal      bool // 0x00002000 Medium metal prices.
	HighPriceMetal        bool // 0x00004000 High metal prices.
	NoTradeEquipment      bool // 0x00000000 Won't trade in equipment.
	LowPriceEquipment     bool // 0x00000100 Low equipment prices.
	MediumPriceEquipment  bool // 0x00000200 Medium equipment prices.
	HighPriceEquipment    bool // 0x00000400 High equipment prices.

	AnimatedFirstFrameEveryOtherFrame bool // 0x0001 For an animated stellar, the first frame will be shown after each subsequent frame.
	AnimatedRandomFrames              bool // 0x0002 For an animated stellar, the next frame in the sequence will be picked at random. The same frame will not be picked twice in a row. Note that this can be combined with the previous flag and the Frame0Bias field to create interesting effects such as random planetary lightning or lights twinkling.
	SoundLoop                         bool // 0x0010 Play this stellar's sound in a continuous loop.
	AllYourBaseAreBelongToUs          bool // 0x0020 Stellar is always dominated (all your base are belong to us).
	StartDestroyed                    bool // 0x0040 Stellar starts the game destroyed.
	AnimatedWhenDestroyed             bool // 0x0080 For an animated stellar, the stellar's graphic is animated after it's been destroyed and static when it is not destroyed. The normal behaviour is the opposite of this: static when destroyed and animated when not.
	DeadlyStellar                     bool // 0x0100 Stellar is deadly - all ships that touch it are destroyed immediately.
	OnlyAttackWhenProvoked            bool // 0x0200 If the stellar has a weapon, it will only fire when provoked (i.e. only when the player is trying to dominate it).
	BuybackAnyTechLevel               bool // 0x0400 If the stellar has an outfit shop, it can buy any nonpermanent outfits the player owns, regardless of tech level.
	IsHyperGate                       bool // 0x1000 Stellar is a hypergate - if the player lands on it he will be given a choice of which other hypergate to travel to (see HyperLink1-8 below).
	IsWormHole                        bool // 0x2000 Stellar is a wormhole - if the player lands on it he will be transported to some other random somewhere in the galaxy. If all of the wormhole's HyperLink1-8 fields set to -1, the player will end up at another random wormhole which also has no defined hyper links. If the wormhole has any hyper links defined, the player will end up at one of the wormholes on the other end.
}

type Spob struct {
	ID SpobID

	xPos            int16
	yPos            int16
	Type            PlanetGraphic
	Flags           SpobFlags
	Tribute         Credits
	TechLevel       TechLevel
	SpecialTech     [8]TechLevel
	Govt            GovtID
	MinStatus       Status
	CustPicID       PictID
	CustSndID       SndID
	DefenseDude     DudeID
	DefCount        FleetInfo
	AnimDelay       time.Duration
	Frame0Bias      int16
	HyperLink       [8]SpobID
	TransitionFrame int16
	ExitAngle       int16
	OnDominate      ControlBitFunction
	OnRelease       ControlBitFunction
	Fee             Credits
	Gravity         int16
	Weapon          WeapID
	Strength        HitPoints
	DeadType        PlanetGraphic
	DeadTime        int16
	ExplodeType     ExplodeType
	OnDestroy       ControlBitFunction
	OnRegen         ControlBitFunction
}

func (m Spob) DescID() DescID {
	return DescID(m.ID) - resourcefork.ResourceForkIDOffset + DescIDOffsetStellar
}

func SpobFromResource(resource resourcefork.Resource) *Spob {
	return SpobFromBytes(SpobID(resource.ID), resource.Data)
}

func SpobFromBytes(id SpobID, b []byte) *Spob {
	flags := binary.BigEndian.Uint32(b[6:])
	flags2 := binary.BigEndian.Uint16(b[32:])

	t := &Spob{
		ID:   id,
		xPos: int16(binary.BigEndian.Uint16(b[0:])),
		yPos: int16(binary.BigEndian.Uint16(b[2:])),
		Type: PlanetGraphic(binary.BigEndian.Uint16(b[4:])),
		Flags: SpobFlags{
			CanLand:                           flags&0x00000001 == 0x00000001,
			HasCommodityExchange:              flags&0x00000002 == 0x00000002,
			HasOutfitter:                      flags&0x00000004 == 0x00000004,
			HasShipyard:                       flags&0x00000008 == 0x00000008,
			IsStation:                         flags&0x00000010 == 0x00000010,
			Uninhabited:                       flags&0x00000020 == 0x00000020,
			HasBar:                            flags&0x00000040 == 0x00000040,
			DestroyToLand:                     flags&0x00000080 == 0x00000080,
			NoTradeFood:                       flags&0x70000000 == 0x00000000, // zeroed is true
			LowPriceFood:                      flags&0x10000000 == 0x10000000,
			MediumPriceFood:                   flags&0x20000000 == 0x20000000,
			HighPriceFood:                     flags&0x40000000 == 0x40000000,
			NoTradeIndustrial:                 flags&0x07000000 == 0x00000000, // zeroed is true
			LowPriceIndustrial:                flags&0x01000000 == 0x01000000,
			MediumPriceIndustrial:             flags&0x02000000 == 0x02000000,
			HighPriceIndustrial:               flags&0x04000000 == 0x04000000,
			NoTradeMedical:                    flags&0x00700000 == 0x00000000, // zeroed is true
			LowPriceMedical:                   flags&0x00100000 == 0x00100000,
			MediumPriceMedical:                flags&0x00200000 == 0x00200000,
			HighPriceMedical:                  flags&0x00400000 == 0x00400000,
			NoTradeLuxury:                     flags&0x00070000 == 0x00000000, // zeroed is true
			LowPriceLuxury:                    flags&0x00010000 == 0x00010000,
			MediumPriceLuxury:                 flags&0x00020000 == 0x00020000,
			HighPriceLuxury:                   flags&0x00040000 == 0x00040000,
			NoTradeMetal:                      flags&0x00007000 == 0x00000000, // zeroed is true
			LowPriceMetal:                     flags&0x00001000 == 0x00001000,
			MediumPriceMetal:                  flags&0x00002000 == 0x00002000,
			HighPriceMetal:                    flags&0x00004000 == 0x00004000,
			NoTradeEquipment:                  flags&0x00000700 == 0x00000000, // zeroed is true
			LowPriceEquipment:                 flags&0x00000100 == 0x00000100,
			MediumPriceEquipment:              flags&0x00000200 == 0x00000200,
			HighPriceEquipment:                flags&0x00000400 == 0x00000400,
			AnimatedFirstFrameEveryOtherFrame: flags2&0x0001 == 0x0001,
			AnimatedRandomFrames:              flags2&0x0002 == 0x0002,
			SoundLoop:                         flags2&0x0010 == 0x0010,
			AllYourBaseAreBelongToUs:          flags2&0x0020 == 0x0020,
			StartDestroyed:                    flags2&0x0040 == 0x0040,
			AnimatedWhenDestroyed:             flags2&0x0080 == 0x0080,
			DeadlyStellar:                     flags2&0x0100 == 0x0100,
			OnlyAttackWhenProvoked:            flags2&0x0200 == 0x0200,
			BuybackAnyTechLevel:               flags2&0x0400 == 0x0400,
			IsHyperGate:                       flags2&0x1000 == 0x1000,
			IsWormHole:                        flags2&0x2000 == 0x2000,
		},
		Tribute:   Credits(int16(binary.BigEndian.Uint16(b[10:]))), // Is credits, but not 32 bit
		TechLevel: TechLevel(binary.BigEndian.Uint16(b[12:])),
		SpecialTech: [8]TechLevel{
			TechLevel(binary.BigEndian.Uint16(b[14:])),
			TechLevel(binary.BigEndian.Uint16(b[16:])),
			TechLevel(binary.BigEndian.Uint16(b[18:])),
			TechLevel(binary.BigEndian.Uint16(b[1092:])),
			TechLevel(binary.BigEndian.Uint16(b[1094:])),
			TechLevel(binary.BigEndian.Uint16(b[1096:])),
			TechLevel(binary.BigEndian.Uint16(b[1098:])),
			TechLevel(binary.BigEndian.Uint16(b[1100:])),
		},
		Govt:        GovtID(binary.BigEndian.Uint16(b[20:])),
		MinStatus:   Status(binary.BigEndian.Uint16(b[22:])),
		CustPicID:   PictID(binary.BigEndian.Uint16(b[24:])),
		CustSndID:   SndID(binary.BigEndian.Uint16(b[26:])),
		DefenseDude: DudeID(binary.BigEndian.Uint16(b[28:])),
		DefCount:    FleetInfo{value: int16(binary.BigEndian.Uint16(b[30:]))},
		AnimDelay:   time.Duration(binary.BigEndian.Uint16(b[34:])) * time.Second / 30,
		Frame0Bias:  int16(binary.BigEndian.Uint16(b[36:])),
		HyperLink: [8]SpobID{
			SpobID(binary.BigEndian.Uint16(b[38:])),
			SpobID(binary.BigEndian.Uint16(b[40:])),
			SpobID(binary.BigEndian.Uint16(b[42:])),
			SpobID(binary.BigEndian.Uint16(b[44:])),
			SpobID(binary.BigEndian.Uint16(b[46:])),
			SpobID(binary.BigEndian.Uint16(b[48:])),
			SpobID(binary.BigEndian.Uint16(b[50:])),
			SpobID(binary.BigEndian.Uint16(b[52:])),
		},
		TransitionFrame: int16(binary.BigEndian.Uint16(b[24:])),
		ExitAngle:       int16(binary.BigEndian.Uint16(b[26:])),
		OnDominate:      ControlBitFunction(byteString(b[54:], 255)),
		OnRelease:       ControlBitFunction(byteString(b[309:], 255)),
		Fee:             Credits(binary.BigEndian.Uint32(b[564:])),
		Gravity:         int16(binary.BigEndian.Uint16(b[568:])),
		Weapon:          WeapID(binary.BigEndian.Uint16(b[570:])),
		Strength:        HitPoints(binary.BigEndian.Uint32(b[572:])),
		DeadType:        PlanetGraphic(binary.BigEndian.Uint16(b[576:])),
		DeadTime:        int16(binary.BigEndian.Uint16(b[578:])),
		ExplodeType:     ExplodeType(binary.BigEndian.Uint16(b[580:])),
		OnDestroy:       ControlBitFunction(byteString(b[582:], 255)),
		OnRegen:         ControlBitFunction(byteString(b[837:], 255)),
	}

	return t
}
