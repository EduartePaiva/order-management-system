package discovery

import (
	"context"
	"math/rand"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ServiceConnection(ctx context.Context, serviceName string, registry Registry) (*grpc.ClientConn, error) {
	services, err := registry.Discover(ctx, serviceName)
	if err != nil {
		return nil, err
	}
	service := services[rand.Intn(len(services))]
	return grpc.NewClient(service, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
