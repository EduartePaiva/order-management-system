package main

import (
	"context"

	"github.com/eduartepaiva/order-management-system/common"
	pb "github.com/eduartepaiva/order-management-system/common/api"
)

type service struct {
	store OrdersStore
}

func NewService(store OrdersStore) *service {
	return &service{store}
}

func (s *service) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.Order, error) {
	items, err := s.ValidadeOrder(ctx, p)
	if err != nil {
		return nil, err
	}

	return &pb.Order{
		Items:      items,
		CustomerID: p.CustomerID,
		ID:         "42",
		Status:     "pending",
	}, nil
}

func (s *service) ValidadeOrder(ctx context.Context, p *pb.CreateOrderRequest) ([]*pb.Item, error) {
	if len(p.Items) == 0 {
		return nil, common.ErrNoItems
	}
	mergedItems := mergeItemsQuantities(p.Items)

	// TEMPORARY
	itemsWithPrice := make([]*pb.Item, 0)
	for _, item := range mergedItems {
		itemsWithPrice = append(itemsWithPrice, &pb.Item{
			ID:       item.ID,
			Quantity: item.Quantity,
			PriceID:  "stripe_price_id",
		})
	}

	// validate with the stock service
	return itemsWithPrice, nil
}

func mergeItemsQuantities(items []*pb.ItemsWithQuantity) []*pb.ItemsWithQuantity {
	mergedItems := make(map[string]*pb.ItemsWithQuantity)
	for _, i := range items {
		if _, ok := mergedItems[i.ID]; ok {
			mergedItems[i.ID].Quantity += i.Quantity
		} else {
			mergedItems[i.ID] = i
		}
	}
	newItems := make([]*pb.ItemsWithQuantity, 0, len(mergedItems))
	for _, item := range mergedItems {
		newItems = append(newItems, item)
	}
	return newItems

}
