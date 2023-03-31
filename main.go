// package main

// import (
// 	"fmt"
// 	"log"

// 	"github.com/hashicorp/consul/api"
// )

// func main() {
// 	// Tạo một kết nối đến Consul
// 	config := api.DefaultConfig()
// 	config.Address = "localhost:8500"
// 	client, err := api.NewClient(config)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Đăng ký dịch vụ
// 	serviceID := "my-service-1"
// 	serviceName := "my-service"
// 	servicePort := 8080
// 	reg := &api.AgentServiceRegistration{
// 		ID:   serviceID,
// 		Name: serviceName,
// 		Port: servicePort,
// 	}
// 	err = client.Agent().ServiceRegister(reg)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("Registered service '%s' on port %d\n", serviceName, servicePort)

// 	// Tìm kiếm dịch vụ
// 	lookup := &api.AgentServiceLookupOptions{}
// 	services, _, err := client.Agent().ServicesWithOptions(*lookup)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	for _, service := range services {
// 		if service.Service == serviceName {
// 			fmt.Printf("Found service '%s' at %s:%d\n", serviceName, service.Address, service.Port)
// 		}
// 	}
// }

package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
)

const (
	ttl     = time.Second * 8
	CheckID = "check_health"
)

type Service struct {
	consulClient *api.Client
}

func NewService() *Service {
	client, err := api.NewClient(&api.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return &Service{
		consulClient: client,
	}
}

func (s *Service) Start() {
	ln, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}

	s.registerService()
	go s.updateHealthCheck()
	s.acceptLoop(ln)
}

func (s *Service) acceptLoop(ln net.Listener) {
	for {
		_, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (s *Service) updateHealthCheck() {
	ticker := time.NewTicker(time.Second * 5)
	for {
		err := s.consulClient.Agent().UpdateTTL(CheckID, "online", api.HealthPassing)
		if err != nil {
			log.Fatal(err)
		}

		<-ticker.C
	}
}

func (s *Service) registerService() {
	check := &api.AgentServiceCheck{
		DeregisterCriticalServiceAfter: ttl.String(),
		TLSSkipVerify:                  true,
		TTL:                            ttl.String(),
		CheckID:                        CheckID,
	}

	register := &api.AgentServiceRegistration{
		ID:      "login_service",
		Name:    "mycluster",
		Tags:    []string{"login"},
		Address: "127.0.0.1",
		Port:    3000,
		Check:   check,
	}

	query := map[string]any{
		"type":        "service",
		"service":     "mycluster",
		"passingonly": true,
	}

	plan, err := watch.Parse(query)
	if err != nil {
		log.Fatal(err)
	}

	plan.HybridHandler = func(index watch.BlockingParamVal, result any) {
		switch msg := result.(type) {
		case []*api.ServiceEntry:
			for _, entry := range msg {
				fmt.Println("new member joined", entry.Service)
			}
		}
		fmt.Println("update cluster", result)
	}
	go func() {
		plan.RunWithConfig("", &api.Config{})
	}()

	if err := s.consulClient.Agent().ServiceRegister(register); err != nil {
		log.Fatal(err)
	}
}

func main() {
	s := NewService()
	s.Start()
}
