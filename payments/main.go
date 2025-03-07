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
	"github.com/eduartepaiva/order-management-system/payments/processor/stripe"
	"google.golang.org/grpc"
)

var (
	grpcAddr   = common.EnvString("GRPC_ADDR", "localhost:2001")
	consulAddr = common.EnvString("CONSUL_ADDR", "localhost:8500")
	ampqUser   = common.EnvString("RABBITMQ_USER", "guest")
	ampqPass   = common.EnvString("RABBITMQ_PASS", "guest")
	ampqHost   = common.EnvString("RABBITMQ_HOST", "localhost")
	ampqPort   = common.EnvString("RABBITMQ_PORT", "5672")
)

const (
	serviceName = "payment"
)

func main() {
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	// register consul
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
	// broker connection
	amqpCh, close := broker.Connect(ampqUser, ampqPass, ampqHost, ampqPort)
	defer func() {
		close()
		amqpCh.Close()
	}()

	stripeProcessor := stripe.NewProcessor()
	svc := NewService(stripeProcessor)
	amqpConsumer := NewConsumer(svc)
	go amqpConsumer.Listen(amqpCh)
	// gRPC server
	grpcServer := grpc.NewServer()

	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer l.Close()

	log.Println("Grpc server started at", grpcAddr)
	if err := grpcServer.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}
