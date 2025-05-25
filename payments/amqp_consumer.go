package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	pb "github.com/eduartepaiva/order-management-system/common/api"
	"github.com/eduartepaiva/order-management-system/common/broker"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
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

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)

		// Extract the headers
		ctx := broker.ExtractAMQPHeader(context.Background(), d.Headers)

		tr := otel.Tracer("amqp")
		_, messageSpan := tr.Start(ctx, fmt.Sprintf("AMQP - consume - %s", q.Name))

		order := &pb.Order{}
		err := json.Unmarshal(d.Body, order)
		if err != nil {
			d.Nack(false, false)
			log.Printf("failed to Unmarshal order: %v", err)
			continue
		}

		paymentLink, err := c.service.CreatePayment(context.Background(), order)
		if err != nil {
			log.Printf("failed to create checkout link: %v", err)

			if err = broker.HandleRetry(ch, &d); err != nil {
				log.Printf("Error handling retry: %v", err)
			}
			d.Nack(false, false)
			continue
		}

		messageSpan.AddEvent(fmt.Sprintf("payments.created: %s", paymentLink))
		messageSpan.End()

		log.Printf("payment link: %s", paymentLink)
		d.Ack(false)
	}
}
