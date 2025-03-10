package inmem

import (
	pb "github.com/eduartepaiva/order-management-system/common/api"
)

type inmem struct{}

func NewInmem() *inmem {
	return &inmem{}
}

func (i *inmem) CreatePaymentLink(*pb.Order) (string, error) {
	return "dummy-link", nil
}
