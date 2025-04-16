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
	"go.opentelemetry.io/otel"
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

	createDLQAndDLX(ch)
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

	return ch.PublishWithContext(
		context.Background(),
		d.Exchange,
		d.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Headers:      d.Headers,
			Body:         d.Body,
			DeliveryMode: amqp.Persistent,
		},
	)
}

func CalcExponentialBackoffTime(base time.Duration, retryCount int64) time.Duration {
	// generate a random number between 0.8 to 1.2 for the jitter
	randNum, _ := rand.Int(rand.Reader, big.NewInt(400))
	numInBetween := 0.8 + (float64(randNum.Int64()) / 1000.0)

	// generate time factor
	delta := math.Pow(2, float64(retryCount)) * numInBetween

	return base * time.Duration(float64(base)*delta)
}

func createDLQAndDLX(ch *amqp.Channel) error {
	q, err := ch.QueueDeclare(
		"main_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Declare DLX
	dlx := "dlx_main"
	err = ch.ExchangeDeclare(
		dlx,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Bind main queue to DLX
	err = ch.QueueBind(
		q.Name,
		"",
		dlx,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Declare DLQ
	_, err = ch.QueueDeclare(
		DQL,
		true,
		false,
		false,
		false,
		nil,
	)

	return err
}

type AmqpHeaderCarrier map[string]interface{}

func (a AmqpHeaderCarrier) Get(k string) string {
	value, ok := a[k]
	if !ok {
		return ""
	}

	return value.(string)
}

func (a AmqpHeaderCarrier) Set(k string, v string) {
	a[k] = v
}

func (a AmqpHeaderCarrier) Keys() []string {
	keys := make([]string, len(a))
	i := 0

	for k := range a {
		keys[i] = k
		i++
	}

	return keys
}

func ExtractAMQPHeader(ctx context.Context, headers map[string]interface{}) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, AmqpHeaderCarrier(headers))
}
