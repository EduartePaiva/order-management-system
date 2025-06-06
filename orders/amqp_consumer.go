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
	service OrdersService
}

func NewConsumer(service OrdersService) *consumer {
	return &consumer{service}
}

func (c *consumer) Listen(ch *amqp.Channel) {
	q, err := ch.QueueDeclare("", true, false, true, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.QueueBind(q.Name, "", broker.OrderPaidEvent, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	for d := range msgs {
		log.Printf("Received message: %s", d.Body)

		// Extract the headers
		ctx := broker.ExtractAMQPHeader(context.Background(), d.Headers)

		tr := otel.Tracer("amqp")
		_, messageSpan := tr.Start(ctx, fmt.Sprintf("AMQP - consume - %s", q.Name))

		o := &pb.Order{}
		if err := json.Unmarshal(d.Body, o); err != nil {
			d.Nack(false, false)
			log.Printf("failed to unmarshal order: %v", err)
			continue
		}

		// Extract the headers

		_, err := c.service.UpdateOrder(context.Background(), o)
		if err != nil {
			log.Printf("failed to update order: %v", err)

			if err := broker.HandleRetry(ch, &d); err != nil {
				log.Printf("Error handling retry: %v", err)
			}

			continue
		}

		messageSpan.AddEvent("order.updated")
		messageSpan.End()

		log.Println("Order has been updated from AMQP")
		d.Ack(false)
	}

}
