package main

import (
	"fmt"
	"strconv"
)

type VMInstance interface {
	CreateInstances(requiredVMs RequiredVMs) []string
	Delete(instances []string)
}

type FakeVMInstance struct {
	instancePrefix string
}

func (fakeVMInstance FakeVMInstance) CreateInstances(requiredVMs RequiredVMs) []string {
	var arr []string

	var counter int = 1
	for k, v := range requiredVMs {
		for i := 0; i < v; i++ {
			t := strconv.Itoa(counter)
			arr = append(arr, k+fakeVMInstance.instancePrefix+t)
			counter++
		}
	}

	return arr
}

func (fakeVMInstance FakeVMInstance) Delete(instances []string) {

	for i := 0; i < len(instances); i++ {
		fmt.Println("Deleting:", instances[i])
	}

}

func CreateRequiredInstances(vmInstance VMInstance, requiredVMs RequiredVMs) []string {
	return vmInstance.CreateInstances(requiredVMs)
}

func DeleteCreatedInstances(vmInstance VMInstance, instances []string) {
	vmInstance.Delete(instances)
}
