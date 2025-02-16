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

func (s *service) CreateOrder(context.Context) error {
	return nil
}

func (s *service) ValidadeOrder(ctx context.Context, p *pb.CreateOrderRequest) error {
	if len(p.Items) == 0 {
		return common.ErrNoItems
	}
	// TODO: is this really something to be validated? probably not, or if I reassign the p.Items
	mergedItems := mergeItemsQuantities(p.Items)
	// This is some extra because it doesn't make sense to merge, mutate values and not update here
	p.Items = mergedItems

	// validate with the stock service
	return nil
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
