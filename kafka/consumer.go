package kafka

import (
	"example/unit-test-hello-world/localWebsocket"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func ConsumerInit() {
	log.Println("Running kafka consumer...")
	conf := kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "kafka-go-getting-started",
		"auto.offset.reset": "end",
	}

	c, err := kafka.NewConsumer(&conf)
	if err != nil {
		fmt.Printf("Failed to create consumer: %s", err)
		os.Exit(1)
	}

	topic := "purchases"
	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		fmt.Printf("ERROR: Expected %v but got %v", nil, err)
	}

	// Set up a channel for handling Ctrl-C, etc
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Process messages
	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev, err := c.ReadMessage(100 * time.Millisecond)
			if err != nil {
				// Errors are informational and automatically handled by the consumer
				continue
			}

			msg := fmt.Sprintf("topic = %s key = %-10s value = %s\n",
				*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))

			// send to ws
			localWebsocket.Writer("hello from kafka consumer!")
			// end of send to ws

			fmt.Println("from kafka producer: " + msg)
		}
	}

	c.Close()
}
