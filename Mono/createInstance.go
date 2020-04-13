package main

import (
	"context"
	"fmt"

	micro "github.com/micro/go-micro/v2"
	
)


func MakeInstances(patchName string) {
	// Create a new service
	service := micro.NewService(micro.Name("mono.client"))
	// Initialise the client and parse command line flags
	service.Init()

	// Create new greeter client
	greeter := NewDevbenchService("Machines", service.Client())

	// Call the greeter
	_, err := greeter.Create(context.TODO(), &Name{Name: patchName})
	if err != nil {
		fmt.Println("Failed")
		fmt.Println(err)
	}

}
