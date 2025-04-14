package main

import (
	"context"
	"crypto/rand"
	"errors"
	"sync"

	pb "github.com/eduartepaiva/order-management-system/common/api"
)

type store struct {
	// add here the mongoDB
	sync.RWMutex
}

var orders = make([]*pb.Order, 0)

func NewStore() *store {
	return &store{}
}

func (s *store) Create(ctx context.Context, o *pb.CreateOrderRequest, items []*pb.Item) (string, error) {
	s.Lock()
	defer s.Unlock()

	id := rand.Text()
	orders = append(orders, &pb.Order{
		ID:         id,
		CustomerID: o.CustomerID,
		Status:     "pending",
		Items:      items,
	})
	return id, nil
}

func (s *store) Get(ctx context.Context, id, customerID string) (*pb.Order, error) {
	s.RLock()
	defer s.RUnlock()

	for _, order := range orders {
		if order.ID == id && order.CustomerID == customerID {
			return order, nil
		}
	}

	return &pb.Order{}, errors.New("order not found")
}

func (s *store) Update(ctx context.Context, id string, order *pb.Order) error {
	s.Lock()
	defer s.Unlock()

	for i := range orders {
		if orders[i].ID == id {
			orders[i].Status = order.Status
			orders[i].PaymentLink = order.PaymentLink
			return nil
		}
	}
	return nil
}
