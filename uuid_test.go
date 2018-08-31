package uuid

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"testing"

	"github.com/globalsign/mgo/bson"
)

func TestUUID(t *testing.T) {
	uuid := New()
	println("UUID:   ", uuid.String())
	data, err := json.Marshal(uuid)
	if err != nil {
		t.Error(err)
	}
	println("JSON:  ", string(data))
	var newUUID UUID
	if err := json.Unmarshal(data, &newUUID); err != nil {
		t.Error(err)
	}
	println("RESTORE:", newUUID.String())
	if !uuid.Equal(newUUID) {
		t.Error("bad restore")
	}
	if newUUID.Version() != 4 {
		t.Error("bad version", newUUID.Version())
	}

	data, err = bson.Marshal(uuid)
	if err != nil {
		t.Fatal(err)
	}
	if err := bson.Unmarshal(data, &newUUID); err != nil {
		t.Fatal(err)
	}
	println("RESTORE:", newUUID.String())
	if !uuid.Equal(newUUID) {
		t.Error("bad restore")
	}

	var buf bytes.Buffer
	err = gob.NewEncoder(&buf).Encode(uuid)
	if err != nil {
		t.Error(err)
	}
	err = gob.NewDecoder(&buf).Decode(&newUUID)
	if err != nil {
		t.Error(err)
	}
	println("RESTORE:", newUUID.String())

	data, err = bson.Marshal(bson.Binary{
		Kind: 0x05,
		Data: uuid.Bytes(),
	})
	if err != nil {
		t.Error(err)
	}
	err = bson.Unmarshal(data, &newUUID)
	if err == nil {
		t.Error("bad SetBSON")
	}
	data, err = bson.Marshal(bson.Binary{
		Kind: 0x04,
		Data: uuid.Bytes()[1:],
	})
	if err != nil {
		t.Error(err)
	}
	err = bson.Unmarshal(data, &newUUID)
	if err == nil {
		t.Error("bad SetBSON")
	}
}

func TestUUIDUnmarshal(t *testing.T) {
	for _, uuidStr := range []string{
		"6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		"{6ba7b810-9dad-11d1-80b4-00c04fd430c8}",
		"urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		"6ba7b8109dad11d180b400c04fd430c8",
	} {
		_, err := Parse(uuidStr)
		if err != nil {
			t.Error(err)
		}
	}

	var uuid UUID
	if uuid.UnmarshalText([]byte("12345678")) == nil {
		t.Error("bad unmarshal")
	}
	if uuid.UnmarshalText([]byte("6ba7b8109dad11d180b400c04fd430cw")) == nil {
		t.Error("bad unmarshal")
	}
}
