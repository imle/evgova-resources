package resources

import (
	"encoding/binary"
	"image/color"

	"github.com/imle/resourcefork"
)

// The pers resource defines the characteristics of an AI personality – that is, a specific person the player can
// encounter in the game. These AI-people have their names (which are also the names of the associated përs resource)
// displayed on the target-info display in place of the name of their ship class. When ships are created, there is
// a 5% chance that a specific AI-person will also be created. (obviously, as AI-people are killed off, they cease
// to appear in the game.)

type PersID IDType

type PersAggression uint8

const (
	PersAggressionClose PersAggression = iota + 1
	PersAggressionMedium
	PersAggressionFar
)

type PersFlags struct {
	Grudge                          bool // 0x0001 The special ship will hold a grudge if attacked, and will subsequently attack the player wherever the twain shall meet.
	EscapePodAndAfterburner         bool // 0x0002 Uses escape pod & has afterburner .
	HailQuoteForGrudge              bool // 0x0004 HailQuote only shown when ship has a grudge against the player.
	HailQuoteIfLiked                bool // 0x0008 HailQuote only shown when ship likes player.
	HailQuoteIfAttacking            bool // 0x0010 Only show HailQuote when ship begins to attack the player.
	HailQuoteIfDisabled             bool // 0x0020 Only show HailQuote when ship is disabled.
	UseAsSpecialShipForLinkMission  bool // 0x0040 When LinkMission is accepted with a single SpecialShip, replace it with this ship while removing this one from play. This is generally only useful for escort and refuel-a-ship missions.
	CommQuoteOnce                   bool // 0x0080 Only show quote once.
	DeactivateOnLinkMissionAccept   bool // 0x0100 Deactivate ship (i.e. don't make it show up again) after accepting its LinkMission.
	LinkMissionOnBoard              bool // 0x0200 Offer ship's LinkMission when boarding it instead of when hailing it.
	CommQuoteOnlyIfLinkMissionAvail bool // 0x0400 Don't show quote when ship's LinkMission is not available.
	LeaveOnLinkMissionAccept        bool // 0x0800 Make ship leave after accepting its LinkMission.
	NoLinkMissionOnWimpyTrader      bool // 0x1000 Don't offer if player is flying a wimpy freighter (aiType 1).
	NoLinkMissionOnBeefyTrader      bool // 0x2000 Don't offer if player is flying a beefy freighter (aiType 2).
	NoLinkMissionOnWarship          bool // 0x4000 Don't offer if player is flying a warship (aiType 3).
	DisasterInfoOnHail              bool // 0x8000 Show disaster info when hailing.
	StartsWithoutFuel               bool // 0x0001 This person starts with zero fuel.
}

type PersLinkSyst int16

type Pers struct {
	ID PersID

	LinkSyst    PersLinkSyst
	Govt        GovtID
	AIType      AIType
	Aggression  PersAggression
	Coward      int16
	ShipType    ShipID
	WeapType    [8]WeapID
	WeapCount   [8]int16
	AmmoLoad    [8]int16
	Credits     Credits
	ShieldMod   int16
	HailPict    PictID
	CommQuote   StrAID
	HailQuote   StrAID
	LinkMission MisnID
	Flags       PersFlags
	ActiveOn    ControlBitTest
	Subtitle    string
	GrantClass  int16
	GrantProb   int16
	GrantCount  int16
	Colour      color.Color
}

func PersFromResource(resource resourcefork.Resource) *Pers {
	return PersFromBytes(PersID(resource.ID), resource.Data)
}

func PersFromBytes(id PersID, b []byte) *Pers {
	flags := binary.BigEndian.Uint16(b[50:])

	t := &Pers{
		ID:         id,
		LinkSyst:   PersLinkSyst(binary.BigEndian.Uint16(b[2:])),
		Govt:       GovtID(binary.BigEndian.Uint16(b[2:])),
		AIType:     AIType(binary.BigEndian.Uint16(b[4:])),
		Aggression: PersAggression(binary.BigEndian.Uint16(b[6:])),
		Coward:     int16(binary.BigEndian.Uint16(b[8:])),
		ShipType:   ShipID(binary.BigEndian.Uint16(b[10:])),
		WeapType: [8]WeapID{
			WeapID(binary.BigEndian.Uint16(b[12:])),
			WeapID(binary.BigEndian.Uint16(b[14:])),
			WeapID(binary.BigEndian.Uint16(b[16:])),
			WeapID(binary.BigEndian.Uint16(b[18:])),
		},
		WeapCount: [8]int16{
			int16(binary.BigEndian.Uint16(b[20:])),
			int16(binary.BigEndian.Uint16(b[22:])),
			int16(binary.BigEndian.Uint16(b[24:])),
			int16(binary.BigEndian.Uint16(b[26:])),
		},
		AmmoLoad: [8]int16{
			int16(binary.BigEndian.Uint16(b[28:])),
			int16(binary.BigEndian.Uint16(b[30:])),
			int16(binary.BigEndian.Uint16(b[32:])),
			int16(binary.BigEndian.Uint16(b[34:])),
		},
		Credits:     Credits(binary.BigEndian.Uint32(b[36:])),
		ShieldMod:   int16(binary.BigEndian.Uint16(b[40:])),
		HailPict:    PictID(binary.BigEndian.Uint16(b[42:])),
		CommQuote:   StrAID(binary.BigEndian.Uint16(b[44:])),
		HailQuote:   StrAID(binary.BigEndian.Uint16(b[46:])),
		LinkMission: MisnID(binary.BigEndian.Uint16(b[48:])),
		Flags: PersFlags{
			Grudge:                          flags&0x0001 == 0x0001,
			EscapePodAndAfterburner:         flags&0x0002 == 0x0002,
			HailQuoteForGrudge:              flags&0x0004 == 0x0004,
			HailQuoteIfLiked:                flags&0x0008 == 0x0008,
			HailQuoteIfAttacking:            flags&0x0010 == 0x0010,
			HailQuoteIfDisabled:             flags&0x0020 == 0x0020,
			UseAsSpecialShipForLinkMission:  flags&0x0040 == 0x0040,
			CommQuoteOnce:                   flags&0x0080 == 0x0080,
			DeactivateOnLinkMissionAccept:   flags&0x0100 == 0x0100,
			LinkMissionOnBoard:              flags&0x0200 == 0x0200,
			CommQuoteOnlyIfLinkMissionAvail: flags&0x0400 == 0x0400,
			LeaveOnLinkMissionAccept:        flags&0x0800 == 0x0800,
			NoLinkMissionOnWimpyTrader:      flags&0x1000 == 0x1000,
			NoLinkMissionOnBeefyTrader:      flags&0x2000 == 0x2000,
			NoLinkMissionOnWarship:          flags&0x4000 == 0x4000,
			DisasterInfoOnHail:              flags&0x8000 == 0x8000,
			StartsWithoutFuel:               flags&0x0001 == 0x0001,
		},
		ActiveOn:   ControlBitTest(byteString(b[52:], 254)),
		Subtitle:   byteString(b[314:], 63),
		GrantClass: int16(binary.BigEndian.Uint16(b[308:])),
		GrantProb:  int16(binary.BigEndian.Uint16(b[312:])),
		GrantCount: int16(binary.BigEndian.Uint16(b[310:])),
		Colour: color.RGBA{
			A: b[378],
			R: b[379],
			G: b[380],
			B: b[381],
		},
	}

	return t
}
