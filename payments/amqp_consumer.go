package main

import (
	"log"

	"github.com/eduartepaiva/order-management-system/common/broker"
	amqp "github.com/rabbitmq/amqp091-go"
)

type consumer struct {
	service PaymentService
}

func NewConsumer(service PaymentService) *consumer {
	return &consumer{service: service}
}

func (c *consumer) Listen(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)
	}
}
