package main

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/labstack/echo"
	"gopkg.in/resty.v1"
	"log"
	"encoding/json"
	"net/http"
	"os"
)

func main() {
	// Echo instance
	e := echo.New()

	// Routes
	e.GET("/", hello)
	e.GET("/healthcheck", healthcheck)
	e.GET("/order/list/:userId", orderList)

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

func orderList(c echo.Context) error {
	
	type Item struct {
		OrderId string `json:"orderId"`
    	Price   int `json:"price"`
	}

	type User struct {
		FullName string `json:"fullName"`
	}

	type Order struct {
		User User `json:"user"`
		Items []Item `json:"items"`
	}

	item1 := Item {
        OrderId: "123",
        Price:  100000,
    }

  item2 := Item {
        OrderId: "456",
        Price:  200000,
    }

  order := Order{}
  
  order.Items = append(order.Items, item1)
  order.Items = append(order.Items, item2)

  // call user service - get user info from user service
  add, _ := LookupServiceWithConsul("user-service")
	
  client := resty.New()
	res, _ := client.R().
		Get(fmt.Sprintf("%s%s", add, "/user/info"))

	json.Unmarshal([]byte(res.String()), &order.User)

	fmt.Println(order.User)

	return c.JSON(http.StatusOK, order)
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
	services, err := client.Agent().Services() // array, slice
	if err != nil {
		return "", err
	}

	fmt.Sprint(services)

	srvc := services[serviceID]
	address := srvc.Address
	port := srvc.Port
	return fmt.Sprintf("http://%s:%v", address, port), nil
}
