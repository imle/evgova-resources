package resources

import (
	"image"
	"log"

	"github.com/imle/gomacimage"
	"github.com/imle/resourcefork"
)

type PictID IDType

const (
	OutfPictIDOffset = 6000
)

type Pict struct {
	ID    PictID
	Image *image.NRGBA
}

func PictFromResource(resource resourcefork.Resource) *Pict {
	pict, e := PictFromBytes(PictID(resource.ID), resource.Data)
	if e != nil {
		log.Fatal(e)
	}

	return pict
}

func PictFromBytes(id PictID, b []byte) (*Pict, error) {
	img, err := gomacimage.PictFromBytes(b)
	if err != nil {
		return nil, err
	}

	t := &Pict{
		ID:    id,
		Image: img.(*image.NRGBA),
	}

	return t, nil
}
