package resources

import (
	"encoding/binary"
	"image"
	"image/color"

	"github.com/imle/resourcefork"
)

// The cölr resource allows you to customize some game-wide interface options.
//
// The various interface buttons that appear are drawn on the fly. Nova uses PICT resources 7500-7502 for the left,
// centre, and right pieces of the "up" buttons, PICT resources 7503-7505 for the "down" button pieces, and PICT
// resources 7506-7508 for the greyed-out button pieces. Corresponding mask images are stored in PICTs 7600-7608.
// STR# resource 150 is used to store the text that appears on each button type.
//
// Note that all colour fields in the cölr resource are encoded the same as HTML colours,
// and that only the first cölr resource is loaded.

type ColrID IDType

type Colr struct {
	ID ColrID

	ButtonUp   color.Color // Normal button text colour.
	ButtonDown color.Color // Pressed button text colour.
	ButtonGrey color.Color // Greyed-out button text colour.

	MenuFont     string      // Main menu font name.
	MenuFontSize int16       // Size of main menu font.
	MenuColor1   color.Color // Bright colour for main menu.
	MenuColor2   color.Color // Dim colour for main menu.

	GridDim    color.Color // Shipyard/outfit dialog grid colour.
	GridBright color.Color // Shipyard/outfit dialog selection square colour.

	ProgressBar image.Rectangle // Position and shape of the loading progress bar, relative to the centre of the window.

	ProgBright  color.Color // Bright progress bar colour.
	ProgDim     color.Color // Dim progress bar colour.
	ProgOutline color.Color // Progress bar outline colour.

	// Position of the six main menu buttons, relative to the top left corner of a 1024x768 main menu background.
	Button1 image.Point
	Button2 image.Point
	Button3 image.Point
	Button4 image.Point
	Button5 image.Point
	Button6 image.Point

	FloatingMap  color.Color // Floating hyperspace map / escort menu border colour.
	ListText     color.Color // List text colour.
	ListBkgnd    color.Color // List background colour.
	ListHilite   color.Color // List hilite colour.
	EscortHilite color.Color // Escort menu item hilite colour.

	ButtonFont   string // Button font name.
	ButtonFontSz int16  // Size of button font.

	// Logo animation x/y position.
	Logo image.Point
	// Rollover animation x/y position.
	Rollover image.Point
	// Sliding button x/y positions.
	Slide1 image.Point
	Slide2 image.Point
	Slide3 image.Point
}

func ColrFromResource(resource resourcefork.Resource) *Colr {
	return ColrFromBytes(ColrID(resource.ID), resource.Data)
}

func ColrFromBytes(id ColrID, b []byte) *Colr {
	t := &Colr{
		ID:           id,
		ButtonUp:     color.RGBA{A: b[0], R: b[1], G: b[2], B: b[3]},
		ButtonDown:   color.RGBA{A: b[4], R: b[5], G: b[6], B: b[7]},
		ButtonGrey:   color.RGBA{A: b[8], R: b[9], G: b[10], B: b[11]},
		MenuFont:     byteString(b[12:], 63),
		MenuFontSize: int16(binary.BigEndian.Uint16(b[76:])),
		MenuColor1:   color.RGBA{A: b[78], R: b[79], G: b[80], B: b[81]},
		MenuColor2:   color.RGBA{A: b[82], R: b[83], G: b[84], B: b[85]},
		GridBright:   color.RGBA{A: b[86], R: b[87], G: b[88], B: b[89]},
		GridDim:      color.RGBA{A: b[90], R: b[91], G: b[92], B: b[93]},
		ProgressBar: image.Rect(
			int(int16(binary.BigEndian.Uint16(b[94:]))),
			int(int16(binary.BigEndian.Uint16(b[96:]))),
			int(int16(binary.BigEndian.Uint16(b[98:]))),
			int(int16(binary.BigEndian.Uint16(b[100:]))),
		),
		ProgBright:   color.RGBA{A: b[102], R: b[103], G: b[104], B: b[105]},
		ProgDim:      color.RGBA{A: b[106], R: b[107], G: b[108], B: b[109]},
		ProgOutline:  color.RGBA{A: b[110], R: b[111], G: b[112], B: b[113]},
		Button1:      image.Point{X: int(b[114]), Y: int(b[116])},
		Button2:      image.Point{X: int(b[118]), Y: int(b[120])},
		Button3:      image.Point{X: int(b[122]), Y: int(b[124])},
		Button4:      image.Point{X: int(b[126]), Y: int(b[128])},
		Button5:      image.Point{X: int(b[130]), Y: int(b[132])},
		Button6:      image.Point{X: int(b[134]), Y: int(b[136])},
		FloatingMap:  color.RGBA{A: b[138], R: b[139], G: b[140], B: b[141]},
		ListText:     color.RGBA{A: b[142], R: b[143], G: b[144], B: b[145]},
		ListBkgnd:    color.RGBA{A: b[146], R: b[147], G: b[148], B: b[149]},
		ListHilite:   color.RGBA{A: b[150], R: b[151], G: b[152], B: b[153]},
		EscortHilite: color.RGBA{A: b[154], R: b[155], G: b[156], B: b[157]},
		ButtonFont:   byteString(b[158:], 63),
		ButtonFontSz: int16(binary.BigEndian.Uint16(b[222:])),
		Logo:         image.Point{X: int(b[224]), Y: int(b[226])},
		Rollover:     image.Point{X: int(b[228]), Y: int(b[230])},
		Slide1:       image.Point{X: int(b[232]), Y: int(b[234])},
		Slide2:       image.Point{X: int(b[236]), Y: int(b[238])},
		Slide3:       image.Point{X: int(b[240]), Y: int(b[242])},
	}

	return t
}
