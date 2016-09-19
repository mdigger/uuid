package uuid

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"testing"

	"gopkg.in/mgo.v2/bson"
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
		t.Fatal(err)
	}
	err = gob.NewDecoder(&buf).Decode(&newUUID)
	if err != nil {
		t.Fatal(err)
	}
	println("RESTORE:", newUUID.String())

}

func TestUUIDUnmarshal(t *testing.T) {
	var uuid UUID
	for _, uuidStr := range []string{
		"6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		"{6ba7b810-9dad-11d1-80b4-00c04fd430c8}",
		"urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		"6ba7b8109dad11d180b400c04fd430c8",
	} {
		err := uuid.UnmarshalText([]byte(uuidStr))
		if err != nil {
			t.Error(err)
		}
	}
	if uuid.UnmarshalText([]byte("12345678")) == nil {
		t.Error("bad unmarshal")
	}
	if uuid.UnmarshalText([]byte("6ba7b8109dad11d180b400c04fd430cw")) == nil {
		t.Error("bad unmarshal")
	}
}
