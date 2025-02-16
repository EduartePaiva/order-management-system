package main

import (
	"context"

	pb "github.com/eduartepaiva/order-management-system/common/api"
)

type OrdersService interface {
	CreateOrder(context.Context) error
	ValidadeOrder(context.Context, *pb.CreateOrderRequest) error
}

type OrdersStore interface {
	Create(context.Context) error
}
