package main

import (
	"context"
	"testing"

	pb "github.com/eduartepaiva/order-management-system/common/api"
	inmemRegistry "github.com/eduartepaiva/order-management-system/common/discovery/inmem"
	"github.com/eduartepaiva/order-management-system/payments/gateway"
	"github.com/eduartepaiva/order-management-system/payments/processor/inmem"
)

func TestService(t *testing.T) {
	processor := inmem.NewInmem()
	registry := inmemRegistry.NewRegistry()

	gatway := gateway.NewGRPCGateway(registry)

	svc := NewService(processor, gatway)

	t.Run("should create a payment link", func(t *testing.T) {
		link, err := svc.CreatePayment(context.Background(), &pb.Order{})
		if err != nil {
			t.Errorf("CreatePayment() error = %v, want nil", err)
		}

		if link == "" {
			t.Error("CreatePayment() link is empty")
		}
	})
}
