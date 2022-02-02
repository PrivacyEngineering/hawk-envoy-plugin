package main

import (
	"context"
	"github.com/streadway/amqp"
	"log"
)

type message []byte

type session struct {
	*amqp.Connection
	*amqp.Channel
}

// Close tears the connection down, taking the channel with it.
func (s session) Close() error {
	if s.Connection == nil {
		return nil
	}
	return s.Connection.Close()
}

// redial continually connects to the URL, exiting the program when no longer possible
func redial(ctx context.Context, url string) chan chan session {
	log.Printf("Executing readial")
	sessions := make(chan chan session)

	go func() {
		sess := make(chan session)
		defer close(sessions)

		for {
			select {
			case sessions <- sess:
			case <-ctx.Done():
				log.Println("shutting down session factory")
				return
			}

			conn, err := amqp.Dial(url)
			if err != nil {
				log.Fatalf("cannot (re)dial: %v: %q", err, url)
			}

			ch, err := conn.Channel()
			if err != nil {
				log.Fatalf("cannot create channel: %v", err)
			}
			log.Printf("Creating new session")

			select {
			case sess <- session{conn, ch}:
			case <-ctx.Done():
				log.Println("shutting down new session")
				return
			}
		}
	}()

	return sessions
}

// publish messages to a reconnecting session to a fanout exchange.
// It receives from the application specific source of messages.
func publish(sessions chan chan session, messages <-chan message) {
	for session := range sessions {
		var (
			running bool
			reading = messages
			pending = make(chan message, 1)
			confirm = make(chan amqp.Confirmation, 1)
		)

		pub := <-session

		// publisher confirms for this channel/connection
		if err := pub.Confirm(false); err != nil {
			log.Printf("publisher confirms not supported")
			close(confirm) // confirms not supported, simulate by always nacking
		} else {
			pub.NotifyPublish(confirm)
		}

		log.Printf("publishing...")

	Publish:
		for {
			var body message
			select {
			case confirmed, ok := <-confirm:
				if !ok {
					break Publish
				}
				if !confirmed.Ack {
					log.Printf("nack message %d, body: %q", confirmed.DeliveryTag, string(body))
				}
				reading = messages

			case body = <-pending:
				err := pub.Publish("", dlq, false, false, amqp.Publishing{
					ContentType: "application/json",
					Body:        body,
				})
				// Retry failed delivery on the next session
				if err != nil {
					pending <- body
					_ = pub.Close()
					break Publish
				}

			case body, running = <-reading:
				// all messages consumed
				if !running {
					return
				}
				// work on pending delivery until ack'd
				pending <- body
				reading = nil
			}
		}
	}
}

func consuming(sessions chan chan session, msgs chan amqp.Delivery) {
	for session := range sessions {
		var (
			receiver chan *amqp.Error
		)

		ch := <-session

		receiver = ch.Channel.NotifyClose(make(chan *amqp.Error, 1))
		Consumer:
		for {
			log.Printf("Consumer in action")
			body, err := ch.Consume(routingKey, "", true, false, false, false, nil)
			if err != nil {
				log.Printf("Unable to consume %s", routingKey)
				return
			}

			select {
			case read := <-body:
				log.Println("Read message active")
				msgs <- read

			case <-receiver:
				log.Println("Notify close active")
				break Consumer
			}
		}
	}
}
