package main

import (
	"context"
	"fmt"

	pb "github.com/eduartepaiva/order-management-system/common/api"
	"go.opentelemetry.io/otel/trace"
)

type TelemetryWithMiddleware struct {
	next PaymentService
}

// CreatePayment implements PaymentService.
func (t TelemetryWithMiddleware) CreatePayment(ctx context.Context, p *pb.Order) (string, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("CreatePayment: %v", p))

	return t.next.CreatePayment(ctx, p)
}

func NewTelemetryMiddleware(svc PaymentService) PaymentService {
	return TelemetryWithMiddleware{svc}
}
