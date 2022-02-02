package main

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/streadway/amqp"
)

var (
	lines = make(chan message)
	msgs  = make(chan amqp.Delivery, 1)
)

const (
	routingKey  = "queue.collector"
	dlq         = "queue.collector.dlq"
	rabbitMqUrl = "amqp://guest:guest@rabbitmq-service:5672/"
)

func main() {
	ctx, _ := context.WithCancel(context.Background())
	go func() {
		publish(redial(ctx, rabbitMqUrl), lines)
	}()

	go func() {
		consuming(redial(ctx, rabbitMqUrl), msgs)
	}()

	log.Printf("Ready to consume messages")
	for {
		select {
		case d := <- msgs:
			propagate(d.Body)
		}
	}
}

func propagate(data []byte) {
	var usages []map[string]interface{}
	err := json.Unmarshal(data, &usages)
	if err != nil {
		log.Printf("Unable to parse message: %s", err)
		toDlq(data)
		return
	}

	for _, usage := range usages {
		metadata,ok := usage["metadata"].(map[string]interface{})
		if !ok {
			log.Printf("Unable to transform metadata to map: %s", err)
			toDlq(data)
			return
		}
		unixTS,ok := metadata["timestamp"].(string)
		if !ok {
			log.Printf("Unable to transform metadata.timestamp to string. Metadata: %v, Error: %s", metadata, err)
			toDlq(data)
			return
		}

		i, err := strconv.ParseFloat(unixTS, 64)
		if err != nil {
			log.Printf("Unable to transform unix timestamp: %s", err)
			toDlq(data)
			return
		}
		i = i / 1e9

		tm := time.Unix(int64(i), 0).Format(time.RFC3339)

		metadata["timestamp"] = tm
		log.Printf("Timestamp converted: %s", tm)
	}

	parsed,err := json.Marshal(usages)
	if err != nil {
		log.Printf("Unable to marshall json: %s", err)
		toDlq(data)
		return;
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(parsed).
		Post("http://collector.collector-ns/api/usages/batch")

	if err != nil || resp.StatusCode() != 200 {
		log.Printf("Back to dlq for [%d, %s]", resp.StatusCode(), err)
		toDlq(parsed)
	}
}

func toDlq(data []byte) {
	// to dlq
	lines <- data
}
