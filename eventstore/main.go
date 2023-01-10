package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

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

	for {
		expectedRevision := esdb.StreamExists{}
		eventId, _ := uuid.NewV4()
		data := map[string]interface{}{
			"id":             "63283f04f84bcc8e1893b09c",
			"created_at":     1663581956789,
			"created_by":     "62efee8d1bbe176aa8feeab9",
			"partnership_id": "630720682f9a32ef3356b3ab",
			"ticket_id":      "6321521b1f9c7f6e8665c437",
			"time":           1663581956789,
		}
		bData, _ := json.Marshal(data)
		eventsData := []esdb.EventData{
			{
				EventID:     eventId,
				EventType:   "order",
				ContentType: esdb.JsonContentType,
				Data:        bData,
			},
		}
		appendStream, err := db.AppendToStream(
			context.Background(),
			"$checkin-logs-stream-1",
			esdb.AppendToStreamOptions{ExpectedRevision: expectedRevision},
			eventsData...,
		)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("(Save) stream: {%+v}", appendStream)
		time.Sleep(time.Second)
	}
}
