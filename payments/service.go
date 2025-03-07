package main

import (
	"context"

	pb "github.com/eduartepaiva/order-management-system/common/api"
	"github.com/eduartepaiva/order-management-system/payments/processor"
)

type service struct {
	processor processor.PaymentProcessor
}

func NewService(processor processor.PaymentProcessor) *service {
	return &service{processor: processor}
}

func (s *service) CreatePayment(ctx context.Context, order *pb.Order) (string, error) {
	return s.processor.CreatePaymentLink(order)
}
