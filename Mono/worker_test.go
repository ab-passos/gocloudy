package main

import (
	"reflect"
	"testing"
)

func TestFindMachine(t *testing.T) {

	var requiredVMs = []string{
		"MachineTypeA",
		"MachineTypeB",
		"MachineTypeC",
	}

	result := findMachine(requiredVMs, "MachineTypeB")

	if result != 1 {
		t.Errorf("findMachineFailed")
	}
}

func TestFindMachineNotExisting(t *testing.T) {

	var requiredVMs = []string{
		"MachineTypeA",
		"MachineTypeB",
		"MachineTypeC",
	}

	result := findMachine(requiredVMs, "MachineTypeD")

	if result != -1 {
		t.Errorf("findMachineFailed")
	}
}

func TestRemoveMachine(t *testing.T) {

	var requiredVMs = []string{
		"MachineTypeA",
		"MachineTypeB",
		"MachineTypeC",
	}

	var expectedResult = []string{
		"MachineTypeA",
		"MachineTypeC",
	}

	result := remove(requiredVMs, 1)

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Remove machine failed")
	}
}

func TestFindAvailableMachine(t *testing.T) {

	var machineIds = []string{
		"MachineTypeA1",
		"MachineTypeB2",
		"MachineTypeC3",
	}

	result := findAvailableMachine(machineIds, "MachineTypeA")

	if result != "MachineTypeA1" {
		t.Errorf("Find Available machines failed, it shoud be MachineTypeA1 and is %s", result)
	}
}
