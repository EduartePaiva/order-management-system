package consul

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	sdk "github.com/hashicorp/consul/api"
)

type registry struct {
	client *sdk.Client
}

func NewRegistry(addr, serviceName string) (*registry, error) {
	config := sdk.DefaultConfig()
	config.Address = addr

	client, err := sdk.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &registry{client: client}, nil
}
func (r *registry) Register(ctx context.Context, instanceID, serviceName, hostPort string) error {
	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 {
		return errors.New("invalid host:port format. Eg: localhost:3000")
	}
	host := parts[0]
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}

	return r.client.Agent().ServiceRegister(&sdk.AgentServiceRegistration{
		ID:      instanceID,
		Address: host,
		Port:    port,
		Name:    serviceName,
		Check: &sdk.AgentServiceCheck{
			CheckID:                        instanceID,
			TLSSkipVerify:                  true,
			TTL:                            "5s",
			Timeout:                        "1s",
			DeregisterCriticalServiceAfter: "10s",
		},
	})

}
func (r *registry) Deregister(ctx context.Context, instanceID, serviceName string) error {
	log.Printf("Deregistering service %s", instanceID)

	return r.client.Agent().CheckDeregister(instanceID)
}
func (r *registry) HealthCheck(instanceID, serviceName string) error {
	return r.client.Agent().UpdateTTL(instanceID, "online", sdk.HealthPassing)
}
func (r *registry) Discover(ctx context.Context, serviceName string) ([]string, error) {
	entries, _, err := r.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	}

	instances := make([]string, 0, len(entries))
	for _, entry := range entries {
		instances = append(instances, fmt.Sprintf("%s:%d", entry.Service.Address, entry.Service.Port))
	}
	return instances, nil
}
