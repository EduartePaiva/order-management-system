package stripe

import (
	"github.com/eduartepaiva/order-management-system/common"
	pb "github.com/eduartepaiva/order-management-system/common/api"
	sdk "github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
)

var (
	gatewayHTTPAddr = common.EnvString("GATEWAY_HTTP_ADDR", "http://localhost:8080")
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

func (s *stripe) CreatePaymentLink(order *pb.Order) (string, error) {
	lineItems := make([]*sdk.CheckoutSessionLineItemParams, 0, len(order.Items))

	for _, item := range order.Items {
		lineItems = append(lineItems, &sdk.CheckoutSessionLineItemParams{
			// Provide the exact Price ID (for example, pr_1234) of the product you want to sell
			Price:    sdk.String(item.PriceID),
			Quantity: sdk.Int64(int64(item.Quantity)),
		})
	}
	params := &sdk.CheckoutSessionParams{
		LineItems:  lineItems,
		Mode:       sdk.String(string(sdk.CheckoutSessionModePayment)),
		SuccessURL: sdk.String(gatewayHTTPAddr + "/success"),
		CancelURL:  sdk.String(gatewayHTTPAddr + "/cancel"),
	}
	checkoutSession, err := session.New(params)
	if err != nil {
		return "", err
	}
	return checkoutSession.URL, nil
}
