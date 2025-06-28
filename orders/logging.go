package main

import (
	"context"
	"time"

	pb "github.com/eduartepaiva/order-management-system/common/api"
	"go.uber.org/zap"
)

type LoggingMiddleware struct {
	next OrdersService
}

// CreateOrder implements OrdersService.
func (t *LoggingMiddleware) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.Order, error) {
	start := time.Now()

	defer func() {
		zap.L().Info("CreateOrder", zap.Duration("took", time.Since(start)))
	}()

	return t.next.CreateOrder(ctx, p)
}

// GetOrder implements OrdersService.
func (t *LoggingMiddleware) GetOrder(ctx context.Context, p *pb.GetOrderRequest) (*pb.Order, error) {
	start := time.Now()

	defer func() {
		zap.L().Info("GetOrder", zap.Duration("took", time.Since(start)))
	}()

	return t.next.GetOrder(ctx, p)
}

// UpdateOrder implements OrdersService.
func (t *LoggingMiddleware) UpdateOrder(ctx context.Context, o *pb.Order) (*pb.Order, error) {
	start := time.Now()

	defer func() {
		zap.L().Info("UpdateOrder", zap.Duration("took", time.Since(start)))
	}()

	return t.next.UpdateOrder(ctx, o)
}

// ValidateOrder implements OrdersService.
func (t *LoggingMiddleware) ValidateOrder(ctx context.Context, p *pb.CreateOrderRequest) ([]*pb.Item, error) {
	start := time.Now()

	defer func() {
		zap.L().Info("ValidateOrder", zap.Duration("took", time.Since(start)))
	}()

	return t.next.ValidateOrder(ctx, p)
}

func NewLoggingMiddleware(next OrdersService) OrdersService {
	return &LoggingMiddleware{next}
}
