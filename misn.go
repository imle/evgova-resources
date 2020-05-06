package resources

import (
	"encoding/binary"

	"github.com/imle/resourcefork"
)

// Missions are the crown jewel of the Nova datafile, as well as the largest and most complex resources in the game.
// Each misn resource corresponds to a single mission that the player can undertake, with the name of the mission
// (which the player sees in the mission list) being the name of the associated misn resource.

type MisnID IDType

type MisnAvailLoc int16

const (
	MisnAvailLocMissionComputer MisnAvailLoc = iota // 0 - From the mission computer.
	MisnAvailLocBar                                 // 1 - In the bar.
	MisnAvailLocShip                                // 2 - Offered from ship (must set up associated përs resource as well).
	MisnAvailLocMainSpaceport                       // 3 - In the main spaceport dialog.
	MisnAvailLocTrading                             // 4 - In the trading dialog.
	MisnAvailLocShipyard                            // 5 - In the shipyard dialog.
	MisnAvailLocOutfit                              // 6 - In the outfit dialog.
)

type MisnShipGoal int16

const (
	MisnShipGoalDestroy  MisnShipGoal = iota // 0 - Destroy all the ships.
	MisnShipGoalDisable                      // 1 - Disable but don't destroy them.
	MisnShipGoalBoardAll                     // 2 - Board them.
	MisnShipGoalEscort                       // 3 - Escort them (keep them from getting killed).
	MisnShipGoalObserve                      // 4 - Observe them (for ships that can cloak, at least one must be visible onscreen - for ships that cannot cloak, the player must merely be in the same system as them).
	MisnShipGoalRescue                       // 5 - Rescue them (they start out disabled and stay that way until you board them; to make a rescue mission where the ship stays disabled, give the special ships a govt with the "always disabled" flag set).
	MisnShipGoalChaseOff                     // 6 - Chase them off (either kill them or scare the into jumping out of the system).
)

type MisnPickupGoal int16

const (
	MisnPickupGoalIgnore           MisnPickupGoal = iota - 1 // -1 - Ignored.
	MisnPickupGoalAtStart                                    // 0  - Pick up at mission start.
	MisnPickupGoalAtTravelStel                               // 1  - Pick up at TravelStel.
	MisnPickupGoalBoardSpecialShip                           // 2  - Pick up when boarding special ship.
)

type MisnDropOffMode int16

const (
	MisnDropOffModeIgnore       MisnDropOffMode = iota - 1 // -1 - Ignored.
	MisnDropOffModeAtTravelStel                            // 0  - Drop off at TravelStel.
	MisnDropOffModeAtReturnStel                            // 1  - Drop off at mission end (ReturnStel).
)

type MisnShipBehav int16

const (
	MisnShipBehavIgnore              MisnShipBehav = iota - 1 // -1 - Ignored (they use their standard AI routines).
	MisnShipBehavAlwaysAttackPlayer                           // 0  - Special ships will always attack the player.
	MisnShipBehavProtectPlayer                                // 1  - Special ships will protect the player.
	MisnShipBehavAttackEnemyStellars                          // 2  - Special ships will attempt to destroy enemy stellars.
)

type MisnShipStart int16

const (
	MisnShipStartDefaultNav4     MisnShipStart = iota - 4 // Appear on top of nav default 4.
	MisnShipStartDefaultNav3                              // Appear on top of nav default 3.
	MisnShipStartDefaultNav2                              // Appear on top of nav default 2.
	MisnShipStartDefaultNav1                              // Appear on top of nav default 1.
	MisnShipStartRandomly                                 // Appear randomly in the system (as usual).
	MisnShipStartJumpIn                                   // Jump in from hyperspace after a short delay.
	MisnShipStartRandomlyCloaked                          // Appear randomly, cloaked.
)

type MisnFlags struct {
	AutoAbort                  bool // 0x0001 Marks the mission as an auto-aborting mission, which will automatically abort itself after it is accepted. (sometimes useful to create special ships) Any control bits pointed to by the mission's CompBitSet fields will be automatically set when the mission aborts.
	NoDestinationArrowsOnMap   bool // 0x0002 Don't show the red destination arrows on the map.
	CannotRefuse               bool // 0x0004 Can't refuse the mission.
	Take100FuelOnAutoAbort     bool // 0x0008 Mission takes away 100 units of fuel upon auto-abort. (mission won't be offered if player has less than 100 units of fuel).
	InfiniteAuxShips           bool // 0x0010 Infinite auxShips.
	FailIfScanned              bool // 0x0020 Mission fails if you're scanned.
	Negative5CompRewardOnAbort bool // 0x0040 Apply -5x CompReward reversal on abort.
	GlobalPenaltyIfJettisoning bool // 0x0080 Global penalty when jettisoning mission cargo in space (currently ignored).
	ShowGreenArrowInitialBrief bool // 0x0100 Show green arrow on map in initial briefing.
	ShowArrowForShipSyst       bool // 0x0200 Show an additional arrow on the map for the ShipSyst.
	InvisibleInMissionDialog   bool // 0x0400 Mission is invisible and won't appear in the mission info dialog. (be careful with this!).
	RandomSpecialShipLocks     bool // 0x0800 The special ships' type will be selected at mission start and then kept the same whenever the special ships for that mission are created, until the mission ends. This can be used for (e.g.) "attack pirate" missions where you want the type of the enemy ship to be random at first but you don't want it to change every time the player lands or re-enters the system.
	UnavailableForFreighters   bool // 0x2000 Mission unavailable if player's ship is of inherentAI type 1 or 2 (cargo ships).
	UnavailableForWarships     bool // 0x4000 Mission unavailable if player's ship is of inherentAI type 3 or 4 (warships).
	FailIfBoardByPirates       bool // 0x8000 Mission will fail if player is boarded by pirates.

	UnavailableIfNotEnoughSpace bool // 0x0001 Don't offer mission if the player doesn't have enough cargo space to hold the mission cargo (even if the mission cargo won't be picked up until later).
	PayOnAutoAbort              bool // 0x0002 Apply mission Pay on auto-abort.
	FailIfDisabledOrDestroyed   bool // 0x0004 Mission fails if player is disabled or destroyed.

	CanAbort bool // Not part of actual flag-bits, but is flag
}

type Payment int32

type Misn struct {
	ID            MisnID
	AvailStel     int16
	AvailLoc      MisnAvailLoc
	AvailRecord   int16
	AvailRating   int16
	AvailRandom   int16
	TravelStel    int16
	ReturnStel    int16
	CargoType     int16
	CargoQty      int16
	PickupMode    MisnPickupGoal
	DropOffMode   MisnDropOffMode
	ScanMask      FlagMask16
	PayVal        Payment
	ShipCount     int16
	ShipSyst      int16
	ShipDude      DudeID
	ShipGoal      MisnShipGoal
	ShipBehav     MisnShipBehav
	ShipNameID    StrAID
	ShipStart     MisnShipStart
	CompGovt      GovtID
	CompReward    int16
	ShipSubtitle  StrAID
	BriefText     DescID
	QuickBrief    DescID
	LoadCargText  DescID
	DumpCargoText DescID
	CompText      DescID
	FailText      DescID
	ShipDoneText  DescID
	TimeLimit     int16
	AuxShipCount  int16
	AuxShipDude   DudeID
	AuxShipSyst   int16
	Flags         MisnFlags
	AvailShipType int16
	RefuseText    DescID
	AvailBits     ControlBitTest
	OnAccept      ControlBitFunction
	OnRefuse      ControlBitFunction
	OnSuccess     ControlBitFunction
	OnFailure     ControlBitFunction
	OnAbort       ControlBitFunction
	OnShipDone    ControlBitFunction
	Require       FlagMask64
	DatePostInc   int16
	AcceptButton  string
	RefuseButton  string
	DispWeight    int16
}

func (m Misn) DescID() DescID {
	return DescID(m.ID) - resourcefork.ResourceForkIDOffset + DescIDOffsetMission
}

func MisnFromResource(resource resourcefork.Resource) *Misn {
	return MisnFromBytes(MisnID(resource.ID), resource.Data)
}

func MisnFromBytes(id MisnID, b []byte) *Misn {
	flags := binary.BigEndian.Uint32(b[80:])
	flags2 := binary.BigEndian.Uint32(b[82:])

	t := &Misn{
		ID:            id,
		AvailStel:     int16(binary.BigEndian.Uint16(b[0:])),
		AvailLoc:      MisnAvailLoc(binary.BigEndian.Uint16(b[4:])),
		AvailRecord:   int16(binary.BigEndian.Uint16(b[6:])),
		AvailRating:   int16(binary.BigEndian.Uint16(b[8:])),
		AvailRandom:   int16(binary.BigEndian.Uint16(b[10:])),
		TravelStel:    int16(binary.BigEndian.Uint16(b[12:])),
		ReturnStel:    int16(binary.BigEndian.Uint16(b[14:])),
		CargoType:     int16(binary.BigEndian.Uint16(b[16:])),
		CargoQty:      int16(binary.BigEndian.Uint16(b[18:])),
		PickupMode:    MisnPickupGoal(binary.BigEndian.Uint16(b[20:])),
		DropOffMode:   MisnDropOffMode(binary.BigEndian.Uint16(b[22:])),
		ScanMask:      FlagMask16(binary.BigEndian.Uint16(b[24:])),
		PayVal:        Payment(binary.BigEndian.Uint32(b[30:])),
		ShipCount:     int16(binary.BigEndian.Uint16(b[32:])),
		ShipSyst:      int16(binary.BigEndian.Uint16(b[34:])),
		ShipDude:      DudeID(binary.BigEndian.Uint16(b[36:])),
		ShipGoal:      MisnShipGoal(binary.BigEndian.Uint16(b[38:])),
		ShipBehav:     MisnShipBehav(binary.BigEndian.Uint16(b[40:])),
		ShipNameID:    StrAID(binary.BigEndian.Uint16(b[42:])),
		ShipStart:     MisnShipStart(binary.BigEndian.Uint16(b[44:])),
		CompGovt:      GovtID(binary.BigEndian.Uint16(b[46:])),
		CompReward:    int16(binary.BigEndian.Uint16(b[48:])),
		ShipSubtitle:  StrAID(binary.BigEndian.Uint16(b[50:])),
		BriefText:     DescID(binary.BigEndian.Uint16(b[52:])),
		QuickBrief:    DescID(binary.BigEndian.Uint16(b[54:])),
		LoadCargText:  DescID(binary.BigEndian.Uint16(b[56:])),
		DumpCargoText: DescID(binary.BigEndian.Uint16(b[58:])),
		CompText:      DescID(binary.BigEndian.Uint16(b[60:])),
		FailText:      DescID(binary.BigEndian.Uint16(b[62:])),
		TimeLimit:     int16(binary.BigEndian.Uint16(b[64:])),
		ShipDoneText:  DescID(binary.BigEndian.Uint16(b[68:])),
		AuxShipCount:  int16(binary.BigEndian.Uint16(b[72:])),
		AuxShipDude:   DudeID(binary.BigEndian.Uint16(b[74:])),
		AuxShipSyst:   int16(binary.BigEndian.Uint16(b[76:])),
		Flags: MisnFlags{
			CanAbort:                    int16(binary.BigEndian.Uint16(b[66:])) == 1,
			AutoAbort:                   flags&0x0001 == 0x0001,
			NoDestinationArrowsOnMap:    flags&0x0002 == 0x0002,
			CannotRefuse:                flags&0x0004 == 0x0004,
			Take100FuelOnAutoAbort:      flags&0x0008 == 0x0008,
			InfiniteAuxShips:            flags&0x0010 == 0x0010,
			FailIfScanned:               flags&0x0020 == 0x0020,
			Negative5CompRewardOnAbort:  flags&0x0040 == 0x0040,
			GlobalPenaltyIfJettisoning:  flags&0x0080 == 0x0080,
			ShowGreenArrowInitialBrief:  flags&0x0100 == 0x0100,
			ShowArrowForShipSyst:        flags&0x0200 == 0x0200,
			InvisibleInMissionDialog:    flags&0x0400 == 0x0400,
			RandomSpecialShipLocks:      flags&0x0800 == 0x0800,
			UnavailableForFreighters:    flags&0x2000 == 0x2000,
			UnavailableForWarships:      flags&0x4000 == 0x4000,
			FailIfBoardByPirates:        flags&0x8000 == 0x8000,
			UnavailableIfNotEnoughSpace: flags2&0x0001 == 0x0001,
			PayOnAutoAbort:              flags2&0x0002 == 0x0002,
			FailIfDisabledOrDestroyed:   flags2&0x0004 == 0x0004,
		},
		AvailShipType: int16(binary.BigEndian.Uint16(b[90:])),
		RefuseText:    DescID(binary.BigEndian.Uint16(b[88:])),
		AvailBits:     ControlBitTest(byteString(b[92:], 254)),
		OnAccept:      ControlBitFunction(byteString(b[347:], 254)),
		OnRefuse:      ControlBitFunction(byteString(b[602:], 254)),
		OnSuccess:     ControlBitFunction(byteString(b[857:], 254)),
		OnFailure:     ControlBitFunction(byteString(b[1112:], 254)),
		OnAbort:       ControlBitFunction(byteString(b[1367:], 254)),
		OnShipDone:    ControlBitFunction(byteString(b[1632:], 254)),
		Require:       FlagMask64(binary.BigEndian.Uint64(b[1622:])),
		DatePostInc:   int16(binary.BigEndian.Uint16(b[1631:])),
		AcceptButton:  byteString(b[1887:], 31),
		RefuseButton:  byteString(b[1919:], 31),
		DispWeight:    int16(binary.BigEndian.Uint16(b[1952:])),
	}

	return t
}

func (m Misn) OperateCompText() {
	panic("implement me!")
}

func (m Misn) OperateQuickBrief() {
	panic("implement me!")
}
