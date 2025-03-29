package stripe

import (
	"fmt"

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
	gatewaySuccessURL := fmt.Sprintf("%s/success.html?customerID=%s&orderID=%s", gatewayHTTPAddr, order.CustomerID, order.ID)
	gatewayCancelURL := fmt.Sprintf("%s/cancel.html", gatewayHTTPAddr)

	params := &sdk.CheckoutSessionParams{
		LineItems:  lineItems,
		Mode:       sdk.String(string(sdk.CheckoutSessionModePayment)),
		SuccessURL: sdk.String(gatewaySuccessURL),
		CancelURL:  sdk.String(gatewayCancelURL),
		Metadata: map[string]string{
			"OrderID":    order.ID,
			"CustomerID": order.CustomerID,
		},
	}
	checkoutSession, err := session.New(params)
	if err != nil {
		return "", err
	}
	return checkoutSession.URL, nil
}
