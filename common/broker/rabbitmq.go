package broker

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math"
	"math/big"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	MaxRetryCount = 3
	DQL           = "dlq_main"
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

func HandleRetry(ch *amqp.Channel, d *amqp.Delivery) error {
	if d.Headers == nil {
		d.Headers = amqp.Table{}
	}

	retryCount, ok := d.Headers["x-retry-count"].(int64)
	if !ok {
		retryCount = 0
	}
	retryCount++

	d.Headers["x-retry-count"] = retryCount

	log.Printf("Retrying message %s, retry count: %d", d.Body, retryCount)

	if retryCount >= MaxRetryCount {
		// DQL

		log.Printf("Moving message to DLQ %s", DQL)

		return ch.PublishWithContext(context.Background(), "", DQL, false, false, amqp.Publishing{
			ContentType:  "application/json",
			Headers:      d.Headers,
			Body:         d.Body,
			DeliveryMode: amqp.Persistent,
		})
	}
	time.Sleep(CalcExponentialBackoffTime(time.Second, retryCount))

	return nil
}

func CalcExponentialBackoffTime(base time.Duration, retryCount int64) time.Duration {
	// generate a random number between 0.8 to 1.2 for the jitter
	randNum, _ := rand.Int(rand.Reader, big.NewInt(400))
	numInBetween := 0.8 + (float64(randNum.Int64()) / 1000.0)

	// generate time factor
	delta := math.Pow(2, float64(retryCount)) * numInBetween

	return base * time.Duration(float64(base)*delta)
}
