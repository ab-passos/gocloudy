package main

import (
	"errors"
)

type VMDetails struct {
	MachineType string
	Os string
	DevbenchType string
	Baseline string
}

type VMInstance struct {
	DevbenchName string
	VirtualMachineDetails VMDetails 
	StartupScript string
}

type VirtualMachine struct {
	Provider VirtualMachineProvider
	Database VirtualMachineDatabase
}

type VirtualMachineProvider interface {
	CreateVirtualMachine(vmInstance VMInstance) error
	DestroyVirtualMachine(vmName string) error
}

type VirtualMachineDatabase interface {
	MachineExists(VirtualMachineDetails VMDetails) bool
}

func NewVirtualMachine(provider VirtualMachineProvider, database VirtualMachineDatabase) *VirtualMachine {
	return &VirtualMachine {
		Provider: provider,
		Database: database,
	}
}

func (v *VirtualMachine) CreateVirtualMachine(vmInstance VMInstance) error {
	exits := v.Database.MachineExists(vmInstance.VirtualMachineDetails)
	if exits {
		err := v.Provider.CreateVirtualMachine(vmInstance)
		return err
	}
	return errors.New("Requested VM type not available")
}

func (v *VirtualMachine) DestroyVirtualMachine(vmName string) error {
	return v.Provider.DestroyVirtualMachine(vmName)
}