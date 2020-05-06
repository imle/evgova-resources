package resources

import (
	"image"
	"log"

	"github.com/imle/gomacimage"
	"github.com/imle/resourcefork"
)

type RleDID IDType

const (
	RleDIDOffsetSpob RleDID = 2000
)

type RleD struct {
	ID          RleDID
	Image       image.Image
	Rectangle   image.Rectangle
	CountAcross int
	CountDown   int
}

func RleDFromResource(resource resourcefork.Resource) *RleD {
	RleD, e := RleDFromBytes(RleDID(resource.ID), resource.Data)
	if e != nil {
		log.Fatal(e)
	}

	return RleD
}

func RleDFromBytes(id RleDID, b []byte) (*RleD, error) {
	rle, err := gomacimage.RleFromBytes(b)
	if err != nil {
		return nil, err
	}

	t := &RleD{
		ID:          id,
		Image:       rle.Image,
		Rectangle:   rle.Rectangle,
		CountAcross: rle.CountAcross,
		CountDown:   rle.CountDown,
	}

	return t, nil
}
