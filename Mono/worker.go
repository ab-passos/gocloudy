package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
	"sync"
)

func findMachine(machineIds []string, machine string) int {
	for i, machineId := range machineIds {
		if strings.Contains(machineId, machine) {
			return i
		}
	}
	return -1
}

func remove(slice []string, i int) []string {
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}

func findAvailableMachine(machineIds []string, machine string) string {
	position := findMachine(machineIds, machine)
	var usedMachineId string
	if position != -1 {
		usedMachineId = machineIds[position]
		remove(machineIds, position)
	}

	return usedMachineId
}

type Executer interface {
	Execute(instructions []byte) error
}

type ExecuterTPF struct {
	commandToExecute string
}

func (executerPPF ExecuterTPF) Execute(instructions []byte) error {
	os.Stdout.Write(instructions)
	return nil
}

func ExecuteOnVm(exec Executer, instructions []byte) {
	err := exec.Execute(instructions)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
}

func processTaskWithXML(testConfiguration TestConfiguration, machineId string, wg *sync.WaitGroup, exec Executer) {
	fmt.Println("-> Executing on machine ", machineId)
	fmt.Println("----------------")

	output, err := xml.MarshalIndent(testConfiguration, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	ExecuteOnVm(exec, output)

	fmt.Println("----------------")
	fmt.Println("<- Exiting for machine ", machineId)
	wg.Done()
}

func ProcessConfigurations(
	testConfiguration TestConfiguration,
	namedInstances []string,
	exec Executer) {

	fmt.Println("=> Started execution")
	var wg sync.WaitGroup
	for i := 0; i < len(testConfiguration.Configurations); i++ {

		wg.Add(1)
		testConfigurationPart := testConfiguration

		testConfigurationPart.Configurations =
			[]Configuration{testConfigurationPart.Configurations[i]}

		machineToBeUsed := findAvailableMachine(
			namedInstances,
			testConfiguration.Configurations[i].Machine_type)

		go processTaskWithXML(testConfigurationPart, machineToBeUsed, &wg, exec)
	}
	wg.Wait()
	fmt.Println("<= Ended execution")

}
