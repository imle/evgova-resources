package resources

import (
	"encoding/binary"
	"image/color"

	"github.com/imle/resourcefork"
)

// A govt resource defines the parameters for a government, which is in turn defined as "any collection of ships and
// planets that react collectively to the actions of the player and other ships." Governments keep track of how they
// feel toward you, and they can also have set enemies and allies.

type GovtID IDType

type GovtFlags struct {
	Xenophobic                     bool // 0x0001 Xenophobic (Warships of this govt attack everyone except their allies. Useful for making pirates and other nasties.)
	AlwaysAttackIfPlayerIsCriminal bool // 0x0002 Ships of this govt will attack the player in non-allied systems if he's a criminal there (useful for making one govt care only about the player's actions on its home turf, while another is nosy and enforces its own laws everywhere it goes).
	AlwaysAttack                   bool // 0x0004 Always attacks player.
	ImmuneToPlayerWeapons          bool // 0x0008 Player's shots won't hit ships of this govt.
	RetreatAtQuarterShields        bool // 0x0010 Warships of this govt will retreat when their shields drop below 25% - otherwise they fight to the death.
	NoHelpFromNonAllies            bool // 0x0020 Nosy ships of other non-allied governments ignore ships of this govt that are under attack.
	NeverAttacksPlayer             bool // 0x0040 Never attacks player (also, player's weapons can't hit them).
	FreighterBetterJam             bool // 0x0080 Freighters (i.e. AiTypes 1 and 2) for this particular government have 50% of the standard InherentJam value for warships (AiType 3) of the same government.
	PersNeverDie                   bool // 0x0100 'pers' ships of this govt won't use escape pod, but will act as if they did.
	WarshipBriberyAccepted         bool // 0x0200 Warships will take bribes.
	NeverRespondsToPlayer          bool // 0x0400 Can't hail ships of this govt. (if a ship type has an inherent attributes govt which includes this flag, all ships of that type will inherit this property)
	StartDerelict                  bool // 0x0800 Ships of this govt start out disabled (derelicts). Note that ships of other governments don't care if you attack or board derelict govt ships.
	PlunderThenKill                bool // 0x1000 Warships will plunder non-mission, non-player enemies before destroying them.
	FreighterBriberyAccepted       bool // 0x2000 Freighters will take bribes.
	PlanetBriberyAccepted          bool // 0x4000 Planets of this govt will take bribes
	BigMoneyBribes                 bool // 0x8000 Ships of this govt taking bribes will demand a larger percentage of your cash supply, and their planets will always take bribes (useful for pirates).

	NoMercyOrAssist      bool // 0x0001 When hailing ships of this govt, the request assistance / beg for mercy button is disabled and the govt is not talkative.
	MinorGovt            bool // 0x0002 This govt is considered "minor" for the purposes of drawing the political boundaries on the map.
	NoBoundaryAffect     bool // 0x0004 This govt's systems don't affect the political boundaries on the map.
	NoDistressOrGreeting bool // 0x0008 Ships of this govt don't send distress messages and don't respond with greetings when hailed (if a ship type has an inherent attributes govt which includes this flag, all ships of that type will inherit this property)
	RoadsideAssistance   bool // 0x0010 Roadside Assistance - Ships of this govt will always repair or refuel the player for free.
	NoHyperGate          bool // 0x0020 Ships of this govt don't use hypergates.
	PreferHyperGate      bool // 0x0040 Ships of this govt prefer to use hypergates instead of jumping out.
	PreferWormhole       bool // 0x0080 Ships of this govt prefer to use wormholes instead of jumping out.
}

type Govt struct {
	ID GovtID

	VoiceType    int16
	Flags        GovtFlags
	SkillMult    int16
	CrimeTol     int16
	ScanFine     int16
	SmugPenalty  int16
	DisabPenalty int16
	BoardPenalty int16
	KillPenalty  int16
	ShootPenalty int16 // Ignored
	InitialRec   int16
	MaxOdds      int16
	Class        [4]int16
	Ally         [4]int16
	Enemy        [4]int16
	Interface    IntfID
	NewsPic      PictID
	ScanMask     FlagMask16
	Require      FlagMask64
	InhJam       [4]int16
	MediumName   string
	Colour       color.Color
	ShipColour   color.Color
	CommName     string
	TargetCode   string
}

func GovtFromResource(resource resourcefork.Resource) *Govt {
	return GovtFromBytes(GovtID(resource.ID), resource.Data)
}

func GovtFromBytes(id GovtID, b []byte) *Govt {
	flags1 := binary.BigEndian.Uint16(b[2:])
	flags2 := binary.BigEndian.Uint16(b[4:])

	t := &Govt{
		ID:        id,
		VoiceType: int16(binary.BigEndian.Uint16(b[0:])),
		Flags: GovtFlags{
			Xenophobic:                     flags1&0x0001 == 0x0001,
			AlwaysAttackIfPlayerIsCriminal: flags1&0x0002 == 0x0002,
			AlwaysAttack:                   flags1&0x0004 == 0x0004,
			ImmuneToPlayerWeapons:          flags1&0x0008 == 0x0008,
			RetreatAtQuarterShields:        flags1&0x0010 == 0x0010,
			NoHelpFromNonAllies:            flags1&0x0020 == 0x0020,
			NeverAttacksPlayer:             flags1&0x0040 == 0x0040,
			FreighterBetterJam:             flags1&0x0080 == 0x0080,
			PersNeverDie:                   flags1&0x0100 == 0x0100,
			WarshipBriberyAccepted:         flags1&0x0200 == 0x0200,
			NeverRespondsToPlayer:          flags1&0x0400 == 0x0400,
			StartDerelict:                  flags1&0x0800 == 0x0800,
			PlunderThenKill:                flags1&0x1000 == 0x1000,
			FreighterBriberyAccepted:       flags1&0x2000 == 0x2000,
			PlanetBriberyAccepted:          flags1&0x4000 == 0x4000,
			BigMoneyBribes:                 flags1&0x8000 == 0x8000,
			NoMercyOrAssist:                flags2&0x0001 == 0x0001,
			MinorGovt:                      flags2&0x0002 == 0x0002,
			NoBoundaryAffect:               flags2&0x0004 == 0x0004,
			NoDistressOrGreeting:           flags2&0x0008 == 0x0008,
			RoadsideAssistance:             flags2&0x0010 == 0x0010,
			NoHyperGate:                    flags2&0x0020 == 0x0020,
			PreferHyperGate:                flags2&0x0040 == 0x0040,
			PreferWormhole:                 flags2&0x0080 == 0x0080,
		},
		ScanFine:     int16(binary.BigEndian.Uint16(b[6:])),
		CrimeTol:     int16(binary.BigEndian.Uint16(b[8:])),
		SmugPenalty:  int16(binary.BigEndian.Uint16(b[10:])),
		DisabPenalty: int16(binary.BigEndian.Uint16(b[12:])),
		BoardPenalty: int16(binary.BigEndian.Uint16(b[14:])),
		KillPenalty:  int16(binary.BigEndian.Uint16(b[16:])),
		ShootPenalty: int16(binary.BigEndian.Uint16(b[18:])),
		InitialRec:   int16(binary.BigEndian.Uint16(b[20:])),
		MaxOdds:      int16(binary.BigEndian.Uint16(b[22:])),
		Class: [4]int16{
			int16(binary.BigEndian.Uint16(b[24:])),
			int16(binary.BigEndian.Uint16(b[26:])),
			int16(binary.BigEndian.Uint16(b[28:])),
			int16(binary.BigEndian.Uint16(b[30:])),
		},
		Ally: [4]int16{
			int16(binary.BigEndian.Uint16(b[32:])),
			int16(binary.BigEndian.Uint16(b[34:])),
			int16(binary.BigEndian.Uint16(b[36:])),
			int16(binary.BigEndian.Uint16(b[38:])),
		},
		Enemy: [4]int16{
			int16(binary.BigEndian.Uint16(b[40:])),
			int16(binary.BigEndian.Uint16(b[42:])),
			int16(binary.BigEndian.Uint16(b[44:])),
			int16(binary.BigEndian.Uint16(b[46:])),
		},
		SkillMult: int16(binary.BigEndian.Uint16(b[48:])),
		ScanMask:  FlagMask16(binary.BigEndian.Uint16(b[50:])),
		Require:   FlagMask64(binary.BigEndian.Uint64(b[84:])),
		InhJam: [4]int16{
			int16(binary.BigEndian.Uint16(b[92:])),
			int16(binary.BigEndian.Uint16(b[94:])),
			int16(binary.BigEndian.Uint16(b[96:])),
			int16(binary.BigEndian.Uint16(b[98:])),
		},
		Interface: IntfID(binary.BigEndian.Uint16(b[172:])),
		NewsPic:   PictID(binary.BigEndian.Uint16(b[174:])),
		Colour: color.RGBA{
			A: b[164],
			R: b[165],
			G: b[166],
			B: b[167],
		},
		ShipColour: color.RGBA{
			A: b[168],
			R: b[169],
			G: b[170],
			B: b[171],
		},
		CommName:   byteString(b[52:], 16),
		TargetCode: byteString(b[68:], 15),
		MediumName: byteString(b[100:], 31),
	}

	return t
}
