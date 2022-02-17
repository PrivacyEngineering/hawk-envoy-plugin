package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
)

var (
	lines = make(chan message)
)

var (
	routingKey   = get("AMQP_QUEUE", "queue.collector")
	dlq          = get("AMQP_DLQ", "queue.collector.dlq")
	rabbitMqUrl  = get("AMQP_CONNECTION", "amqp://guest:guest@localhost:5672/")
	collectorUrl = get("COLLECTOR_URL", "http://collector.collector-ns/api/usages/batch")
)

func get(key string, def string) string {
	if val, has := os.LookupEnv(key); has {
		return val
	} else {
		log.Printf("No env variable %s found. Using default connection: %s", key, def)
		return def
	}
}

func main() {
	ctx, done := context.WithCancel(context.Background())
	go func() {
		publish(redial(ctx, rabbitMqUrl), lines)
	}()

	cli, err := NewClient()
	failOnError(err, "New rabbitmq client failed")

	select {}
	log.Printf("shutting down")

	if err := cli.Shutdown(); err != nil {
		log.Fatalf("error during shutdown: %s", err)
	}
	done()
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
		fix(data, usage)
	}
	send(data, usages)
}

func toDlq(data []byte) {
	// to dlq
	lines <- data
}
