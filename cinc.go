package resources

import (
	"image"
	"log"

	"github.com/imle/gomacimage"
	"github.com/imle/resourcefork"
)

type CicnID IDType

type Cicn struct {
	ID    CicnID
	Image image.Image
}

func CicnFromResource(resource resourcefork.Resource) *Cicn {
	Cicn, e := CicnFromBytes(CicnID(resource.ID), resource.Data)
	if e != nil {
		log.Fatal(e)
	}

	return Cicn
}

func CicnFromBytes(id CicnID, b []byte) (*Cicn, error) {
	img, err := gomacimage.CicnFromBytes(b)
	if err != nil {
		return nil, err
	}

	t := &Cicn{
		ID:    id,
		Image: img,
	}

	return t, nil
}
