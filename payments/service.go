package main

import (
	"context"

	pb "github.com/eduartepaiva/order-management-system/common/api"
	"github.com/eduartepaiva/order-management-system/payments/gateway"
	"github.com/eduartepaiva/order-management-system/payments/processor"
)

type service struct {
	processor processor.PaymentProcessor
	gateray   gateway.OrdersGateway
}

func NewService(processor processor.PaymentProcessor, gateway gateway.OrdersGateway) *service {
	return &service{processor: processor, gateray: gateway}
}

func (s *service) CreatePayment(ctx context.Context, order *pb.Order) (string, error) {
	link, err := s.processor.CreatePaymentLink(order)
	if err != nil {
		return "", err
	}

	err = s.gateray.UpdateOrderAfterPaymentLink(ctx, order.ID, link)
	if err != nil {
		return "", err
	}
	return link, nil
}
