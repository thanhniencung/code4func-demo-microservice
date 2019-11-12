package main

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"os"
)

func main() {
	// Echo instance
	e := echo.New()

	// Routes
	e.GET("/", hello)
	e.GET("/healthcheck", healthcheck)
	e.GET("/order/myorder", MyOrder)

	registerServiceWithConsul()

	// Start server
	e.Logger.Fatal(e.Start(":3001"))
}

// Handler
func hello(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"code":  http.StatusOK,
		"message": "Welcome to Order Service",
	})
}

func healthcheck(c echo.Context) error {
	return c.String(http.StatusOK, "Good!")
}

func MyOrder(c echo.Context) error {
	add, _ := LookupServiceWithConsul("user-service")
	fmt.Println(add)
	return c.JSON(http.StatusOK, echo.Map{
		"orderId": "abzjf4df",
		"total": 500,
	})
}

func registerServiceWithConsul() {
	config := consulapi.DefaultConfig()
	consul, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatalln(err)
	}
	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = "order-service" //replace with service id
	registration.Name = "order-service" //replace with service name
	address := hostname()
	registration.Address = address
	if err != nil {
		log.Fatalln(err)
	}
	registration.Port = 3001
	registration.Check = new(consulapi.AgentServiceCheck)
	registration.Check.HTTP = fmt.Sprintf("http://%s:%v/healthcheck",
		address, 3001)
	registration.Check.Interval = "5s"
	registration.Check.Timeout = "3s"
	consul.Agent().ServiceRegister(registration)
}

func hostname() string {
	hn, err := os.Hostname()
	if err != nil {
		log.Fatalln(err)
	}
	return hn
}

func LookupServiceWithConsul(serviceID string) (string, error) {
	config := consulapi.DefaultConfig()
	client, err := consulapi.NewClient(config)
	if err != nil {
		return "", err
	}
	services, err := client.Agent().Services()
	if err != nil {
		return "", err
	}

	fmt.Sprint(services)

	srvc := services[serviceID]
	address := srvc.Address
	port := srvc.Port
	return fmt.Sprintf("http://%s:%v", address, port), nil
}