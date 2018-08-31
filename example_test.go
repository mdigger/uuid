package uuid_test

import (
	"encoding/json"
	"log"

	"github.com/globalsign/mgo/bson"
	"github.com/mdigger/uuid"
)

func Example() {
	uuidData := uuid.New()
	println("UUID:   ", uuidData.String())
	data, err := json.Marshal(uuidData)
	if err != nil {
		log.Fatal(err)
	}
	println("JSON:  ", string(data))
	var newUUID uuid.UUID
	if err := json.Unmarshal(data, &newUUID); err != nil {
		log.Fatal(err)
	}
	println("RESTORE:", newUUID.String())
	data, err = bson.Marshal(uuidData)
	if err != nil {
		log.Fatal(err)
	}
	if err := bson.Unmarshal(data, &newUUID); err != nil {
		log.Fatal(err)
	}
	println("RESTORE:", newUUID.String())
}
