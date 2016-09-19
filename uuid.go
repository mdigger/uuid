// Package uuid contains functions for creating and working with unique IDs in
// RFC 4122.
//
// The main difference from other similar packages:
//
// 1. support only versions of UUID V4
//
// 2. full support for serialization/deserialization to text and binary form,
// including JSON, BSON, XML and databases.
package uuid

import (
	"bytes"
	"crypto/rand"
	"database/sql/driver"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"gopkg.in/mgo.v2/bson"
)

// UUID describes the format of the unique identifier corresponding to RFC 4122.
type UUID [16]byte

// NewUUID returns a new random unique identifier.
func New() (uuid UUID) {
	if _, err := io.ReadFull(rand.Reader, uuid[:]); err != nil {
		panic(err)
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // set version byte
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // set high order byte 0b10{8,9,a,b}
	return
}

// Equal returns true if the UUID is equal to the current compare.
func (u UUID) Equal(uuid UUID) bool {
	return bytes.Equal(u[:], uuid[:])
}

// Version returns the version of the algorithm used to generate the UUID.
func (u UUID) Version() uint {
	return uint(u[6] >> 4)
}

// Bytes returns a byte representation of the UUID.
func (u UUID) Bytes() []byte {
	return u[:]
}

// String returns the canonical string representation of a UUID:
//  xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
func (u UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

// MarshalText provides the HMDI supports the interface encoding.TextMarshaler.
// The result of the encoding corresponds exactly to the canonical string
// representation.
func (u UUID) MarshalText() ([]byte, error) {
	return []byte(u.String()), nil
}

// UnmarshalText provides support for the interface encoding.TextUnmarshaler.
// The following formats are supported:
//  "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
//  "{6ba7b810-9dad-11d1-80b4-00c04fd430c8}",
//  "urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8"
func (u *UUID) UnmarshalText(text []byte) (err error) {
	if len(text) < 32 {
		return fmt.Errorf("uuid: invalid UUID string: %s", text)
	}
	if bytes.Equal(text[:9], []byte("urn:uuid:")) {
		text = text[9:]
	} else if text[0] == '{' {
		text = text[1:]
	}
	b := u[:]
	for _, byteGroup := range []int{8, 4, 4, 4, 12} {
		if text[0] == '-' {
			text = text[1:]
		}
		if _, err = hex.Decode(b[:byteGroup/2], text[:byteGroup]); err != nil {
			return err
		}
		text = text[byteGroup:]
		b = b[byteGroup/2:]
	}
	return
}

// MarshalBinary provides the HMDI supports the interface
// encoding.BinaryMarshaler.
func (u UUID) MarshalBinary() (data []byte, err error) {
	return u.Bytes(), nil
}

// UnmarshalBinary provides support for the interface encoding.BinaryUnmarshaler.
// Returns an error if data size is not equal to 16 bytes.
func (u *UUID) UnmarshalBinary(data []byte) error {
	if len(data) != 16 {
		return fmt.Errorf("uuid: UUID must be exactly 16 bytes long, got %d bytes", len(data))
	}
	copy(u[:], data)
	return nil
}

// Value provides support for the interface driver.Valuer.
func (u UUID) Value() (driver.Value, error) {
	return u.String(), nil
}

// Scan provides support for the sql interface.Scanner.
// For the 16 byte sequence is used UnmarshalBinary, whereas the longer
// sequence, or string is used UnmarshalText.
func (u *UUID) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		if len(src) == 16 {
			return u.UnmarshalBinary(src)
		}
		return u.UnmarshalText(src)
	case string:
		return u.UnmarshalText([]byte(src))
	default:
		return fmt.Errorf("uuid: cannot convert %T to UUID", src)
	}
}

// Parse parses and returns a UUID from its string representation.
func Parse(s string) (uuid UUID, err error) {
	err = uuid.UnmarshalText([]byte(s))
	return
}

// GetBSON returns a representation of the unique identifier in the form of the
// BSON binary object with the set type UUID.
func (u UUID) GetBSON() (interface{}, error) {
	return bson.Binary{
		Kind: 0x04,      // тип UUID
		Data: u.Bytes(), // содержимое уникального идентификатора
	}, nil
}

// SetBSON deserializes the UUID from the internal binary representation of JSON.
func (uuid *UUID) SetBSON(raw bson.Raw) error {
	var bin = new(bson.Binary)
	if err := raw.Unmarshal(bin); err != nil {
		return err
	}
	if bin.Kind != 0x04 {
		return errors.New("bson: bad UUID binary type")
	}
	return uuid.UnmarshalBinary(bin.Data)
}
