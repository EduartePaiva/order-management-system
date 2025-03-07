package processor

import (
	pb "github.com/eduartepaiva/order-management-system/common/api"
)

type PaymentProcessor interface {
	CreatePaymentLink(*pb.Order) (string, error)
}
