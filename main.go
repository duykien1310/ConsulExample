package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/consul/api"
)

func main() {
	// Tạo một kết nối đến Consul
	config := api.DefaultConfig()
	config.Address = "localhost:8500"
	client, err := api.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	// Đăng ký dịch vụ
	serviceID := "my-service-1"
	serviceName := "my-service"
	servicePort := 8080
	reg := &api.AgentServiceRegistration{
		ID:   serviceID,
		Name: serviceName,
		Port: servicePort,
	}
	err = client.Agent().ServiceRegister(reg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Registered service '%s' on port %d\n", serviceName, servicePort)

	// Tìm kiếm dịch vụ
	lookup := &api.AgentServiceLookupOptions{}
	services, _, err := client.Agent().ServicesWithOptions(*lookup)
	if err != nil {
		log.Fatal(err)
	}
	for _, service := range services {
		if service.Service == serviceName {
			fmt.Printf("Found service '%s' at %s:%d\n", serviceName, service.Address, service.Port)
		}
	}
}
