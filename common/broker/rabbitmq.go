package broker

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Connect(user, pass, host, port string) (*amqp.Channel, func() error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s", user, pass, host, port))
	if err != nil {
		log.Fatal(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	err = ch.ExchangeDeclare(OrderCreatedEvent, amqp.ExchangeDirect, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = ch.ExchangeDeclare(OrderPaidEvent, amqp.ExchangeFanout, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	return ch, conn.Close
}
