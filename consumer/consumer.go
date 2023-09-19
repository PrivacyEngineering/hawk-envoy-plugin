package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

type RabbitClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func NewClient() (*RabbitClient, error) {
	cli := &RabbitClient{
		conn:    nil,
		channel: nil,
		tag:     "hawk-consumer-cli",
		done:    make(chan error),
	}
	var err error
	cli.conn, err = amqp.Dial(rabbitMqUrl)
	if err != nil {
		return nil, fmt.Errorf("dial: %v", err)
	}

	go func() {
		fmt.Printf("closing: %s", <-cli.conn.NotifyClose(make(chan *amqp.Error)))
	}()

	cli.channel, err = cli.conn.Channel()
	if err != nil {
		defer func(conn *amqp.Connection) {
			_ = conn.Close()
		}(cli.conn)
		return nil, fmt.Errorf("channel: %v", err)
	}

	err = cli.channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	_, err = cli.channel.QueueDeclare(routingKey, true, false, false, false, nil)
	if err != nil {
		defer func(channel *amqp.Channel) {
			_ = channel.Close()
		}(cli.channel)
		defer func(conn *amqp.Connection) {
			_ = conn.Close()
		}(cli.conn)
		return nil, fmt.Errorf("queue declare: %v", err)
	}

	_, err = cli.channel.QueueDeclare(dlq, true, false, false, false, nil)
	if err != nil {
		defer func(channel *amqp.Channel) {
			_ = channel.Close()
		}(cli.channel)
		defer func(conn *amqp.Connection) {
			_ = conn.Close()
		}(cli.conn)
		return nil, fmt.Errorf("dlq queue declare: %v", err)
	}

	deliveries, err := cli.channel.Consume(routingKey, cli.tag, false, false, false, false, nil)
	if err != nil {
		defer func(channel *amqp.Channel) {
			_ = channel.Close()
		}(cli.channel)
		defer func(conn *amqp.Connection) {
			_ = conn.Close()
		}(cli.conn)
		return nil, fmt.Errorf("consume: %v", err)
	}

	go handle(deliveries, cli.done)

	return cli, nil
}

func (c *RabbitClient) Shutdown() error {
	// will close() the deliveries channel
	if err := c.channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("consumer cancel failed: %s", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer log.Printf("AMQP shutdown OK")

	// wait for handle() to exit
	return <-c.done
}

func handle(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		log.Printf(
			"got %dB delivery: [%v] %q",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)
		propagate(d.Body)
		err := d.Ack(false)
		if err != nil {
			toDlq(d.Body)
		}
	}
	log.Printf("handle: deliveries channel closed")
	done <- nil
}
