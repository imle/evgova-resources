package resources

import (
	"github.com/imle/resourcefork"
)

// TODO

type SndID IDType

type Snd struct {
	ID SndID
}

func SndFromResource(resource resourcefork.Resource) *Snd {
	panic("not yet implemented")
}
