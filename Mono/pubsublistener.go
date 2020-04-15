package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	pubsub "cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
)

type Parser interface {
	GetBucketName(raw map[string]interface{}) string
	GetPatchName(raw map[string]interface{}) string
}

type JsonParser struct {
}

func (jsonParser *JsonParser) GetBucketName(raw map[string]interface{}) string {
	bucket, _ := json.Marshal(raw["bucket"])
	return string(bucket)
}

func (jsonParser *JsonParser) GetPatchName(raw map[string]interface{}) string {
	patch, _ := json.Marshal(raw["name"])
	return string(patch)
}

func CreateJsonParser() *JsonParser {
	return &JsonParser{}
}

type PubSub interface {
	ListenToBucketAtGCP() error
}

type PubSubGCP struct {
	parser Parser
}

func CreatePubSubGCP() *PubSubGCP {
	return &PubSubGCP{
		parser: CreateJsonParser(),
	}
}

func subscribe() (*pubsub.Subscription, context.Context) {
	log.Printf("Listenning")
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, GetProjectId())
	if err != nil {
		panic(err)
	}
	return client.Subscription(GetSubscriptionForPatchAdded()), ctx
}

func (pubSubGCP *PubSubGCP) ListenToBucketAtGCP() error {

	sub, ctx := subscribe()

	log.Printf("Subscribed")
	err2 := sub.Receive(ctx,
		func(ctx context.Context, m *pubsub.Message) {
			log.Printf("Got message: %s", m.Data)
			var raw map[string]interface{}
			if err := json.Unmarshal(m.Data, &raw); err != nil {
				panic(err)
			}
			bucket := pubSubGCP.parser.GetBucketName(raw)
			log.Printf("Bucket is %s\n", bucket)
			obj := pubSubGCP.parser.GetPatchName(raw)

			MakeInstances(obj)
			m.Ack()

		})
	return err2
}

type BucketListener struct {
	pubSub PubSub
	wg     *sync.WaitGroup
}

func (bucketListener *BucketListener) ListenToBucket() {
	err := bucketListener.pubSub.ListenToBucketAtGCP()
	if err != nil {
		bucketListener.wg.Done()
		panic(err)
	}
	bucketListener.wg.Done()
}

func CreateBucketListener(waitGroup *sync.WaitGroup) *BucketListener {
	return &BucketListener{
		pubSub: CreatePubSubGCP(),
		wg:     waitGroup,
	}
}

//TODO

func DownloadFromBucket(bucket string, object string) error {
	bucket = bucket[1 : len(bucket)-1]
	object = object[1 : len(object)-1]
	println(bucket)
	println(object)
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	data, err := read(client, bucket, object)
	if err != nil {
		log.Fatalf("Cannot read object: %v", err)
	}
	fmt.Printf("Object contents: %s\n", data)

	return err
}

func read(client *storage.Client, bucket, object string) ([]byte, error) {
	// [START download_file]
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	rc, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	return data, nil
	// [END download_file]
}

func DownloadFile(url string, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	println(url)

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		println("GetUrl")
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
