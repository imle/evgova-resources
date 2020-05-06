package resources

import (
	"encoding/binary"
	"time"

	"github.com/imle/resourcefork"
)

// Shan (ship animation) resources contain sprite info for ship graphics,
// which are too complex for the more rudimentary spin resource.

type ShanID IDType

type BlinkMode int16

const (
	BlinkModeNone BlinkMode = iota - 1
	BlinkModeIgnored
	BlinkModeSquareWave
	BlinkModeTriangleWave
	BlinkModeRandomPulsing
)

type ShanFlags struct {
	Banking                bool // 0x0001 Extra frames in base image are used to display banking. The first set of sprites is used for level flight, the second for banking left, and the third for banking right.
	AnimatedParts          bool // 0x0002 Extra frames in base image are used for animated ship parts such as for folding/unfolding wings. The sprites will be cycled upon landing, taking off, and entering/exiting hyperspace.
	NoKeyCarried           bool // 0x0004 The second set of frames in the base image are displayed when the ship is not carrying any of its KeyCarried type ships onboard.
	Sequence               bool // 0x0008 Extra frames in base image are shown in sequence, just like the sprites in the alternating image. The AnimDelay field has the same effect in this case.
	StopAnimationsDisabled bool // 0x0010 Stop the ships' animations when it is disabled.
	HideAltDisabled        bool // 0x0020 Hide alt sprites when the ship is disabled.
	HideLightsDisabled     bool // 0x0040 Hide running light sprites when the ship is disabled.
	UnfoldWhenFiring       bool // 0x0080 Ship unfolds when firing weapons, and folds back up when not firing.
	AdjustForSkew          bool // 0x0100 Adjust ship's visual presentation to correct for the skew caused by graphics that are
}

type Shan struct {
	ID ShanID

	BaseImageID  RleDID // The resource ID of the basic sprite images for this ship.
	BaseMaskID   RleDID // The ID of the corresponding sprite masks (ignored if the base image is an rleD/rle8 resource).
	BaseSetCount int16  // The number of sprite sets for the basic sprite images. A sprite set is usually 36 sprite images, and the graphics for all of a ship's basic sprite sets are stored in the same PICT/rleD/rle8 resource, referred to in BaseImageID.
	BaseXSize    int16  // The X size of each basic sprite image.
	BaseYSize    int16  // The Y size of each basic sprite image.

	BaseTransp int16 // The inherent transparency of the basic sprite images, from 0 (no transparency) to 32 (fully transparent).

	AltImageID  RleDID // The resource ID of the alternating sprite images for this ship. Sprites from the alt sprite sets can be displayed on top of the basic sprite for the ship, cycling through each available sprite set at a rate defined in the Delay field, below. Set to zero if unused.
	AltMaskID   RleDID // The corresponding mask ID. Set to zero if unused.
	AltSetCount int16  // The number of sprite sets for the alternating sprites.
	AltXSize    int16  // The X size of the alt sprite image.
	AltYSize    int16  // The Y size of the alt sprite image.

	GlowImageID RleDID // Engine glow.
	GlowMaskID  RleDID
	GlowXSize   int16
	GlowYSize   int16

	LightImageID RleDID // Running lights.
	LightMaskID  RleDID
	LightXSize   int16
	LightYSize   int16

	WeapImageID RleDID // Weapon effects.
	WeapMaskID  RleDID
	WeapXSize   int16
	WeapYSize   int16

	ShieldImageID RleDID // Shield bubble (shield sprite have a number of frames exactly equal to 1, FramesPer, or BaseSetCount*FramesPer).
	ShieldMaskID  RleDID
	ShieldXSize   int16
	ShieldYSize   int16

	// Note that the first four flags in this field are mutually exclusive - i.e. you can have a ship that banks, unfolds, changes appearance when it is carrying a certain other ship type, or animates in sequence, but these effects can't be combined. The only exception is that having both flags 0x0001 and 0x0002 set is treated specially - it results in a ship whose extra frames are used for banking and which always displays its engine glow when it is turning, whether or not it is actually accelerating. (This something that got thrown in at some point when I realized that it would be necessary to have in order to replicate the behaviour of a certain type of ship from a certain TV show).
	Flags ShanFlags

	AnimDelay     uint16        // The delay between frames of the sprite animations, in 30ths of a second.
	AnimDelayTime time.Duration // The delay between frames of the sprite animations, in 30ths of a second.
	WeapDecay     int16         // The rate at which the weapon glow sprite fades out to transparency, if applicable. 50 is a good median number - lower numbers yield slower decays.
	FramesPer     int16         // The number of frames for one rotation of this ship - usually 36 is a good number, but larger ships can benefit from having more frames per rotation to make their turning animation look smoother. Be sure this value is equal to the actual number of frames per revolution in your images, or bad things will happen!

	// 0 or -1 Ignored.
	// 1       Square-wave blinking:
	//             BlinkValA is the light on-time.
	//             BlinkValB is the delay between blinks.
	//             BlinkValC is the number of blinks in a group.
	//             BlinkValD is the delay between groups.
	// 2       Triangle-wave pulsing:
	//             BlinkValA is the minimum intensity (1-32).
	//             BlinkValB is the intensity increase per frame, x100.
	//             BlinkValC is the maximum intensity (1-32).
	//             BlinkValD is the intensity decrease per frame, x100.
	// 3       Random pulsing:
	//             BlinkValA is the minimum intensity (1-32).
	//             BlinkValB is the maximum intensity (1-32).
	//             BlinkValC is the delay between intensity changes.
	//             BlinkValD is ignored.
	BlinkMode BlinkMode
	BlinkValA int16
	BlinkValB int16
	BlinkValC int16
	BlinkValD int16

	// Here you can set the exit points on the ship sprite for four.
	// different classes of weapons. Note that The "Gun" "Beam" etc.
	// designations are for convenience only, since which set of
	// weapon exit points is used by a given weapon are defined in
	// that weapon's ExitType field. The x & y positions of each
	// weapon exit point are measured in pixels from the centre of the
	// ship when the ship is pointing straight up (frame index 0). See
	// the next four fields if you need to account for any perspective corrections in your sprites.
	GunPosX    [4]int16
	GunPosY    [4]int16
	TurretPosX [4]int16
	TurretPosY [4]int16
	GuidedPosX [4]int16
	GuidedPosY [4]int16
	BeamPosX   [4]int16
	BeamPosY   [4]int16

	// If you have ship sprites that are rendered at an angle, these
	// fields are used to correct for the ships perspective when
	// calculating the weapon exit points (above). If the ship is
	// pointing generally "up" (heading is 0-90 or 270-359) then UpCompressX/Y are
	// used; if the ship is pointing generally "down" (heading is 91-269 degrees)
	// then DnCompressX/Y are used. These values are divided by 100 and then
	// multiplied by the rotated x & y values in the weapon exit point fields to
	// apply a rough correction factor, so values less that 100 will bring the exit
	// points in closer to the ship and values greater than 100 will move the exit
	// points farther out. Experimentation is the best way to learn how this works.
	// Values of zero are interpreted the same as a value of 100, so you can leave
	// this field set to zero if unused.
	UpCompressX int16
	UpCompressY int16
	DnCompressX int16
	DnCompressY int16

	// Here you can set further weapon exit point offsets in order to
	// compensate for skew caused by the z position of a ship
	// graphic's weapon exit point. These values are added to a
	// shot or beam's position after the weapon exit point x & y offsets and the x & y
	// compression factors have been applied, so the effect of these values is not scaled.
	// Positive values here move up the screen, negative values move down the screen.
	// (this is a lot easier to use with the editor than it is to describe).
	GunPosZ    [4]int16
	TurretPosZ [4]int16
	GuidedPosZ [4]int16
	BeamPosZ   [4]int16
}

func ShanFromResource(resource resourcefork.Resource) *Shan {
	return ShanFromBytes(ShanID(resource.ID), resource.Data)
}

func ShanFromBytes(id ShanID, b []byte) *Shan {
	flags := binary.BigEndian.Uint16(b[46:])

	t := &Shan{
		ID:           id,
		BaseImageID:  RleDID(binary.BigEndian.Uint16(b[0:])),
		BaseMaskID:   RleDID(binary.BigEndian.Uint16(b[2:])),
		BaseSetCount: int16(binary.BigEndian.Uint16(b[4:])),
		BaseXSize:    int16(binary.BigEndian.Uint16(b[6:])),
		BaseYSize:    int16(binary.BigEndian.Uint16(b[8:])),
		BaseTransp:   int16(binary.BigEndian.Uint16(b[10:])),
		AltImageID:   RleDID(binary.BigEndian.Uint16(b[12:])),
		AltMaskID:    RleDID(binary.BigEndian.Uint16(b[14:])),
		AltSetCount:  int16(binary.BigEndian.Uint16(b[16:])),
		AltXSize:     int16(binary.BigEndian.Uint16(b[18:])),
		AltYSize:     int16(binary.BigEndian.Uint16(b[20:])),
		GlowImageID:  RleDID(binary.BigEndian.Uint16(b[22:])),
		GlowMaskID:   RleDID(binary.BigEndian.Uint16(b[24:])),
		GlowXSize:    int16(binary.BigEndian.Uint16(b[26:])),
		GlowYSize:    int16(binary.BigEndian.Uint16(b[28:])),
		LightImageID: RleDID(binary.BigEndian.Uint16(b[30:])),
		LightMaskID:  RleDID(binary.BigEndian.Uint16(b[32:])),
		LightXSize:   int16(binary.BigEndian.Uint16(b[34:])),
		LightYSize:   int16(binary.BigEndian.Uint16(b[36:])),
		WeapImageID:  RleDID(binary.BigEndian.Uint16(b[38:])),
		WeapMaskID:   RleDID(binary.BigEndian.Uint16(b[40:])),
		WeapXSize:    int16(binary.BigEndian.Uint16(b[42:])),
		WeapYSize:    int16(binary.BigEndian.Uint16(b[44:])),
		Flags: ShanFlags{
			Banking:                flags&0x0001 == 0x0001,
			AnimatedParts:          flags&0x0002 == 0x0002,
			NoKeyCarried:           flags&0x0004 == 0x0004,
			Sequence:               flags&0x0008 == 0x0008,
			StopAnimationsDisabled: flags&0x0010 == 0x0010,
			HideAltDisabled:        flags&0x0020 == 0x0020,
			HideLightsDisabled:     flags&0x0040 == 0x0040,
			UnfoldWhenFiring:       flags&0x0080 == 0x0080,
			AdjustForSkew:          flags&0x0100 == 0x0100,
		},
		AnimDelay:     binary.BigEndian.Uint16(b[48:]),
		AnimDelayTime: time.Duration(binary.BigEndian.Uint16(b[48:])) * time.Second / 30,
		WeapDecay:     int16(binary.BigEndian.Uint16(b[50:])),
		FramesPer:     int16(binary.BigEndian.Uint16(b[52:])),
		BlinkMode:     BlinkMode(binary.BigEndian.Uint16(b[54:])),
		BlinkValA:     int16(binary.BigEndian.Uint16(b[56:])),
		BlinkValB:     int16(binary.BigEndian.Uint16(b[58:])),
		BlinkValC:     int16(binary.BigEndian.Uint16(b[60:])),
		BlinkValD:     int16(binary.BigEndian.Uint16(b[62:])),
		ShieldImageID: RleDID(binary.BigEndian.Uint16(b[64:])),
		ShieldMaskID:  RleDID(binary.BigEndian.Uint16(b[66:])),
		ShieldXSize:   int16(binary.BigEndian.Uint16(b[68:])),
		ShieldYSize:   int16(binary.BigEndian.Uint16(b[70:])),
		GunPosX: [4]int16{
			int16(binary.BigEndian.Uint16(b[72:])),
			int16(binary.BigEndian.Uint16(b[74:])),
			int16(binary.BigEndian.Uint16(b[76:])),
			int16(binary.BigEndian.Uint16(b[78:])),
		},
		GunPosY: [4]int16{
			int16(binary.BigEndian.Uint16(b[80:])),
			int16(binary.BigEndian.Uint16(b[82:])),
			int16(binary.BigEndian.Uint16(b[84:])),
			int16(binary.BigEndian.Uint16(b[86:])),
		},
		TurretPosX: [4]int16{
			int16(binary.BigEndian.Uint16(b[88:])),
			int16(binary.BigEndian.Uint16(b[90:])),
			int16(binary.BigEndian.Uint16(b[92:])),
			int16(binary.BigEndian.Uint16(b[94:])),
		},
		TurretPosY: [4]int16{
			int16(binary.BigEndian.Uint16(b[96:])),
			int16(binary.BigEndian.Uint16(b[98:])),
			int16(binary.BigEndian.Uint16(b[100:])),
			int16(binary.BigEndian.Uint16(b[102:])),
		},
		GuidedPosX: [4]int16{
			int16(binary.BigEndian.Uint16(b[104:])),
			int16(binary.BigEndian.Uint16(b[106:])),
			int16(binary.BigEndian.Uint16(b[108:])),
			int16(binary.BigEndian.Uint16(b[110:])),
		},
		GuidedPosY: [4]int16{
			int16(binary.BigEndian.Uint16(b[112:])),
			int16(binary.BigEndian.Uint16(b[114:])),
			int16(binary.BigEndian.Uint16(b[116:])),
			int16(binary.BigEndian.Uint16(b[118:])),
		},
		BeamPosX: [4]int16{
			int16(binary.BigEndian.Uint16(b[120:])),
			int16(binary.BigEndian.Uint16(b[122:])),
			int16(binary.BigEndian.Uint16(b[124:])),
			int16(binary.BigEndian.Uint16(b[126:])),
		},
		BeamPosY: [4]int16{
			int16(binary.BigEndian.Uint16(b[128:])),
			int16(binary.BigEndian.Uint16(b[130:])),
			int16(binary.BigEndian.Uint16(b[132:])),
			int16(binary.BigEndian.Uint16(b[134:])),
		},
		UpCompressX: int16(binary.BigEndian.Uint16(b[136:])),
		UpCompressY: int16(binary.BigEndian.Uint16(b[138:])),
		DnCompressX: int16(binary.BigEndian.Uint16(b[140:])),
		DnCompressY: int16(binary.BigEndian.Uint16(b[142:])),
		GunPosZ: [4]int16{
			int16(binary.BigEndian.Uint16(b[144:])),
			int16(binary.BigEndian.Uint16(b[146:])),
			int16(binary.BigEndian.Uint16(b[148:])),
			int16(binary.BigEndian.Uint16(b[150:])),
		},
		TurretPosZ: [4]int16{
			int16(binary.BigEndian.Uint16(b[152:])),
			int16(binary.BigEndian.Uint16(b[154:])),
			int16(binary.BigEndian.Uint16(b[156:])),
			int16(binary.BigEndian.Uint16(b[158:])),
		},
		GuidedPosZ: [4]int16{
			int16(binary.BigEndian.Uint16(b[160:])),
			int16(binary.BigEndian.Uint16(b[162:])),
			int16(binary.BigEndian.Uint16(b[164:])),
			int16(binary.BigEndian.Uint16(b[166:])),
		},
		BeamPosZ: [4]int16{
			int16(binary.BigEndian.Uint16(b[168:])),
			int16(binary.BigEndian.Uint16(b[170:])),
			int16(binary.BigEndian.Uint16(b[172:])),
			int16(binary.BigEndian.Uint16(b[174:])),
		},
	}

	return t
}
