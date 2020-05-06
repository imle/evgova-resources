package resources

import (
	"encoding/binary"

	"github.com/imle/resourcefork"
)

// TODO: Str patching

// String Arrays are constructed as:
//
//  0x0000 - Total strings in array [N]
//    - repeating N times
//    0x00 - String length [L]
//    0x00 * L - Strings are not null terminated

type StrAID IDType

type StrA struct {
	ID StrAID

	Values []*string
}

func StrAFromResource(resource resourcefork.Resource) *StrA {
	return StrAFromBytes(StrAID(resource.ID), resource.Data)
}

func StrAFromBytes(id StrAID, b []byte) *StrA {
	// First word is string count
	strCount := int16(binary.BigEndian.Uint16(b[0:]))

	t := &StrA{
		ID:     id,
		Values: make([]*string, strCount),
	}

	// Start after first word
	pos := int16(2)

	var strLen uint8
	for i := int16(0); i < strCount; i++ {
		strLen = b[pos]
		pos++

		s := string(b[pos : pos+int16(strLen)])
		t.Values[i] = &s
		pos += int16(strLen)
	}

	return t
}
