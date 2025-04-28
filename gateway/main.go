package main

import (
	"context"
	"log"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/eduartepaiva/order-management-system/common"
	"github.com/eduartepaiva/order-management-system/common/discovery"
	"github.com/eduartepaiva/order-management-system/common/discovery/consul"
	"github.com/eduartepaiva/order-management-system/gateway/gateway"
)

var (
	httpAddr   = common.EnvString("HTTP_ADDR", ":8080")
	consulAddr = common.EnvString("CONSUL_ADDR", "localhost:8500")
	jeagerAddr = common.EnvString("JEAGER_ADDR", "localhost:4318")
)

const (
	serviceName = "gateway"
)

func main() {
	if err := common.SetGlobalTracer(context.TODO(), serviceName, jeagerAddr); err != nil {
		log.Fatal("failed to set global tracer")
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)

	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}
	if err := registry.Register(ctx, instanceID, serviceName, httpAddr); err != nil {
		panic(err)
	}
	defer registry.Deregister(ctx, instanceID, serviceName)

	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				log.Fatal("failed to health check")
			}
			time.Sleep(time.Second * 1)
		}
	}()

	mux := http.NewServeMux()

	ordersGateway := gateway.NewGRPCGateway(registry)

	handler := NewHandler(ordersGateway)
	handler.registerRoutes(mux)

	log.Printf("Starting http server at port %s", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("Failed to start http server")
	}
}
