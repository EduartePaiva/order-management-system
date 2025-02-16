package main

import (
	"context"
	"log"

	pb "github.com/eduartepaiva/order-management-system/common/api"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer
	service OrdersService
}

func NewGRPCHandler(grpcServer *grpc.Server, service OrdersService) {
	handler := &grpcHandler{service: service}
	pb.RegisterOrderServiceServer(grpcServer, handler)
}

func (h *grpcHandler) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Printf("New order received! Order %v", p)
	err := h.service.ValidadeOrder(ctx, p)
	if err != nil {
		log.Printf("Invalid order")
	}
	log.Println(p.Items)

	return &pb.Order{ID: "42"}, nil
}
