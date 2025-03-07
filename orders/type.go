package main

import (
	"context"

	pb "github.com/eduartepaiva/order-management-system/common/api"
)

type OrdersService interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.Order, error)
	ValidadeOrder(context.Context, *pb.CreateOrderRequest) ([]*pb.Item, error)
}

type OrdersStore interface {
	Create(context.Context) error
}
