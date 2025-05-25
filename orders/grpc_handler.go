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
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer

	service OrdersService
	channel *amqp.Channel
}

func NewGRPCHandler(grpcServer *grpc.Server, service OrdersService, channel *amqp.Channel) {
	handler := &grpcHandler{service: service, channel: channel}
	pb.RegisterOrderServiceServer(grpcServer, handler)
}

func (h *grpcHandler) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.Order, error) {
	q, err := h.channel.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	tr := otel.Tracer("amqp")
	amqpCtx, messageSpam := tr.Start(ctx, fmt.Sprintf("AMQP - publish - %s", q.Name))
	defer messageSpam.End()

	log.Printf("New order received! Order %v", p)
	order, err := h.service.CreateOrder(amqpCtx, p)
	if err != nil {
		return nil, err
	}

	marshaledOrder, err := json.Marshal(order)
	if err != nil {
		log.Println("Failed to marshal json", err)
		return nil, err
	}

	// Inject the headers
	headers := broker.InjectAMQPHeaders(amqpCtx)

	h.channel.PublishWithContext(amqpCtx, "", q.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         marshaledOrder,
		DeliveryMode: amqp.Persistent,
		Headers:      headers,
	})

	return order, nil
}

func (h *grpcHandler) GetOrder(ctx context.Context, p *pb.GetOrderRequest) (*pb.Order, error) {
	log.Println("getting the order: ", p)
	return h.service.GetOrder(ctx, p)
}

func (h *grpcHandler) UpdateOrder(ctx context.Context, order *pb.Order) (*pb.Order, error) {
	return h.service.UpdateOrder(ctx, order)
}
