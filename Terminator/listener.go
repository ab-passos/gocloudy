package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"cloud.google.com/go/pubsub"
	micro "github.com/micro/go-micro/v2"
)

type BucketListener struct {
	pubSub PubSub
	wg     *sync.WaitGroup
}

type PubSub interface {
	ListenToBucketAtGCP() error
}

type PubSubGCP struct {
}

func CreatePubSubGCP() *PubSubGCP {
	return &PubSubGCP{}
}

func CreateBucketListener(waitGroup *sync.WaitGroup) *BucketListener {
	return &BucketListener{
		pubSub: CreatePubSubGCP(),
		wg:     waitGroup,
	}
}

func subscribe() (*pubsub.Subscription, context.Context) {
	log.Printf("Listenning")
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, "eleanor-270008")
	if err != nil {
		panic(err)
	}
	return client.Subscription("messages-vm"), ctx
}

func (pubSubGCP *PubSubGCP) ListenToBucketAtGCP() error {

	sub, ctx := subscribe()

	log.Printf("Subscribed")
	err2 := sub.Receive(ctx,
		func(ctx context.Context, m *pubsub.Message) {
			log.Printf("Got message: %s", m.Data)
			deleteInstance(string(m.Data))
			m.Ack()
		})
	return err2
}

func (bucketListener *BucketListener) ListenToBucket() {
	err := bucketListener.pubSub.ListenToBucketAtGCP()
	if err != nil {
		bucketListener.wg.Done()
		panic(err)
	}
	bucketListener.wg.Done()
}

func deleteInstance(vmName string) {
	// Create a new service
	service := micro.NewService(micro.Name("terminator.client"))
	// Initialise the client and parse command line flags
	service.Init()

	fmt.Println("Going to call Machines")

	// Create new greeter client
	greeter := NewDevbenchService("Machines", service.Client())

	// Call the greeter
	_, err := greeter.Delete(context.TODO(), &Name{Name: vmName})
	if err != nil {
		fmt.Println("Failed to delete")
		fmt.Println(err)
	}

}
