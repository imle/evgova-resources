package resources

import (
	"encoding/binary"
	"image"
	"image/color"

	"github.com/imle/resourcefork"
)

// The ïntf resource controls the appearance of the status bar by modifying the position and colour of the various
// status bar elements, as well as changing the status bar background image. Additionally, the use of multiple ïntf
// resources, along with proper values in the various gövt resources' Interface fields, allows the appearance of
// status bar to change Intfd on what type of ship the player is piloting - this is a useless but fairly neat effect.

type IntfID IDType

type Intf struct {
	ID IntfID

	BrightText   color.Color     // Bright text colour.
	DimText      color.Color     // Dim text colour.
	RadarArea    image.Rectangle // Rectangular bounds of the radar display.
	BrightRadar  color.Color     // Bright radar pixel colour.
	DimRadar     color.Color     // Dim radar pixel colour (note that having an IFF outfit will override these colours).
	ShieldArea   image.Rectangle // Rectangular bounds of the shield indicator.
	ShieldColor  color.Color     // Shield bar colour.
	ArmorArea    image.Rectangle // Rectangular bounds of the armour indicator.
	ArmorColor   color.Color     // Armour bar colour.
	FuelArea     image.Rectangle // Rectangular bounds of the fuel indicator.
	FuelFull     color.Color     // Colour of the "full jumps" portion of the fuel indicator.
	FuelPartial  color.Color     // Colour of the "partial fuel" portion of the fuel indicator.
	NavArea      image.Rectangle // Rectangular bounds of the navigation display.
	WeapArea     image.Rectangle // Rectangular bounds of the weapon display.
	TargArea     image.Rectangle // Rectangular bounds of the target display.
	CargoArea    image.Rectangle // Rectangular bounds of the cargo display.
	StatusFont   string          // Font to use for the status bar.
	StatFontSize int16           // Normal font size to use.
	SubtitleSize int16           // Font size for ship subtitles.
	StatusBkgnd  PictID          // ID of PICT resource to use as backdrop for status display. Values less than 128 are interpreted as 128.
}

func IntfFromResource(resource resourcefork.Resource) *Intf {
	return IntfFromBytes(IntfID(resource.ID), resource.Data)
}

func IntfFromBytes(id IntfID, b []byte) *Intf {
	t := &Intf{
		ID:         id,
		BrightText: color.RGBA{A: b[0], R: b[1], G: b[2], B: b[3]},
		DimText:    color.RGBA{A: b[4], R: b[5], G: b[6], B: b[7]},
		RadarArea: image.Rect(
			int(int16(binary.BigEndian.Uint16(b[8:]))),
			int(int16(binary.BigEndian.Uint16(b[10:]))),
			int(int16(binary.BigEndian.Uint16(b[12:]))),
			int(int16(binary.BigEndian.Uint16(b[14:]))),
		),
		BrightRadar: color.RGBA{A: b[16], R: b[17], G: b[18], B: b[19]},
		DimRadar:    color.RGBA{A: b[20], R: b[21], G: b[22], B: b[23]},
		ShieldArea: image.Rect(
			int(int16(binary.BigEndian.Uint16(b[24:]))),
			int(int16(binary.BigEndian.Uint16(b[26:]))),
			int(int16(binary.BigEndian.Uint16(b[28:]))),
			int(int16(binary.BigEndian.Uint16(b[30:]))),
		),
		ShieldColor: color.RGBA{A: b[32], R: b[33], G: b[34], B: b[35]},
		ArmorArea: image.Rect(
			int(int16(binary.BigEndian.Uint16(b[36:]))),
			int(int16(binary.BigEndian.Uint16(b[38:]))),
			int(int16(binary.BigEndian.Uint16(b[40:]))),
			int(int16(binary.BigEndian.Uint16(b[42:]))),
		),
		ArmorColor: color.RGBA{A: b[44], R: b[45], G: b[46], B: b[47]},
		FuelArea: image.Rect(
			int(int16(binary.BigEndian.Uint16(b[48:]))),
			int(int16(binary.BigEndian.Uint16(b[50:]))),
			int(int16(binary.BigEndian.Uint16(b[52:]))),
			int(int16(binary.BigEndian.Uint16(b[54:]))),
		),
		FuelFull:    color.RGBA{A: b[56], R: b[57], G: b[58], B: b[59]},
		FuelPartial: color.RGBA{A: b[60], R: b[61], G: b[62], B: b[63]},
		NavArea: image.Rect(
			int(int16(binary.BigEndian.Uint16(b[64:]))),
			int(int16(binary.BigEndian.Uint16(b[66:]))),
			int(int16(binary.BigEndian.Uint16(b[68:]))),
			int(int16(binary.BigEndian.Uint16(b[70:]))),
		),
		WeapArea: image.Rect(
			int(int16(binary.BigEndian.Uint16(b[72:]))),
			int(int16(binary.BigEndian.Uint16(b[74:]))),
			int(int16(binary.BigEndian.Uint16(b[76:]))),
			int(int16(binary.BigEndian.Uint16(b[78:]))),
		),
		TargArea: image.Rect(
			int(int16(binary.BigEndian.Uint16(b[80:]))),
			int(int16(binary.BigEndian.Uint16(b[82:]))),
			int(int16(binary.BigEndian.Uint16(b[84:]))),
			int(int16(binary.BigEndian.Uint16(b[86:]))),
		),
		CargoArea: image.Rect(
			int(int16(binary.BigEndian.Uint16(b[88:]))),
			int(int16(binary.BigEndian.Uint16(b[90:]))),
			int(int16(binary.BigEndian.Uint16(b[92:]))),
			int(int16(binary.BigEndian.Uint16(b[94:]))),
		),
		StatusFont:   byteString(b[96:], 63),
		StatFontSize: int16(binary.BigEndian.Uint16(b[160:])),
		SubtitleSize: int16(binary.BigEndian.Uint16(b[162:])),
		StatusBkgnd:  PictID(binary.BigEndian.Uint16(b[164:])),
	}

	return t
}
