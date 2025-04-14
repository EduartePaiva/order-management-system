package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	pb "github.com/eduartepaiva/order-management-system/common/api"
	"github.com/eduartepaiva/order-management-system/common/broker"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/webhook"
)

type paymentHTTPHandler struct {
	channel *amqp.Channel
}

func NewHTTPHandler(channel *amqp.Channel) *paymentHTTPHandler {
	return &paymentHTTPHandler{channel}
}

func (h *paymentHTTPHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("/webhook", h.handleCheckoutWebhook)
}

func (h *paymentHTTPHandler) handleCheckoutWebhook(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// Pass the request body and Stripe-Signature header to ConstructEvent, along with the webhook signing key
	// Use the secret provided by Stripe CLI for local testing
	// or your webhook endpoint's secret.
	event, err := webhook.ConstructEvent(body, r.Header.Get("Stripe-Signature"), endpointStripeSecret)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
		w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}

	if event.Type == stripe.EventTypeCheckoutSessionCompleted ||
		event.Type == stripe.EventTypeCheckoutSessionAsyncPaymentSucceeded {
		var cs stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &cs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if cs.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid {
			log.Printf("Payment for checkout session %v succeeded!", cs.ID)

			CustomerID := cs.Metadata["CustomerID"]
			ID := cs.Metadata["ID"]

			ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
			defer cancel()

			o := &pb.Order{
				ID:          ID,
				CustomerID:  CustomerID,
				Status:      "paid",
				PaymentLink: "",
			}

			marshledOrder, _ := json.Marshal(o)

			h.channel.PublishWithContext(ctx, broker.OrderPaidEvent, "", false, false, amqp.Publishing{
				ContentType:  "application/json",
				Body:         marshledOrder,
				DeliveryMode: amqp.Persistent,
			})

			log.Println("Message published order.paid")
		}

	}

	w.WriteHeader(http.StatusOK)
}
