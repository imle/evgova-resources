package resources

import (
	"encoding/binary"

	"github.com/imle/resourcefork"
)

// Spaceships are the heart of Nova, so the ship resource contains a lot of info. The name of a ship class, which is
// seen in the targeting display, corresponds to the name of the ship resource.

type ShipID IDType

type ShipFlags struct {
	SlowJump                          bool // 0x0001 Slow jumping (75% normal speed).
	SemiFastJump                      bool // 0x0002 Semi-fast jumping (125%)
	FastJump                          bool // 0x0004 Fast jumping (150%)
	UseFuelRegen                      bool // 0x0008 Player ship takes advantage of FuelRegen property
	DisableAt10Percent                bool // 0x0010 Ship is disabled at 10% armour instead of 33%
	AfterburnerAtAdvancedCombatRating bool // 0x0020 Ship has afterburner when player has an advanced combat rating
	AIHasAfterburner                  bool // 0x0040 Ship always has afterburner (for AIs only)
	AdvancedTargetStats               bool // 0x0100 Show % armour on target display instead of 'Shields Down'
	NoTargetStats                     bool // 0x0200 Don't show armour or shield state on status display
	PlanetTypeShip                    bool // 0x0400 Ship is a planet-type ship, and can only be hit by planet-type weapons
	TurretBlindSpotFront              bool // 0x1000 Ship's turrets have a blind spot to the front
	TurretBlindSpotSides              bool // 0x2000 Ship's turrets have a blind spot to the sides
	TurretBlindSpotRear               bool // 0x4000 Ship's turrets have a blind spot to the rear
	EscapeTypeShip                    bool // 0x8000 Ship is an escape ship type - if the player is carrying any ships of this type and decides to eject, he will fly off in a ship of this type (with random damage) instead of an escape pod.

	SwarmingBehavior        bool // 0x0001 Ship exhibits swarming behaviour.
	StandoffAttacks         bool // 0x0002 Ship prefers standoff attacks.
	NoTarget                bool // 0x0004 Ship can't be targeted.
	NoPointDefenceTargeting bool // 0x0008 Ship can be fired on by point defence systems.
	NoFighterVoices         bool // 0x0010 Don't use fighter voices.
	MustSlowToJump          bool // 0x0020 Ship can jump without slowing down.
	Inertialess             bool // 0x0040 Ship is inertia less.
	AIDockOnNoAmmo          bool // 0x0080 AI ships of this type will run away/dock if out of ammo for all ammo-using weapons.
	AICloakOnReload         bool // 0x0100 AI ships of this type will cloak when their weapon goes into burst reload.
	AICloakOnRetreat        bool // 0x0200 AI ships will cloak when running away.
	AICloakOnHyperspace     bool // 0x0400 AI ships will cloak when hyper spacing.
	AICloakFlying           bool // 0x0800 AI ships will cloak when just flying around.
	AIUncloakNearTarget     bool // 0x1000 AI ships will not uncloak until close to their target.
	AICloakOnWhenDocking    bool // 0x2000 AI ships will cloak when docking.
	AICloakUnderAttack      bool // 0x4000 AI ships will cloak when pre-emptively attacked.

	CanDestroyAsteroids               bool // 0x0001 Ship destroys asteroids.
	CanScoopAsteroidDebris            bool // 0x0002 Ship scoops asteroid debris.
	ShipIgnoresGravity                bool // 0x0010 Ship ignores gravity.
	ShipIgnoresDeadlyStellars         bool // 0x0020 Ship ignores deadly stellars.
	TurretsAboveShip                  bool // 0x0040 Ship's turreted shots appear above the ship instead of below.
	ShowOnlyIfAvailability            bool // 0x0100 Don't show ship in shipyard if Availability is false.
	ShowOnlyIfRequire                 bool // 0x0200 Don't show ship in shipyard if Require bits not met.
	HideShipTypesOfEqualDisplayWeight bool // 0x4000 When this ship is available for sale, it prevents all higher-numbered ship types with equal DispWeight from being made available for sale at the same time.
}

type ShipEscortType int16

const (
	ShipEscortTypeRuntime    ShipEscortType = iota - 1 // Have the game try to figure it out at runtime.
	ShipEscortTypeFighter                              // Fighter.
	ShipEscortTypeMediumShip                           // Medium Ship.
	ShipEscortTypeWarship                              // Warship.
	ShipEscortTypeFreighter                            // Freighter.
)

type InherentGovt IDType

// -1        No inherent combat govt or inherent attributes govt for this ship.
// 128-383   Ship is treated as being inherently of the govt with this ID, both for AI combat and attributes inheritance).
// 1128-1383 Ship has an inherent attributes govt with this ID (minus 1000) but no inherent combat govt.
// 2128-2383 Ship has an inherent combat govt with this ID (minus 2000) but no inherent attributes govt.

func (g InherentGovt) Parse() (govtID GovtID, attributes bool, combat bool) {
	return GovtID(g % 1000), g < 2000, g < 1128 || 2128 < g
}

type Ship struct {
	ID ShipID

	Holds         int16
	Shield        int16
	Accel         int16
	Speed         int16
	Maneuver      int16
	Fuel          int16
	FreeMass      int16
	Armour        int16
	ShieldRech    int16
	WeapType      [8]WeapID
	WeapCount     [8]int16
	AmmoLoad      [8]int16
	MaxGun        int16
	MaxTur        int16
	TechLevel     TechLevel
	Cost          Credits
	DeathDelay    int16
	ArmorRech     int16
	Explode1      BoomID
	Explode2      ExplodeType
	DispWeight    int16
	Mass          int16
	Length        int16
	InherentAI    AIType
	Crew          int16
	Strength      int16
	InherentGovt  InherentGovt
	Flags         ShipFlags
	PodCount      int16
	DefaultItems  [8]OutfID
	ItemCount     [8]int16
	FuelRegen     int16
	SkillVar      int16
	Availability  ControlBitTest
	AppearOn      ControlBitTest
	OnPurchase    ControlBitFunction
	Deionize      int16
	IonizeMax     int16
	KeyCarried    ShipID
	DefaultItems2 [8]OutfID
	ItemCount2    [8]int16
	Contribute    FlagMask64
	Require       FlagMask64
	BuyRandom     int16
	HireRandom    int16
	OnCapture     ControlBitFunction
	OnRetire      ControlBitFunction
	Subtitle      string
	UpgradeTo     ShipID
	EscUpgrdCost  Credits
	EscSellValue  Credits
	EscortType    AIType
	ShortName     string
	CommName      string
	LongName      string
	MovieFile     string
}

func (s ShipID) DescID() DescID {
	return DescID(s) - resourcefork.ResourceForkIDOffset + DescIDOffsetShipClass
}

func ShipFromResource(resource resourcefork.Resource) *Ship {
	return ShipFromBytes(ShipID(resource.ID), resource.Data)
}

func boomID(i int16) BoomID {
	if i == -1 {
		return -1
	}

	return BoomID(i) + resourcefork.ResourceForkIDOffset
}

func ShipFromBytes(id ShipID, b []byte) *Ship {
	flags := binary.BigEndian.Uint16(b[74:])
	flags2 := binary.BigEndian.Uint16(b[98:])
	flags3 := binary.BigEndian.Uint16(b[1830:])

	t := &Ship{
		ID:         id,
		Holds:      int16(binary.BigEndian.Uint16(b[0:])),
		Shield:     int16(binary.BigEndian.Uint16(b[2:])),
		Accel:      int16(binary.BigEndian.Uint16(b[4:])),
		Speed:      int16(binary.BigEndian.Uint16(b[6:])),
		Maneuver:   int16(binary.BigEndian.Uint16(b[8:])),
		Fuel:       int16(binary.BigEndian.Uint16(b[10:])),
		FreeMass:   int16(binary.BigEndian.Uint16(b[12:])),
		Armour:     int16(binary.BigEndian.Uint16(b[14:])),
		ShieldRech: int16(binary.BigEndian.Uint16(b[16:])),
		WeapType: [8]WeapID{
			WeapID(binary.BigEndian.Uint16(b[18:])),
			WeapID(binary.BigEndian.Uint16(b[20:])),
			WeapID(binary.BigEndian.Uint16(b[22:])),
			WeapID(binary.BigEndian.Uint16(b[24:])),
			WeapID(binary.BigEndian.Uint16(b[1742:])),
			WeapID(binary.BigEndian.Uint16(b[1744:])),
			WeapID(binary.BigEndian.Uint16(b[1746:])),
			WeapID(binary.BigEndian.Uint16(b[1748:])),
		},
		WeapCount: [8]int16{
			int16(binary.BigEndian.Uint16(b[26:])),
			int16(binary.BigEndian.Uint16(b[28:])),
			int16(binary.BigEndian.Uint16(b[30:])),
			int16(binary.BigEndian.Uint16(b[32:])),
			int16(binary.BigEndian.Uint16(b[1750:])),
			int16(binary.BigEndian.Uint16(b[1752:])),
			int16(binary.BigEndian.Uint16(b[1754:])),
			int16(binary.BigEndian.Uint16(b[1756:])),
		},
		AmmoLoad: [8]int16{
			int16(binary.BigEndian.Uint16(b[34:])),
			int16(binary.BigEndian.Uint16(b[36:])),
			int16(binary.BigEndian.Uint16(b[38:])),
			int16(binary.BigEndian.Uint16(b[40:])),
			int16(binary.BigEndian.Uint16(b[1758:])),
			int16(binary.BigEndian.Uint16(b[1760:])),
			int16(binary.BigEndian.Uint16(b[1762:])),
			int16(binary.BigEndian.Uint16(b[1764:])),
		},
		MaxGun:       int16(binary.BigEndian.Uint16(b[42:])),
		MaxTur:       int16(binary.BigEndian.Uint16(b[44:])),
		TechLevel:    TechLevel(binary.BigEndian.Uint16(b[46:])),
		Cost:         Credits(binary.BigEndian.Uint16(b[48:])),
		DeathDelay:   int16(binary.BigEndian.Uint16(b[52:])),
		ArmorRech:    int16(binary.BigEndian.Uint16(b[54:])),
		Explode1:     boomID(int16(binary.BigEndian.Uint16(b[56:]))),
		Explode2:     ExplodeType(binary.BigEndian.Uint16(b[58:])),
		DispWeight:   int16(binary.BigEndian.Uint16(b[60:])),
		Mass:         int16(binary.BigEndian.Uint16(b[62:])),
		Length:       int16(binary.BigEndian.Uint16(b[64:])),
		InherentAI:   AIType(binary.BigEndian.Uint16(b[66:])),
		Crew:         int16(binary.BigEndian.Uint16(b[68:])),
		Strength:     int16(binary.BigEndian.Uint16(b[70:])),
		InherentGovt: InherentGovt(binary.BigEndian.Uint16(b[72:])),
		Flags: ShipFlags{
			SlowJump:                          flags&0x0001 == 0x0001,
			SemiFastJump:                      flags&0x0002 == 0x0002,
			FastJump:                          flags&0x0004 == 0x0004,
			UseFuelRegen:                      flags&0x0008 == 0x0008,
			DisableAt10Percent:                flags&0x0010 == 0x0010,
			AfterburnerAtAdvancedCombatRating: flags&0x0020 == 0x0020,
			AIHasAfterburner:                  flags&0x0040 == 0x0040,
			AdvancedTargetStats:               flags&0x0100 == 0x0100,
			NoTargetStats:                     flags&0x0200 == 0x0200,
			PlanetTypeShip:                    flags&0x0400 == 0x0400,
			TurretBlindSpotFront:              flags&0x1000 == 0x1000,
			TurretBlindSpotSides:              flags&0x2000 == 0x2000,
			TurretBlindSpotRear:               flags&0x4000 == 0x4000,
			EscapeTypeShip:                    flags&0x8000 == 0x8000,
			SwarmingBehavior:                  flags2&0x0001 == 0x0001,
			StandoffAttacks:                   flags2&0x0002 == 0x0002,
			NoTarget:                          flags2&0x0004 == 0x0004,
			NoPointDefenceTargeting:           flags2&0x0008 == 0x0008,
			NoFighterVoices:                   flags2&0x0010 == 0x0010,
			MustSlowToJump:                    flags2&0x0020 == 0x0020,
			Inertialess:                       flags2&0x0040 == 0x0040,
			AIDockOnNoAmmo:                    flags2&0x0080 == 0x0080,
			AICloakOnReload:                   flags2&0x0100 == 0x0100,
			AICloakOnRetreat:                  flags2&0x0200 == 0x0200,
			AICloakOnHyperspace:               flags2&0x0400 == 0x0400,
			AICloakFlying:                     flags2&0x0800 == 0x0800,
			AIUncloakNearTarget:               flags2&0x1000 == 0x1000,
			AICloakOnWhenDocking:              flags2&0x2000 == 0x2000,
			AICloakUnderAttack:                flags2&0x4000 == 0x4000,
			CanDestroyAsteroids:               flags3&0x0001 == 0x0001,
			CanScoopAsteroidDebris:            flags3&0x0002 == 0x0002,
			ShipIgnoresGravity:                flags3&0x0010 == 0x0010,
			ShipIgnoresDeadlyStellars:         flags3&0x0020 == 0x0020,
			TurretsAboveShip:                  flags3&0x0040 == 0x0040,
			ShowOnlyIfAvailability:            flags3&0x0100 == 0x0100,
			ShowOnlyIfRequire:                 flags3&0x0200 == 0x0200,
			HideShipTypesOfEqualDisplayWeight: flags3&0x4000 == 0x4000,
		},
		PodCount: int16(binary.BigEndian.Uint16(b[76:])),
		DefaultItems: [8]OutfID{
			OutfID(binary.BigEndian.Uint16(b[78:])),
			OutfID(binary.BigEndian.Uint16(b[80:])),
			OutfID(binary.BigEndian.Uint16(b[82:])),
			OutfID(binary.BigEndian.Uint16(b[84:])),
			OutfID(binary.BigEndian.Uint16(b[880:])),
			OutfID(binary.BigEndian.Uint16(b[882:])),
			OutfID(binary.BigEndian.Uint16(b[884:])),
			OutfID(binary.BigEndian.Uint16(b[886:])),
		},
		ItemCount: [8]int16{
			int16(binary.BigEndian.Uint16(b[86:])),
			int16(binary.BigEndian.Uint16(b[88:])),
			int16(binary.BigEndian.Uint16(b[90:])),
			int16(binary.BigEndian.Uint16(b[92:])),
			int16(binary.BigEndian.Uint16(b[888:])),
			int16(binary.BigEndian.Uint16(b[890:])),
			int16(binary.BigEndian.Uint16(b[892:])),
			int16(binary.BigEndian.Uint16(b[894:])),
		},
		FuelRegen:    int16(binary.BigEndian.Uint16(b[94:])),
		SkillVar:     int16(binary.BigEndian.Uint16(b[96:])),
		Availability: ControlBitTest(byteString(b[:108], 254)),
		AppearOn:     ControlBitTest(byteString(b[:363], 254)),
		OnPurchase:   ControlBitFunction(byteString(b[:618], 255)),
		Deionize:     int16(binary.BigEndian.Uint16(b[874:])),
		IonizeMax:    int16(binary.BigEndian.Uint16(b[876:])),
		KeyCarried:   ShipID(binary.BigEndian.Uint16(b[878:])),
		Contribute:   FlagMask64(binary.BigEndian.Uint64(b[100:])),
		Require:      FlagMask64(binary.BigEndian.Uint64(b[896:])),
		BuyRandom:    int16(binary.BigEndian.Uint16(b[904:])),
		HireRandom:   int16(binary.BigEndian.Uint16(b[906:])),
		OnCapture:    ControlBitFunction(byteString(b[:908], 255)),
		OnRetire:     ControlBitFunction(byteString(b[:1163], 255)),
		Subtitle:     byteString(b[:1766], 64),
		UpgradeTo:    ShipID(binary.BigEndian.Uint16(b[1832:])),
		EscUpgrdCost: Credits(binary.BigEndian.Uint32(b[1834:])),
		EscSellValue: Credits(binary.BigEndian.Uint32(b[1838:])),
		EscortType:   AIType(binary.BigEndian.Uint16(b[1829:])),
		ShortName:    byteString(b[1486:], 64),
		CommName:     byteString(b[1550:], 32),
		LongName:     byteString(b[1582:], 128),
		MovieFile:    byteString(b[1710:], 32),
	}

	return t
}
