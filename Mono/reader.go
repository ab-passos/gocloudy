package main

import (
//	"flag"
//	"os"
	"sync"
)

type MainTestExecutor interface {
	Execute(filePtr *string)
}

func ExecuteTestCase(testExecutor MainTestExecutor, filePtr *string) {
	testExecutor.Execute(filePtr)
}

func createFakeVMInstance() VMInstance {
	return FakeVMInstance{"_"}
}

func createFakeExecutor() Executer {
	return ExecuterTPF{"TPF.py"}
}

func execute(filePtr *string) {
	var testConfiguration = ParseConfigurationXML(*filePtr)

	requiredVMS := GetRequiredVMs(testConfiguration)

	namedInstances := CreateRequiredInstances(
		createFakeVMInstance(),
		requiredVMS)

	copyOfNamedInstances := make([]string, len(namedInstances))
	copy(copyOfNamedInstances, namedInstances)

	ProcessConfigurations(
		testConfiguration,
		namedInstances,
		createFakeExecutor())

	DeleteCreatedInstances(createFakeVMInstance(), copyOfNamedInstances)
}

func main() {

	var wg sync.WaitGroup
	wg.Add(1)
	bucketListener := CreateBucketListener(&wg)
	bucketListener.ListenToBucket()
	wg.Wait()
}
