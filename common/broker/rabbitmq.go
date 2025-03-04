package broker

import (
	"fmt"
	"log"

	ampq "github.com/rabbitmq/amqp091-go"
)

func Connect(user, pass, host, port string) (*ampq.Channel, func() error) {
	conn, err := ampq.Dial(fmt.Sprintf("ampq://%s:%s@%s:%s", user, pass, host, port))
	if err != nil {
		log.Fatal(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	err = ch.ExchangeDeclare(OrderCreatedEvent, ampq.ExchangeDirect, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = ch.ExchangeDeclare(OrderCreatedPaid, ampq.ExchangeFanout, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	return ch, conn.Close
}
