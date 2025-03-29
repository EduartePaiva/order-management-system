package main

import pb "github.com/eduartepaiva/order-management-system/common/api"

type CreateOrderRequest struct {
	Order         *pb.Order `json:"order"`
	RedirectToURL string    `json:"redirectToURL"`
}
