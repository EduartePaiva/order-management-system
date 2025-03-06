package main

import (
	"context"

	pb "github.com/eduartepaiva/order-management-system/common/api"
)

type service struct{}

func NewService() *service {
	return &service{}
}

func (s *service) CreatePayment(ctx context.Context, order *pb.Order) (string, error) {
	// connect to payment processor
	return "", nil
}
