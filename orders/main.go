package main

import (
	"context"
	"log"
	"net"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/eduartepaiva/order-management-system/common"
	"github.com/eduartepaiva/order-management-system/common/discovery"
	"github.com/eduartepaiva/order-management-system/common/discovery/consul"
	"google.golang.org/grpc"
)

var (
	grpcAddr   = common.EnvString("GRPC_ADDR", "localhost:2000")
	consulAddr = common.EnvString("CONSUL_ADDR", "localhost:8500")
)

const (
	serviceName = "orders"
)

func main() {
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

	grpcServer := grpc.NewServer()

	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer l.Close()

	store := NewStore()
	svc := NewService(store)
	NewGRPCHandler(grpcServer, svc)

	svc.CreateOrder(context.Background())

	log.Println("Grpc server started at", grpcAddr)

	if err := grpcServer.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}
