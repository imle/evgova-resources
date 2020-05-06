package resources

import (
	"encoding/binary"
	"image/color"

	"github.com/imle/resourcefork"
)

// The weap resource, surprisingly, stores info on Nova's weapons. The name of the weap resource is used as the weapon
// name in the weaponry section of the status display.

type WeapID IDType

type WeapGuidance int8

const (
	WeapGuidanceUnguided            WeapGuidance = iota - 1 // -1 Unguided projectile.
	WeapGuidanceBeam                                        // 0 Beam weapon (see below).
	WeapGuidanceHoming                                      // 1 Homing weapon (see Seeker field below).
	_                                                       // 2 (unused).
	WeapGuidanceTurretBeam                                  // 3 Turreted beam.
	WeapGuidanceTurret                                      // 4 Turreted, unguided projectile.
	WeapGuidanceFreeFallBomb                                // 5 Freefall bomb (launched at 80% of the ship's current velocity, "weathervanes" into the "wind."
	WeapGuidanceFreeFlightRocket                            // 6 Freeflight rocket (launched straight ahead, accelerates to its maximum velocity).
	WeapGuidanceTurretQuadrantFront                         // 7 Front-quadrant turret, (can fire +/-45¡ off the ship's nose) fires straight ahead if no target.
	WeapGuidanceTurretQuadrantRear                          // 8 Rear-quadrant turret (can fire +/-45¡ off the ship's tail).
	WeapGuidanceTurretPointDefence                          // 9 Point defence turret (fires automatically at incoming guided weapons and nearby ships).
	WeapGuidanceBeamPointDefence                            // 10 Point defence beam (fires automatically at incoming guided weapons and nearby ships).
	WeapGuidanceCarriedShip         = 99                    // 99 Carried ship (AmmoType is the ID of the ship class).
)

type SpinWeapIndex SpinID // 0-255 => 3000-3255

func (i SpinWeapIndex) SpinID() SpinID {
	return SpinID(i) + 3000
}

type SndWeapIndex SndID // 0-63 => 200-263

func (i SndWeapIndex) SndID() SndID {
	return SndID(i) + 200
}

type WeapFlags struct {
	SpinWeaponGraphic            bool // 0x0001 Spin the weapon's graphic continuously (rate of frame advance is controlled by the BeamWidth field as detailed below).
	SecondTriggerWeapon          bool // 0x0002 Weapon fired by second trigger.
	StartOnFirstFrameOfAnimation bool // 0x0004 For cycling weapons, always start on the first frame of the animation.
	DontFireAtFastShips          bool // 0x0008 For guided weapons, don't fire at fast ships (ships with turn rate > 3).
	LoopedSound                  bool // 0x0010 Weapon's sound is looped rather than played repeatedly.
	IgnoreShields                bool // 0x0020 Weapon passes through shields (use sparingly!).
	MultipleOfTypeFire           bool // 0x0040 Multiple weapons of this type fire simultaneously.
	NoPointDefenceTargeting      bool // 0x0080 Weapon can't be targeted by point defence systems (works only for homing weapons).
	BlastDoesNoPlayerDamage      bool // 0x0100 Weapon's blast doesn't hurt the player.
	SmallSmoke                   bool // 0x0200 Weapon generates small smoke.
	BigSmoke                     bool // 0x0400 Weapon generates big smoke.
	LongerSmokeLifetime          bool // 0x0800 Weapon's smoke trail is more persistent.
	TurretBlindFront             bool // 0x1000 Turreted weapon has a blind spot to the front.
	TurretBlindSides             bool // 0x2000 Turreted weapon has a blind spot to the sides.
	TurretBlindRear              bool // 0x4000 Turreted weapon has a blind spot to the rear .
	ShotDetonatesOnLastFrame     bool // 0x8000 Shot detonates at the end of its lifespan (useful for flak-type weapons).

	AnimationStayOnFirstFrameUntilProxSafetyExpired bool // 0x0001 For cycling weapons, keep the graphic on the first frame until the weapon's ProxSafety count has expired.
	AnimationStopOnLastFrame                        bool // 0x0002 For cycling weapons, stop the graphic on the last frame.
	ProximityDetonatorIgnoreAsteroids               bool // 0x0004 Proximity detonator ignores asteroids.
	ProximityDetonatorTriggeredByNonTargets         bool // 0x0008 Proximity detonator is triggered by ships other than the target (for guided weapons).
	SubmunitionsFireOnNearestTarget                 bool // 0x0010 Submunitions fire toward nearest valid target.
	NoSubmunitionsOnShotExpire                      bool // 0x0020 Don't launch submunitions when the shot expires.
	NoAmmoStats                                     bool // 0x0040 Don't show weapon's ammo quantity on the status display.
	OnlyFireWithKeyCarriedAboard                    bool // 0x0080 This weapon can only be fired when there is at least one ship of this ship's KeyCarried type aboard.
	AIDoesNotUse                                    bool // 0x0100 AI ships won't use this weapon.
	UseShipWeaponSprite                             bool // 0x0200 This weapon uses the ship's weapon sprite, if applicable.
	PlanetTypeWeapon                                bool // 0x0400 Weapon is a planet-type weapon, and can only hit planet-type ships or destroyable stellars.
	NoSelectOutOfAmmo                               bool // 0x0800 Don't allow this weapon to be selected or displayed if it is out of ammo.
	NoDestroy                                       bool // 0x1000 Weapon can disable but not destroy.
	BeamUnderneathShip                              bool // 0x2000 For beam weapons, display the beam underneath ships instead of on top of them.
	CloakedFire                                     bool // 0x4000 Weapon can be fired while cloaked.
	TenfoldDamageToAsteroids                        bool // 0x8000 Weapon does x10 mass damage to asteroids.

	AmmoUsedAtEndOfBurstCycle bool // 0x0001 Weapon will only use ammo at the end of a burst cycle.
	ShotsAreTranslucent       bool // 0x0002 Weapon's shots are translucent
	OnlyOneShotAtATime        bool // 0x0004 Firing ship can't fire another shot of this type until the previous one expires or hits something.
	ClosestExitPoint          bool // 0x0010 Weapon fires from whatever weapon exit point is closest to the target.
	ExclusiveWeapon           bool // 0x0020 Weapon is exclusive - no other weapons on the ship can fire while this weapon is firing or reloading.
}

type WeapSeeker struct {
	IgnoresAsteroids       bool // 0x0001 Passes over asteroids.
	DecoyedAsteroids       bool // 0x0002 Decoyed by asteroids. ?????
	ConfusedByInterference bool // 0x0008 Confused by sensor interference.
	TurnsAwayIfJammed      bool // 0x0010 Turns away if jammed.
	NoFireIfIonized        bool // 0x0020 Can't fire if ship is ionized.
	NoLockIfNotDirectedAt  bool // 0x4000 Loses lock if target not directly ahead.
	SelfDamageIfJammed     bool // 0x8000 May attack parent ship if jammed.
}

type WeapExitType int32

const (
	WeapExitTypeIgnored WeapExitType = iota - 1
	WeapExitTypeGunPosXY
	WeapExitTypeTurretPosXY
	WeapExitTypeGuidedPosXY
	WeapExitTypeBeamPosXY
)

type Weap struct {
	ID WeapID

	Reload       FrameCount
	Count        FrameCount
	MassDmg      int16
	EnergyDmg    int16
	Guidance     WeapGuidance
	Speed        int16
	AmmoType     int16
	Graphic      SpinWeapIndex
	Inaccuracy   int16
	Sound        SndWeapIndex
	Impact       int16
	ExplodeType  ExplodeType
	ProxRadius   int16
	BlastRadius  int16
	Flags        WeapFlags
	Seeker       WeapSeeker
	SmokeSet     CicnID
	Decay        FrameCount
	Particles    int16
	PartVel      int16
	PartLifeMin  int16
	PartLifeMax  int16
	PartColor    color.Color
	BeamLength   int16
	BeamWidth    int16
	Falloff      int16
	BeamColor    color.Color
	CoronaColor  color.Color
	SubCount     int16
	SubType      WeapID
	SubTheta     int16
	SubLimit     int16
	ProxSafety   int16
	Ionization   int16
	HitParticles int16
	HitPartLife  int16
	HitPartVel   int16
	HitPartColor color.Color
	ExitType     WeapExitType
	BurstCount   int16
	BurstReload  int16
	JamVuln1     int16
	JamVuln2     int16
	JamVuln3     int16
	JamVuln4     int16
	Durability   int16
	GuidedTurn   int16
	MaxAmmo      int16
	Recoil       int16
	LiDensity    int16
	LiAmplitude  int16
	IonizeColor  color.Color
}

func WeapFromResource(resource resourcefork.Resource) *Weap {
	return WeapFromBytes(WeapID(resource.ID), resource.Data)
}

func WeapFromBytes(id WeapID, b []byte) *Weap {
	flags := binary.BigEndian.Uint16(b[28:])
	flags2 := binary.BigEndian.Uint16(b[72:])
	flags3 := binary.BigEndian.Uint16(b[102:])
	flagsSeeker := binary.BigEndian.Uint16(b[30:])

	t := &Weap{
		ID:          id,
		Reload:      FrameCount(binary.BigEndian.Uint16(b[0:])),
		Count:       FrameCount(binary.BigEndian.Uint16(b[2:])),
		MassDmg:     int16(binary.BigEndian.Uint16(b[4:])),
		EnergyDmg:   int16(binary.BigEndian.Uint16(b[6:])),
		Guidance:    WeapGuidance(binary.BigEndian.Uint16(b[8:])),
		Speed:       int16(binary.BigEndian.Uint16(b[10:])),
		AmmoType:    int16(binary.BigEndian.Uint16(b[12:])),
		Graphic:     SpinWeapIndex(binary.BigEndian.Uint16(b[14:])),
		Inaccuracy:  int16(binary.BigEndian.Uint16(b[16:])),
		Sound:       SndWeapIndex(binary.BigEndian.Uint16(b[18:])),
		Impact:      int16(binary.BigEndian.Uint16(b[20:])),
		ExplodeType: ExplodeType(binary.BigEndian.Uint16(b[22:])),
		ProxRadius:  int16(binary.BigEndian.Uint16(b[24:])),
		BlastRadius: int16(binary.BigEndian.Uint16(b[26:])),
		Flags: WeapFlags{
			SpinWeaponGraphic:            flags&0x0001 == 0x0001,
			SecondTriggerWeapon:          flags&0x0002 == 0x0002,
			StartOnFirstFrameOfAnimation: flags&0x0004 == 0x0004,
			DontFireAtFastShips:          flags&0x0008 == 0x0008,
			LoopedSound:                  flags&0x0010 == 0x0010,
			IgnoreShields:                flags&0x0020 == 0x0020,
			MultipleOfTypeFire:           flags&0x0040 == 0x0040,
			NoPointDefenceTargeting:      flags&0x0080 == 0x0080,
			BlastDoesNoPlayerDamage:      flags&0x0100 == 0x0100,
			SmallSmoke:                   flags&0x0200 == 0x0200,
			BigSmoke:                     flags&0x0400 == 0x0400,
			LongerSmokeLifetime:          flags&0x0800 == 0x0800,
			TurretBlindFront:             flags&0x1000 == 0x1000,
			TurretBlindSides:             flags&0x2000 == 0x2000,
			TurretBlindRear:              flags&0x4000 == 0x4000,
			ShotDetonatesOnLastFrame:     flags&0x8000 == 0x8000,
			AnimationStayOnFirstFrameUntilProxSafetyExpired: flags2&0x0001 == 0x0001,
			AnimationStopOnLastFrame:                        flags2&0x0002 == 0x0002,
			ProximityDetonatorIgnoreAsteroids:               flags2&0x0004 == 0x0004,
			ProximityDetonatorTriggeredByNonTargets:         flags2&0x0008 == 0x0008,
			SubmunitionsFireOnNearestTarget:                 flags2&0x0010 == 0x0010,
			NoSubmunitionsOnShotExpire:                      flags2&0x0020 == 0x0020,
			NoAmmoStats:                                     flags2&0x0040 == 0x0040,
			OnlyFireWithKeyCarriedAboard:                    flags2&0x0080 == 0x0080,
			AIDoesNotUse:                                    flags2&0x0100 == 0x0100,
			UseShipWeaponSprite:                             flags2&0x0200 == 0x0200,
			PlanetTypeWeapon:                                flags2&0x0400 == 0x0400,
			NoSelectOutOfAmmo:                               flags2&0x0800 == 0x0800,
			NoDestroy:                                       flags2&0x1000 == 0x1000,
			BeamUnderneathShip:                              flags2&0x2000 == 0x2000,
			CloakedFire:                                     flags2&0x4000 == 0x4000,
			TenfoldDamageToAsteroids:                        flags2&0x8000 == 0x8000,
			AmmoUsedAtEndOfBurstCycle:                       flags3&0x0001 == 0x0001,
			ShotsAreTranslucent:                             flags3&0x0002 == 0x0002,
			OnlyOneShotAtATime:                              flags3&0x0004 == 0x0004,
			ClosestExitPoint:                                flags3&0x0010 == 0x0010,
			ExclusiveWeapon:                                 flags3&0x0020 == 0x0020,
		},
		Seeker: WeapSeeker{
			IgnoresAsteroids:       flagsSeeker&0x0001 == 0x0001,
			DecoyedAsteroids:       flagsSeeker&0x0002 == 0x0002,
			ConfusedByInterference: flagsSeeker&0x0008 == 0x0008,
			TurnsAwayIfJammed:      flagsSeeker&0x0010 == 0x0010,
			NoFireIfIonized:        flagsSeeker&0x0020 == 0x0020,
			NoLockIfNotDirectedAt:  flagsSeeker&0x4000 == 0x4000,
			SelfDamageIfJammed:     flagsSeeker&0x8000 == 0x8000,
		},
		SmokeSet:    CicnID(binary.BigEndian.Uint16(b[32:])),
		Decay:       FrameCount(binary.BigEndian.Uint16(b[34:])),
		Particles:   int16(binary.BigEndian.Uint16(b[36:])),
		PartVel:     int16(binary.BigEndian.Uint16(b[38:])),
		PartLifeMin: int16(binary.BigEndian.Uint16(b[40:])),
		PartLifeMax: int16(binary.BigEndian.Uint16(b[42:])),
		PartColor: color.RGBA{
			A: b[44],
			R: b[45],
			G: b[46],
			B: b[47],
		},
		BeamLength: int16(binary.BigEndian.Uint16(b[48:])),
		BeamWidth:  int16(binary.BigEndian.Uint16(b[50:])),
		Falloff:    int16(binary.BigEndian.Uint16(b[52:])),
		BeamColor: color.RGBA{
			A: b[54],
			R: b[55],
			G: b[56],
			B: b[57],
		},
		CoronaColor: color.RGBA{
			A: b[58],
			R: b[59],
			G: b[60],
			B: b[61],
		},
		SubCount:     int16(binary.BigEndian.Uint16(b[62:])),
		SubType:      WeapID(binary.BigEndian.Uint16(b[64:])),
		SubTheta:     int16(binary.BigEndian.Uint16(b[66:])),
		SubLimit:     int16(binary.BigEndian.Uint16(b[68:])),
		ProxSafety:   int16(binary.BigEndian.Uint16(b[70:])),
		Ionization:   int16(binary.BigEndian.Uint16(b[74:])),
		HitParticles: int16(binary.BigEndian.Uint16(b[76:])),
		HitPartLife:  int16(binary.BigEndian.Uint16(b[78:])),
		HitPartVel:   int16(binary.BigEndian.Uint16(b[80:])),
		HitPartColor: color.RGBA{
			A: b[82],
			R: b[83],
			G: b[84],
			B: b[85],
		},
		ExitType:    WeapExitType(binary.BigEndian.Uint16(b[88:])),
		BurstCount:  int16(binary.BigEndian.Uint16(b[90:])),
		BurstReload: int16(binary.BigEndian.Uint16(b[92:])),
		JamVuln1:    int16(binary.BigEndian.Uint16(b[94:])),
		JamVuln2:    int16(binary.BigEndian.Uint16(b[96:])),
		JamVuln3:    int16(binary.BigEndian.Uint16(b[98:])),
		JamVuln4:    int16(binary.BigEndian.Uint16(b[100:])),
		Durability:  int16(binary.BigEndian.Uint16(b[104:])),
		GuidedTurn:  int16(binary.BigEndian.Uint16(b[106:])),
		MaxAmmo:     int16(binary.BigEndian.Uint16(b[108:])),
		Recoil:      int16(binary.BigEndian.Uint16(b[86:])),
		LiDensity:   int16(binary.BigEndian.Uint16(b[110:])),
		LiAmplitude: int16(binary.BigEndian.Uint16(b[112:])),
		IonizeColor: color.RGBA{
			A: b[114],
			R: b[115],
			G: b[116],
			B: b[117],
		},
	}

	return t
}
