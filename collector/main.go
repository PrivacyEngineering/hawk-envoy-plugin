package main

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	routingKey  = "queue.collector"
	rabbitMqUrl = "amqp://guest:guest@rabbitmq-service:5672/"
)

var (
	lines = make(chan message)
)

func main() {
	ctx, done := context.WithCancel(context.Background())
	go func() {
		publish(redial(ctx, rabbitMqUrl), lines)
	}()

	http.HandleFunc("/", rootHandler)
	log.Println(http.ListenAndServe(":8080", nil))
	done()
}

func rootHandler(_ http.ResponseWriter, request *http.Request) {
	var bodyBytes []byte

	if request.Body == nil {
		return
	}

	bodyBytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Printf("Unable to read body. %v", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing body. %v", err)
		}
	}(request.Body)

	if len(bodyBytes) == 0 {
		return
	}

	log.Printf("BODY: %v", string(bodyBytes))
	lines <- bodyBytes
}
