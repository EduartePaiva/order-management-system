package main

import (
	"context"
	"fmt"

	pb "github.com/eduartepaiva/order-management-system/common/api"
)

type Store struct {
	stock map[string]*pb.Item
}

func NewStore() *Store {
	return &Store{
		stock: map[string]*pb.Item{
			"2": {
				ID:       "2",
				Name:     "Potato Chips",
				PriceID:  "price_something",
				Quantity: 10,
			},
			"1": {
				ID:       "1",
				Name:     "Cheese Burger",
				PriceID:  "price_something2",
				Quantity: 20,
			},
		},
	}
}

func (s *Store) GetItem(ctx context.Context, id string) (*pb.Item, error) {
	item, ok := s.stock[id]

	if !ok {
		return nil, fmt.Errorf("Item not found")
	}

	return item, nil
}

func (s *Store) GetItems(ctx context.Context, ids []string) ([]*pb.Item, error) {
	res := make([]*pb.Item, 0)

	for _, id := range ids {
		if i, ok := s.stock[id]; ok {
			res = append(res, i)
		}
	}

	return res, nil
}
