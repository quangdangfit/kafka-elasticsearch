package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/quangdangfit/kafka-elasticsearch/utils"
)

func main() {
	conf := utils.ReadConfig("./config/config.cfg")

	topic := "test"
	p, err := kafka.NewProducer(&conf)

	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		os.Exit(1)
	}

	// Go-routine to handle message delivery reports and
	// possibly other event types (errors, stats, etc)
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Produced event to topic: %s, partition = %d key = %-10s value = %s\n",
						*ev.TopicPartition.Topic, ev.TopicPartition.Partition, string(ev.Key), string(ev.Value))
				}
			}
		}
	}()

	keys := [...]string{"eabara", "jsmith", "sgarcia", "jbernard", "htanaka", "awalther"}
	values := [...]string{"book", "alarm clock", "t-shirts", "gift card", "batteries"}

	for n := 0; n < 10; n++ {
		key := keys[rand.Intn(len(keys))]
		data := values[rand.Intn(len(values))]
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Key:            []byte(key),
			Value:          []byte(data),
		}, nil)
	}

	// Wait for all messages to be delivered
	p.Flush(15 * 1000)
	p.Close()
}
