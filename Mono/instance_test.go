package main

import (
	"testing"
)

func TestInstanceCreation(t *testing.T) {

	requiredVMs := RequiredVMs{
		"MachineTypeA": 3,
		"MachineTypeB": 1,
		"MachineTypeC": 5,
	}

	var fake FakeVMInstance

	result := CreateRequiredInstances(fake, requiredVMs)

	if len(result) != 9 {
		t.Errorf("It should be 9 and it is %s", result[0])
	}
}

func TestInstanceDelete(t *testing.T) {

	var requiredVMs = []string{
		"MachineTypeA",
		"MachineTypeB",
		"MachineTypeC",
	}

	var fake FakeVMInstance

	DeleteCreatedInstances(fake, requiredVMs)
}
