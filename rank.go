package resources

import (
	"encoding/binary"

	"github.com/imle/resourcefork"
)

// The rank resource is used to give the player a feeling of 'belonging' to a given government. It can also be used to
// give the player certain advantages that come with rank. When a rank is made active (which is accomplished through
// any suitable control bit set string) the player is given all the privileges of that rank, whatever they might be,
// and the name of that rank is displayed in the player-info dialog.

type RankID IDType

type RankFlags struct {
	DeactivateAllOtherForGovtOnActivate   bool // 0x0001 Deactivate all other active ranks affiliated with this same govt when this rank is activated (excludes permanent ranks).
	DeactivateAllOtherForGovtOnDeactivate bool // 0x0002 Deactivate all other active ranks affiliated with this same govt when this rank is deactivated (excludes permanent ranks).
	DeactivateOnGovtShipDisableOrDestroy  bool // 0x0004 Deactivate this rank if player destroys or disables a ship of the affiliated government or its allies.
	Permanent                             bool // 0x0008 Rank is permanent and cannot be deactivated except if explicitly done by a control bit eval string.
	DeactivateAllLowerForGovtOnActivate   bool // 0x0010 Deactivate all other active and lower-weighted ranks affiliated with this same govt when this rank is activated (excludes permanent ranks).
	DeactivateAllLowerForGovtOnDeactivate bool // 0x0020 Deactivate all other active and lower-weighted ranks affiliated with this same govt when this rank is deactivated (excludes permanent ranks).
	DeactivateOnCrimeAgainstGovt          bool // 0x0040 Deactivate this rank if the player commits any crime against the affiliated government.
	GovtShipsWillNotAutoAttack            bool // 0x0100 Ships of the affiliated government will not automatically attack the player when he has this rank.
	GovtPlanetsAlwaysAllowDock            bool // 0x0200 All planets of the affiliated government will let the player land when he has this rank, regardless of their MinStatus field.
	CanRequestGovtAssistance              bool // 0x0400 Player can always request battle assistance from ships of the affiliated government, who will also call in reinforcements on the player's behalf if they are available.
	FreeRepairAndFuelFromGovt             bool // 0x0800 Ships allied with the affiliated govt will always repair or refuel the player for free.
}

type Rank struct {
	ID RankID

	Weight     int16
	AffilGovt  GovtID
	Contribute FlagMask64
	Salary     int16
	SalaryCap  int16
	Flags      RankFlags
	PriceMod   int16
	ConvName   string
	ShortName  string
}

func RankFromResource(resource resourcefork.Resource) *Rank {
	return RankFromBytes(RankID(resource.ID), resource.Data)
}

func RankFromBytes(id RankID, b []byte) *Rank {
	flags1 := binary.BigEndian.Uint16(b[22:])

	t := &Rank{
		ID:         id,
		Weight:     int16(binary.BigEndian.Uint16(b[0:])),
		AffilGovt:  GovtID(binary.BigEndian.Uint16(b[2:])),
		PriceMod:   int16(binary.BigEndian.Uint16(b[4:])),
		SalaryCap:  int16(binary.BigEndian.Uint16(b[6:])),
		Salary:     int16(binary.BigEndian.Uint16(b[8:])),
		Contribute: FlagMask64(binary.BigEndian.Uint64(b[14:])),
		Flags: RankFlags{
			DeactivateAllOtherForGovtOnActivate:   flags1&0x0001 == 0x0001,
			DeactivateAllOtherForGovtOnDeactivate: flags1&0x0002 == 0x0002,
			DeactivateOnGovtShipDisableOrDestroy:  flags1&0x0004 == 0x0004,
			Permanent:                             flags1&0x0008 == 0x0008,
			DeactivateAllLowerForGovtOnActivate:   flags1&0x0010 == 0x0010,
			DeactivateAllLowerForGovtOnDeactivate: flags1&0x0020 == 0x0020,
			DeactivateOnCrimeAgainstGovt:          flags1&0x0040 == 0x0040,
			GovtShipsWillNotAutoAttack:            flags1&0x0100 == 0x0100,
			GovtPlanetsAlwaysAllowDock:            flags1&0x0200 == 0x0200,
			CanRequestGovtAssistance:              flags1&0x0400 == 0x0400,
			FreeRepairAndFuelFromGovt:             flags1&0x0800 == 0x0800,
		},
		ConvName:  byteString(b[24:], 63),
		ShortName: byteString(b[88:], 63),
	}

	return t
}
