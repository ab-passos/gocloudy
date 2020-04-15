package main

import "sync"

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	bucketListener := CreateBucketListener(&wg)
	bucketListener.ListenToBucket()
	wg.Wait()
}
