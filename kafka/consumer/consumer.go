package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	conf := kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	}

	conf["group.id"] = "kafka-go-getting-started"
	conf["auto.offset.reset"] = "earliest"

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

	msgToWebSocket := make(map[int]string)
	count := 0

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

			msgToWebSocket[count] = msg
			count++

			// send to ws
			sendToWebSocket(msg)
			// end of send to ws

			fmt.Println("👍 from kafka producer: " + msg)
		}
	}

	// sendMessageToWebsocket(msgToWebSocket)

	c.Close()
}

func sendToWebSocket(msg string) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/kafka", nil)
	if err != nil {
		log.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("msg", msg)

	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}
	defer resp.Body.Close()
}

// func sendMessageToWebsocket(msg map[int]string) {
// 	for i := 0; i < 10; i++ {
// 		// at kafka service:
// 		// - do get requests with msg[i] as the param OR query param
// 		//
// 		// at gin (ws) service:
// 		// - parse param OR query param,
// 		// - send messages

// 		client := &http.Client{}
// 		req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/fromKafka", nil)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		q := req.URL.Query()
// 		q.Add("msg", msg[i])

// 		req.URL.RawQuery = q.Encode()

// 		resp, err := client.Do(req)
// 		if err != nil {
// 			fmt.Println("Errored when sending request to the server")
// 			return
// 		}
// 		defer resp.Body.Close()
// 	}
// }
