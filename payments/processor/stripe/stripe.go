package stripe

import (
	"github.com/eduartepaiva/order-management-system/common"
	pb "github.com/eduartepaiva/order-management-system/common/api"
	sdk "github.com/stripe/stripe-go/v81"
)

// initialize the stripe key
func init() {
	sdk.Key = common.EnvString("STRIPE_KEY", "NOT_A_KEY")
}

type stripe struct {
}

func NewProcessor() *stripe {
	return &stripe{}
}

func (s *stripe) CreatePaymentLink(*pb.Order) (string, error) {
	return "", nil
}
