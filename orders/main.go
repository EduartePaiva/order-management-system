package main

import (
	"context"
	"log"
	"net"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/eduartepaiva/order-management-system/common"
	"github.com/eduartepaiva/order-management-system/common/broker"
	"github.com/eduartepaiva/order-management-system/common/discovery"
	"github.com/eduartepaiva/order-management-system/common/discovery/consul"
	"google.golang.org/grpc"
)

var (
	grpcAddr      = common.EnvString("GRPC_ADDR", "localhost:2000")
	consulAddr    = common.EnvString("CONSUL_ADDR", "localhost:8500")
	ampqUser      = common.EnvString("RABBITMQ_USER", "guest")
	ampqPass      = common.EnvString("RABBITMQ_PASS", "guest")
	ampqHost      = common.EnvString("RABBITMQ_HOST", "localhost")
	ampqPort      = common.EnvString("RABBITMQ_PORT", "5672")
	stripePriceID = common.EnvString("STRIPE_PRICE_ID", "some price id")
	jeagerAddr    = common.EnvString("JEAGER_ADDR", "localhost:4318")
)

const (
	serviceName = "orders"
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
	if err := registry.Register(ctx, instanceID, serviceName, grpcAddr); err != nil {
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

	amqpCh, close := broker.Connect(ampqUser, ampqPass, ampqHost, ampqPort)
	defer func() {
		close()
		amqpCh.Close()
	}()

	grpcServer := grpc.NewServer()

	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer l.Close()

	store := NewStore()
	svc := NewService(store)
	svcWithTelemetry := NewTelemetryMiddleware(svc)

	NewGRPCHandler(grpcServer, svcWithTelemetry, amqpCh)

	amqpConsumer := NewConsumer(svcWithTelemetry)
	go amqpConsumer.Listen(amqpCh)

	log.Println("Grpc server started at", grpcAddr)

	if err := grpcServer.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}
