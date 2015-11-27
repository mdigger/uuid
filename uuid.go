// uuid содержит функции для генерации и работы с уникальными идентификаторами
// в формате RFC 4122.
//
// Основное отличие от других аналогичных пакетов:
//  - поддержка только версии UUID V4
//  - полная поддержка сериализации/десериализации в текстовый и бинарный вид, включая
//    JSON, BSON, XML и базы данных
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

// UUID описывает формат уникального идентификатора, соответствующего формату RFC 4122.
type UUID [16]byte

// NewUUID возвращает новый случайный уникальный идентификатор.
func New() (uuid UUID) {
	if _, err := io.ReadFull(rand.Reader, uuid[:]); err != nil {
		panic(err)
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // set version byte
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // set high order byte 0b10{8,9,a,b}
	return
}

// Equal возвращает true, если сравниваемый UUID равен текущему.
func (u UUID) Equal(uuid UUID) bool {
	return bytes.Equal(u[:], uuid[:])
}

// Version возвращает версию алгоритма, использовавшегося для генерации UUID.
func (u UUID) Version() uint {
	return uint(u[6] >> 4)
}

// Bytes возвращает байтовое представление UUID.
func (u UUID) Bytes() []byte {
	return u[:]
}

// String возвращает каноническое строковое представление UUID:
//  xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
func (u UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

// MarshalText обеспечивает поддержику интерфейса encoding.TextMarshaler.
// Результат кодирования в точности соответствует каноническому строковому представлению.
func (u UUID) MarshalText() ([]byte, error) {
	return []byte(u.String()), nil
}

// UnmarshalText обеспечивает поддержику интерфейса encoding.TextUnmarshaler.
// Поддерживаются следующие форматы:
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

// MarshalBinary обеспечивает поддержику интерфейса encoding.BinaryMarshaler.
func (u UUID) MarshalBinary() (data []byte, err error) {
	return u.Bytes(), nil
}

// UnmarshalBinary обеспечивает поддержику интерфейса encoding.BinaryUnmarshaler.
// Возвращает ошибку, если размер данных не равен 16 байтам.
func (u *UUID) UnmarshalBinary(data []byte) error {
	if len(data) != 16 {
		return fmt.Errorf("uuid: UUID must be exactly 16 bytes long, got %d bytes", len(data))
	}
	copy(u[:], data)
	return nil
}

// Value обеспечивает поддержику интерфейса driver.Valuer.
func (u UUID) Value() (driver.Value, error) {
	return u.String(), nil
}

// Scan обеспечивает поддержику интерфейса sql.Scanner.
// Для 16 байтной последовательности используется UnmarshalBinary, а для более длинной
// последовательности или для строки используется UnmarshalText.
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

// Parse разбирает и возвращает UUID из его строкового представления
func Parse(s string) (uuid UUID, err error) {
	err = uuid.UnmarshalText([]byte(s))
	return
}

// GetBSON возвращает представление уникального идентификатора в виде бинарного объекта
// BSON с установленным типом UUID.
func (u UUID) GetBSON() (interface{}, error) {
	return bson.Binary{
		Kind: 0x04,      // тип UUID
		Data: u.Bytes(), // содержимое уникального идентификатора
	}, nil
}

// SetBSON десериализует UUID из внутреннего бинарного представления JSON.
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
