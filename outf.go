package resources

import (
	"encoding/binary"

	"github.com/imle/resourcefork"
)

// Outf resources store information on the items that you can buy when you choose 'Outfit Ship' at a planet or station.

type OutfID IDType

type OutfMod struct {
	typeVal OutfModType
	value   int16
}

func (o OutfMod) OutfModType() OutfModType {
	return o.typeVal
}

func (o OutfMod) OutfModValue() int16 {
	return o.value
}

type OutfModType int32

const (
	_                                 OutfModType = iota // 0
	OutfModTypeWeapon                                    // 1
	OutfModTypeCargoSpace                                // 2
	OutfModTypeAmmunition                                // 3
	OutfModTypeShieldCapacity                            // 4
	OutfModTypeShieldRechargeSpeed                       // 5
	OutfModTypeArmour                                    // 6
	OutfModTypeAccelerationBooster                       // 7
	OutfModTypeSpeedIncrease                             // 8
	OutfModTypeTurnRateChange                            // 9
	_                                                    // 10
	OutfModTypeEscapePod                                 // 11
	OutfModTypeFuelCapacity                              // 12
	OutfModTypeDensityScanner                            // 13
	OutfModTypeIFF                                       // 14
	OutfModTypeAfterburner                               // 15
	OutfModTypeMap                                       // 16
	OutfModTypeCloakingDevice                            // 17
	OutfModTypeFuelScoop                                 // 18
	OutfModTypeAutoRefueller                             // 19
	OutfModTypeAutoEject                                 // 20
	OutfModTypeCleanLegalRecord                          // 21
	OutfModTypeHyperspaceSpeed                           // 22
	OutfModTypeHyperspaceDistance                        // 23
	OutfModTypeInterferenceMod                           // 24
	OutfModTypeMarines                                   // 25
	_                                                    // 26
	OutfModTypeIncreaseMaximum                           // 27
	OutfModTypeMurkMod                                   // 28
	OutfModTypeFasterArmourRecharge                      // 29
	OutfModTypeCloakScanner                              // 30
	OutfModTypeMiningScoop                               // 31
	OutfModTypeMultiJump                                 // 32
	OutfModTypeJammingType1                              // 33
	OutfModTypeJammingType2                              // 34
	OutfModTypeJammingType3                              // 35
	OutfModTypeJammingType4                              // 36
	OutfModTypeFastJump                                  // 37
	OutfModTypeInertialDampener                          // 38
	OutfModTypeIonDissipater                             // 39
	OutfModTypeIonAbsorber                               // 40
	OutfModTypeGravityResistance                         // 41
	OutfModTypeResistDeadlyStellars                      // 42
	OutfModTypePaint                                     // 43
	OutfModTypeReinforcementInhibitor                    // 44
	OutfModTypeModMaxGuns                                // 45
	OutfModTypeModMaxTurrets                             // 46
	OutfModTypeBomb                                      // 47
	OutfModTypeIFFScrambler                              // 48
	OutfModTypeRepairSystem                              // 49
	OutfModTypeNonLethalBomb                             // 50
)

type OutfFlags struct {
	FixedGun                    bool // 0x0001 This item is a fixed gun.
	Turret                      bool // 0x0002 This item is a turret.
	Persistent                  bool // 0x0004 This item stays with you when you trade ships (persistent).
	CantSell                    bool // 0x0008 This item can't be sold.
	RemoveAfterBuy              bool // 0x0010 Remove any items of this type after purchase (useful for permits and other intangible purchases).
	PersistentMissionSet        bool // 0x0020 This item is persistent in the case where the player's ship is changed by a mission set operator. The item's normal persistence for when the player buys or captures a new ship is still controlled by the 0x0004 bit.
	RequireBitsOrHasOne         bool // 0x0100 Don't show this item unless the player meets the Require bits, or already has at least one of it.
	PriceProportionalToShipMass bool // 0x0200 This item's total price is proportional to the player's ship's mass. (ship class Mass field is multiplied by this item's Cost field)
	MassProportionalToShipMass  bool // 0x0400 This item's total mass (at purchase) is proportional to the player's ship's mass. (ship class Mass field is multiplied by this item's Mass field and then divided by 100) Only works for positive-mass items.
	SellAnywhere                bool // 0x0800 This item can be sold anywhere, regardless of tech level, requirements, or mission bits.
	HideHigherDispWeight        bool // 0x1000 When this item is available for sale, it prevents all higher-numbered items with equal DispWeight from being made available for sale at the same time.
	RankOutfit                  bool // 0x2000 This outfit appears in the Ranks section of the player info dialog instead of in the Extras section.
	AvailableBitsOrHasOne       bool // 0x4000 Don't show this item unless its Availability evaluates to true, or if the player already has at least one of it.
}

type RequireGovtID GovtID

func (r RequireGovtID) All() bool {
	return r == -1
}

func (r RequireGovtID) GovtAndAllies() bool {
	return 128 <= r && r <= 383
}

func (r RequireGovtID) GovtAndAlliesAndIndependent() bool {
	return 1128 <= r && r <= 1383
}

func (r RequireGovtID) NotGovtAndAllies() bool {
	return 2128 <= r && r <= 2383
}

func (r RequireGovtID) NotGovtAndAlliesAndIndependent() bool {
	return 3128 <= r && r <= 3383
}

type Outf struct {
	ID OutfID

	DispWeight   int16
	Mass         int16
	TechLevel    TechLevel
	ModType      [4]OutfMod
	Max          int16
	Flags        OutfFlags
	Cost         Credits
	Availability ControlBitTest
	OnPurchase   ControlBitFunction
	Contribute   FlagMask64
	Require      FlagMask64
	OnSell       ControlBitFunction
	ItemClass    int16
	ScanMask     FlagMask16
	BuyRandom    int16
	ShortName    string
	LCName       string
	LCPlural     string
	RequireGovt  RequireGovtID
}

func (o Outf) DescID() DescID {
	return DescID(o.ID) - resourcefork.ResourceForkIDOffset + DescIDOffsetOutfit
}

func (o Outf) PictID() PictID {
	return PictID(OutfPictIDOffset + int16(o.ID) - resourcefork.ResourceForkIDOffset)
}

func OutfFromResource(resource resourcefork.Resource) *Outf {
	return OutfFromBytes(OutfID(resource.ID), resource.Data)
}

func OutfFromBytes(id OutfID, b []byte) *Outf {
	flags := binary.BigEndian.Uint16(b[12:])

	t := &Outf{
		ID:         id,
		DispWeight: int16(binary.BigEndian.Uint16(b[0:])),
		Mass:       int16(binary.BigEndian.Uint16(b[2:])),
		TechLevel:  TechLevel(binary.BigEndian.Uint16(b[4:])),
		ModType: [4]OutfMod{
			{typeVal: OutfModType(binary.BigEndian.Uint16(b[6:])), value: int16(binary.BigEndian.Uint16(b[8:]))},
			{typeVal: OutfModType(binary.BigEndian.Uint16(b[18:])), value: int16(binary.BigEndian.Uint16(b[20:]))},
			{typeVal: OutfModType(binary.BigEndian.Uint16(b[22:])), value: int16(binary.BigEndian.Uint16(b[24:]))},
			{typeVal: OutfModType(binary.BigEndian.Uint16(b[26:])), value: int16(binary.BigEndian.Uint16(b[28:]))},
		},
		Max: int16(binary.BigEndian.Uint16(b[10:])),
		Flags: OutfFlags{
			FixedGun:                    flags&0x0001 == 0x0001,
			Turret:                      flags&0x0002 == 0x0002,
			Persistent:                  flags&0x0004 == 0x0004,
			CantSell:                    flags&0x0008 == 0x0008,
			RemoveAfterBuy:              flags&0x0010 == 0x0010,
			PersistentMissionSet:        flags&0x0020 == 0x0020,
			RequireBitsOrHasOne:         flags&0x0100 == 0x0100,
			PriceProportionalToShipMass: flags&0x0200 == 0x0200,
			MassProportionalToShipMass:  flags&0x0400 == 0x0400,
			SellAnywhere:                flags&0x0800 == 0x0800,
			HideHigherDispWeight:        flags&0x1000 == 0x1000,
			RankOutfit:                  flags&0x2000 == 0x2000,
			AvailableBitsOrHasOne:       flags&0x4000 == 0x4000,
		},
		Cost:         Credits(binary.BigEndian.Uint16(b[14:])),
		Availability: ControlBitTest(byteString(b[46:], 254)),
		OnPurchase:   ControlBitFunction(byteString(b[301:], 255)),
		Contribute:   FlagMask64(binary.BigEndian.Uint64(b[30:])),
		Require:      FlagMask64(binary.BigEndian.Uint64(b[38:])),
		OnSell:       ControlBitFunction(byteString(b[556:], 255)),
		ItemClass:    int16(binary.BigEndian.Uint16(b[1004:])),
		ScanMask:     FlagMask16(binary.BigEndian.Uint16(b[1006:])),
		BuyRandom:    int16(binary.BigEndian.Uint16(b[1008:])),
		ShortName:    byteString(b[811:], 63),
		LCName:       byteString(b[875:], 63),
		LCPlural:     byteString(b[939:], 64),
		RequireGovt:  RequireGovtID(binary.BigEndian.Uint16(b[1010:])),
	}

	return t
}
