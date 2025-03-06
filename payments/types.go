package main

import (
	"context"

	pb "github.com/eduartepaiva/order-management-system/common/api"
)

type PaymentService interface {
	CreatePayment(context.Context, *pb.Order) (string, error)
}
