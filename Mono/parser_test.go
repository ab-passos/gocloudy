package main

import (
	"testing"
)

func TestParseConfigurationXML(t *testing.T) {
	testConfiguration := ParseConfigurationXML("testSets/only-configs.test_scenario.xml")
	if len(testConfiguration.Configurations) != 4 {
		t.Errorf("Configuration lenght is wrong, got: %d, want: %d.", len(testConfiguration.Configurations), 4)
	}

	array_of_configurations := [4]string{
		"Config1",
		"Config2",
		"Config3",
		"Config4"}

	for i := 0; i < len(array_of_configurations); i++ {
		if testConfiguration.Configurations[i].Configuration_id != array_of_configurations[i] {
			t.Errorf("Configuration Id is wrong, got: %s, want: %s.", testConfiguration.Configurations[i].Configuration_id, array_of_configurations[i])
		}
	}

	devbench := "VM"

	for i := 0; i < len(array_of_configurations); i++ {
		if testConfiguration.Configurations[i].Environment != devbench {
			t.Errorf("Devbench is wrong, got: %s, want: %s.", testConfiguration.Configurations[i].Environment, devbench)
		}
	}

	machine_type := "Debian"

	for i := 0; i < len(array_of_configurations); i++ {
		if testConfiguration.Configurations[i].Machine_type != machine_type {
			t.Errorf("Machine_type is wrong, got: %s, want: %s.", testConfiguration.Configurations[i].Machine_type, machine_type)
		}
	}
}

func TestGetRequiredVMs(t *testing.T) {
	testConfiguration := ParseConfigurationXML("testSets/only-configs.test_scenario.xml")

	requiredVMS := GetRequiredVMs(testConfiguration)
	if len(requiredVMS) != 1 {
		t.Errorf("Number of required VM type is wrong, got: %d, want: %d.", len(requiredVMS), 1)
	}
	if requiredVMS["Debian"] != 4 {
		t.Errorf("Number of required VMs is wrong, got: %d, want: %d.", len(requiredVMS), 4)
	}

}
