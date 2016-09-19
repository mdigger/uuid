package uuid

import (
	"encoding/json"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

func TestUUID(t *testing.T) {
	uuid := New()
	println("UUID:   ", uuid.String())
	data, err := json.Marshal(uuid)
	if err != nil {
		t.Fatal(err)
	}
	println("JSON:  ", string(data))
	var newUUID UUID
	if err := json.Unmarshal(data, &newUUID); err != nil {
		t.Fatal(err)
	}
	println("RESTORE:", newUUID.String())
	println("----------------------------")
	data, err = bson.Marshal(uuid)
	if err != nil {
		t.Fatal(err)
	}
	if err := bson.Unmarshal(data, &newUUID); err != nil {
		t.Fatal(err)
	}
	println("RESTORE:", newUUID.String())
}
