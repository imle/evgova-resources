package resources

import (
	"encoding/binary"

	"github.com/imle/resourcefork"
)

// The boom resource is used to customize the various explosion types. Nova supports up to 64 different explosion
// types. The graphics for the explosions are loaded from spïn resources 400-463, the sounds are loaded from snd
// resources 300-363, and the behaviour of each explosion type is contained within bööm resources 128-191.

type BoomID IDType

// -1        - No explosion.
// 0-63      - This type of explosion.
// 1000-1063 - Explosion type 0-63, plus a random number of type-0 explosions around it.
type ExplodeType IDType

func (p ExplodeType) Get() (exp BoomID, randType0 bool) {
	return BoomID(p % 1000), p >= 1000
}

const BoomSupported int16 = 64
const BoomSndOffset int16 = 300
const BoomSpinOffset int16 = 400

type Boom struct {
	ID BoomID
	// The rate at which the explosion will animate - a value of 100 will cause each frame of the explosion
	// to appear for exactly one frame of the game animation, and lower values will stretch out the explosion
	// animation and make it stay onscreen longer.
	FrameAdvance int16
	// The index (0-63) of the explosion sound to associate with this explosion type, or -1 for a silent explosion.
	// Usually you'd set this to either -1 or to the same value as the explosion index itself (e.g. 1 for bööm
	// resource 129, etc.) but if you want to use the same sound for two different explosion types, you can do
	// that with this field.
	SoundID SndID
	// The index (0-63) of the explosion graphic to associate with this explosion type. Usually you'd set this
	// to the same value as the explosion index itself (e.g. 1 for bööm resource 129, etc.) but if you want to
	// use the same graphic for two different explosion types (like the small weapon explosion and the
	// ship-breaking-up explosion) you can do that with this field.
	GraphicID SpinID
}

func BoomFromResource(resource resourcefork.Resource) *Boom {
	return BoomFromBytes(BoomID(resource.ID), resource.Data)
}

func BoomFromBytes(id BoomID, b []byte) *Boom {
	t := &Boom{
		ID:           id,
		FrameAdvance: int16(binary.BigEndian.Uint16(b[0:])),
		SoundID:      SndID(int16(binary.BigEndian.Uint16(b[2:])) + BoomSndOffset),
		GraphicID:    SpinID(int16(binary.BigEndian.Uint16(b[4:])) + BoomSpinOffset),
	}

	return t
}
