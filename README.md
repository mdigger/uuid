# Unique IDs in RFC 4122

[![GoDoc](https://godoc.org/github.com/mdigger/uuid?status.svg)](https://godoc.org/github.com/mdigger/uuid)
[![Build Status](https://travis-ci.org/mdigger/uuid.svg)](https://travis-ci.org/mdigger/uuid)
[![Coverage Status](https://coveralls.io/repos/github/mdigger/uuid/badge.svg?branch=master)](https://coveralls.io/github/mdigger/uuid?branch=master)

Package uuid contains functions for creating and working with unique IDs in
RFC 4122.

The main difference from other similar packages:

1. support only versions of UUID V4;
2. full support for serialization/deserialization to text and binary form,
including JSON, BSON, XML and databases.

```go
package main

import (
	"encoding/json"
	"log"

	"github.com/mdigger/uuid"
	"gopkg.in/mgo.v2/bson"
)

func main() {
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
```
