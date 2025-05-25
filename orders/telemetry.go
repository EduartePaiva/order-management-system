package main

import (
	"context"
	"fmt"

	pb "github.com/eduartepaiva/order-management-system/common/api"
	"go.opentelemetry.io/otel/trace"
)

type TelemetryMiddleware struct {
	next OrdersService
}

// CreateOrder implements OrdersService.
func (t *TelemetryMiddleware) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.Order, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("CreateOrder: %v", p))

	return t.next.CreateOrder(ctx, p)
}

// GetOrder implements OrdersService.
func (t *TelemetryMiddleware) GetOrder(ctx context.Context, p *pb.GetOrderRequest) (*pb.Order, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("GetOrder: %v", p))

	return t.next.GetOrder(ctx, p)
}

// UpdateOrder implements OrdersService.
func (t *TelemetryMiddleware) UpdateOrder(ctx context.Context, o *pb.Order) (*pb.Order, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("UpdateOrder: %v", o))

	return t.next.UpdateOrder(ctx, o)
}

// ValidateOrder implements OrdersService.
func (t *TelemetryMiddleware) ValidateOrder(ctx context.Context, p *pb.CreateOrderRequest) ([]*pb.Item, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("ValidateOrder: %v", p))

	return t.next.ValidateOrder(ctx, p)
}

func NewTelemetryMiddleware(next OrdersService) OrdersService {
	return &TelemetryMiddleware{next}
}
