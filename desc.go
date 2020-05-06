package resources

import (
	"bytes"
	"encoding/binary"

	"github.com/imle/resourcefork"
)

// Desc resources store null-terminated text strings (descriptions) that are used by Nova in a variety of places. For
// some desc resources, Nova looks for a certain reserved ID number. Other desc resources are pointed to by fields in
// other resources, so their ID numbers are not necessarily fixed, and can be set to virtually anything by the scenario
// designer. The reserved desc ID numbers, along with the maximum length for each type, are below:
//
// 128-2175       Stellar object descriptions, shown when landed on a planet.
// 3000-3511      Outfit item descriptions, shown in ship outfitting dialog.
// 4000-4999      Mission descriptions, shown in mission dialog.
// 13000-13767    Ship class descriptions, shown in the shipyard and requisition-escort dialog.
// 13999          Message shown after the player uses an escape pod.
// 14000-14767    Ship pilot descriptions, shown in the hire-escort dialog.
// 32760-32767    Reserved.
//
// If you wish, you can make a dësc resource mutable via control bits - embedding a special sequence of characters into
// the dësc resource will instruct Nova to change the contents of the text on the fly. This sequence is delimited
// (marked) by the characters "{" and "}", and follows this format:
// >> {bXXX "string one" "string two"}
// Where "XXX" is be replaced by the index of the control bit you wish to test. You can add in a "!" character before
// the "bXXX" test in order to negate the result of the test, but unlike the control bit test strings, you cannot
// perform compound tests in a dësc resource - i.e., no testing of multiple bits at a time.
//
// If the bit test (after being negated, if the "!" character is present) evaluates to true, the first string will be
// substituted in place of all the characters between (and including) the "{" and "}" characters. If the bit test
// evaluates to false and there is a second string in the expression, that second string will be substituted. If
// there is no second string, nothing will be substituted. For example, consider this dësc resource:
// >> This is a {b001 "great and terrific" "lousy, terrible"} example.
// ...if bit 001 is set, the output will be "This is a great and terrific example." If bit 001 is not set, the output
// will be "This is a lousy, terrible example.".
//
// Also note that if you want to include a quotation mark (") character in either of the two strings,
// use standard C syntax to do it:
// >> My name is {b002 "Dave \"pipeline\" Williams"}
//
// This is also works with the player's gender - for example:
// >> This is a test string and the player is {G "a male character" "a female pilot"}.
// ...in this case, the G character signifies that the following text is mutable based on the player's gender;
// if the player is male, the first string is used, otherwise the second string is used.
// Note that the "!" token works here as usual.
//
// You can also change the text based on whether or not the player is registered:
// >> This is a test string you {P "have paid" "haven't paid"}.
// ...in this case, the P character signifies that the following text is mutable based on whether or not the game is
// registered; if the player has registered, the first string is used, otherwise the second string is used. Note that
// the "!" token works here as usual, and you can also append a number to the P character to specify a number of days,
// just as you can with the ncb Pxxx test operator.

type DescID IDType

const (
	DescIDOffsetStellar   DescID = 128
	DescIDOffsetOutfit    DescID = 3000
	DescIDOffsetMission   DescID = 4000
	DescIDOffsetShipClass DescID = 13000
	DescIDOffsetEscapePod DescID = 13999
	DescIDOffsetPilot     DescID = 14000
	DescIDOffsetReserved  DescID = 32760
)

type DescFlags struct {
	MovieAfterBriefing bool // 0x0001 Show the movie after the briefing instead of before.
	MovieDoubleSize    bool // 0x0002 Show movie at double-size.
	CinematicMovie     bool // 0x0004 Cinematic movie - blank background and fade screen before and after .
}

type OperatedString string

type Desc struct {
	ID          DescID
	Description OperatedString

	// This is used to include graphics in mission briefings. If you put in the ID of a valid PICT resource,
	// Nova will display that image along with the description text when it displays a mission dialog box
	// (with the exception of the Mission Computer and Mission Info dialogs).
	Graphic PictID

	// The name of a QuickTime movie file to display before the briefing dialog appears.
	// This file must reside either in "Nova Files" or "Nova Plugins".
	MovieFile string

	Flags DescFlags
}

func DescFromResource(resource resourcefork.Resource) *Desc {
	return DescFromBytes(DescID(resource.ID), resource.Data)
}

const (
	offsetMovie   = 3
	offsetGraphic = 1
)

func DescFromBytes(id DescID, b []byte) *Desc {
	descEnd := bytes.IndexByte(b, 0)
	movStart := descEnd + offsetMovie

	flags := b[36]

	t := &Desc{
		ID:          id,
		Description: OperatedString(bytes.Runes(b[:descEnd])),
		Graphic:     PictID(binary.BigEndian.Uint16(b[descEnd+offsetGraphic:])),
		MovieFile:   string(bytes.Runes(b[movStart : movStart+bytes.IndexByte(b[movStart:], 0)])),
		Flags: DescFlags{
			MovieAfterBriefing: flags&0x0001 == 0x0001,
			MovieDoubleSize:    flags&0x0002 == 0x0002,
			CinematicMovie:     flags&0x0004 == 0x0004,
		},
	}

	return t
}
