package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/EventStore/EventStore-Client-Go/esdb"

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

	subsOpt := esdb.SubscribeToStreamOptions{From: esdb.Start{}}
	subscriber, err := db.SubscribeToStream(context.Background(), "$checkin-logs-stream-1", subsOpt)
	if err != nil {
		log.Fatal(err)
	}
	defer subscriber.Close()

	for {
		subscriptionEvent := subscriber.Recv()
		if err != nil {
			log.Fatal(err)
		}

		var data map[string]interface{}
		event := subscriptionEvent.EventAppeared.Event
		json.Unmarshal(event.Data, &data)
		log.Printf("(Get) stream eventId: %s, data: %+v", event.EventID, data)
	}
}
