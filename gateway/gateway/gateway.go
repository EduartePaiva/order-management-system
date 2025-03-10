package gateway

import (
	"context"

	pb "github.com/eduartepaiva/order-management-system/common/api"
)

type OrdersGateway interface {
	CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.Order, error)
	GetOrder(ctx context.Context, p *pb.GetOrderRequest) (*pb.Order, error)
}
