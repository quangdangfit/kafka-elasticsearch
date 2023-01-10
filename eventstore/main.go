package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/gofrs/uuid"

	"github.com/quangdangfit/kafka-elasticsearch/eventstore/eventstroredb"
)

func main() {
	db, err := eventstroredb.NewEventStoreDB(eventstroredb.EventStoreConfig{
		ConnectionString: "esdb://localhost:2113?tls=false",
	})

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close() // nolint: errcheck

	expectedRevision := esdb.NoStream{}
	eventId, _ := uuid.NewV4()
	data, _ := json.Marshal("data")
	eventsData := []esdb.EventData{
		{
			EventID:   eventId,
			EventType: "order",
			Data:      data,
		},
	}
	appendStream, err := db.AppendToStream(
		context.Background(),
		"streamId",
		esdb.AppendToStreamOptions{ExpectedRevision: expectedRevision},
		eventsData...,
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("(Save) stream: {%+v}", appendStream)
}
