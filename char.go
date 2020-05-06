package resources

import (
	"encoding/binary"
	"time"

	"github.com/imle/resourcefork"
)

// The chär resource is used to allow multiple entry points into the scenario's storyline, by letting the player
// pick a "character template" from a list when a new pilot is started. The player is presented with a list of
// all available chär resources when a new pilot is created, and must choose one (but only one). This allows
// different character types to have different ships, legal records, etc. at the start of the game, and also allows
// for the setting-up of mission stuff by way of a control bit set string that is evaluated when a new pilot is started.

type CharID IDType

type CharFlags struct {
	Default bool
}

type Char struct {
	ID CharID

	Cash Credits // The amount of money a player gets when starting out with this character type.

	ShipType ShipID // ID number of the starting ship type.

	System [4]SystID // ID numbers of up to four possible starting systems for the player. The player will randomly be placed in one of these systems when starting out. Set to -1 if unused (if all four of these fields are set to -1, the player will be placed in system ID 128 as a default).

	Govt   [4]GovtID // For each of the governments whose ID is entered in a Govt1-4 field, the player's legal status in systems owned by that government or one of its allies is set to the value in the corresponding Status field. For systems owned by an enemy of the govt identified in the Govt field, the player's legal status is set to the negative value of the number in the corresponding Status field. Set unused Govt fields to -1.
	Status [4]Status

	CombatRating int16 // The player's starting combat rating.

	IntroPict [4]PictID // IDs of up to four PICT resources to show in sequence when the player starts out with a new pilot of this type. Set to -1 if unused.

	PictDelay [4]time.Duration // Maximum delay time to display each of the four above pictures, in seconds.

	IntroTextID DescID // The ID of the dësc resource to show when the player starts out with a new pilot of this type. Along with the IntroPictID field, this allow you to have different opening sequences for each pilot type (useful to show different sides of the same issue, for example). Using the dësc resource's MovieFile field here lets you have an introduction movie. The intro text (and any associated movie) is displayed after the above PICT resources are shown.

	OnStart ControlBitFunction // A control bit set string that is called when the player starts out with a new pilot of this type.

	Flags CharFlags // Denotes the "default" chär resource, which will be automatically selected in the popup menu. If more than one chär resource has this bit set (there shouldn't be) the one with the lowest ID will be considered the default.

	StartDate  int16  // The starting day of the game.
	StartMonth int16  // The starting month of the game.
	StartYear  int16  // The starting year of the game.
	DatePrefix string // String that is appended to the start of the date whenever it's displayed.
	DateSuffix string // String that is appended to the end of the date whenever it's displayed.
}

func CharFromResource(resource resourcefork.Resource) *Char {
	return CharFromBytes(CharID(resource.ID), resource.Data)
}

func CharFromBytes(id CharID, b []byte) *Char {
	flags1 := binary.BigEndian.Uint16(b[306:])

	t := &Char{
		ID:       id,
		Cash:     Credits(binary.BigEndian.Uint32(b[0:])),
		ShipType: ShipID(binary.BigEndian.Uint16(b[4:])),
		System: [4]SystID{
			SystID(binary.BigEndian.Uint16(b[6:])),
			SystID(binary.BigEndian.Uint16(b[8:])),
			SystID(binary.BigEndian.Uint16(b[10:])),
			SystID(binary.BigEndian.Uint16(b[12:])),
		},
		Govt: [4]GovtID{
			GovtID(binary.BigEndian.Uint16(b[14:])),
			GovtID(binary.BigEndian.Uint16(b[16:])),
			GovtID(binary.BigEndian.Uint16(b[18:])),
			GovtID(binary.BigEndian.Uint16(b[20:])),
		},
		Status: [4]Status{
			Status(binary.BigEndian.Uint16(b[22:])),
			Status(binary.BigEndian.Uint16(b[24:])),
			Status(binary.BigEndian.Uint16(b[26:])),
			Status(binary.BigEndian.Uint16(b[28:])),
		},
		CombatRating: int16(binary.BigEndian.Uint16(b[30:])),
		IntroPict: [4]PictID{
			PictID(binary.BigEndian.Uint16(b[32:])),
			PictID(binary.BigEndian.Uint16(b[34:])),
			PictID(binary.BigEndian.Uint16(b[36:])),
			PictID(binary.BigEndian.Uint16(b[38:])),
		},
		PictDelay: [4]time.Duration{
			time.Duration(binary.BigEndian.Uint16(b[40:])) * time.Second,
			time.Duration(binary.BigEndian.Uint16(b[42:])) * time.Second,
			time.Duration(binary.BigEndian.Uint16(b[44:])) * time.Second,
			time.Duration(binary.BigEndian.Uint16(b[46:])) * time.Second,
		},
		IntroTextID: DescID(binary.BigEndian.Uint16(b[48:])),
		OnStart:     ControlBitFunction(byteString(b[50:], 255)),
		Flags: CharFlags{
			Default: flags1&0x0001 == 0x0001,
		},
		StartDate:  int16(binary.BigEndian.Uint16(b[308:])),
		StartMonth: int16(binary.BigEndian.Uint16(b[310:])),
		StartYear:  int16(binary.BigEndian.Uint16(b[312:])),
		DatePrefix: byteString(b[314:], 15),
		DateSuffix: byteString(b[330:], 15),
	}

	return t
}
