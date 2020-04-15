package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/api/compute/v1"
)

const (
	imageURL        = "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-7-wheezy-v20140606"
	someDescription = "compute sample instance"
)

type InstanceData struct {
	patchName    string
	projectId    string
	instanceName string
	zone         string
}

type GCPVirtualMachineProvider struct {
}

func generateStartupScript(patchName string, instanceName string) string {

	script := "#! /bin/bash\n\n" +
		"cd ~\n" +
		"gsutil cp gs://patch-store-bucket/" + patchName + " " + patchName + "\n" +
		"tar -xzf " + patchName + "\n" +
		"python example-testing.py &> outputfile.txt\n" +
		"gsutil cp outputfile.txt gs://vm-tooling/outputfile-" + instanceName + ".txt\n" +
		"INSTANCENAME=" + instanceName + "\n" +
		"gsutil cp gs://vm-tooling/rex-watch-dog rex-watch-dog\n" +
		"chmod 777 rex-watch-dog\n" +
		"echo $INSTANCENAME\n" +
		"./rex-watch-dog -project eleanor-270008 -topic vm-notification -message $INSTANCENAME\n"

	return script
}

func generateUniqueInstanceName(instanceName string) string {
	return instanceName + time.Now().Format("20060102150405")
}

func generatedUniqueDiskName() string {
	return "my-root-pd" + time.Now().Format("20060102150405")
}

func createInstance(service *compute.Service, instanceData InstanceData) {
	prefix := "https://www.googleapis.com/compute/v1/projects/" + instanceData.projectId
	generateInstanceName := generateUniqueInstanceName(instanceData.instanceName)
	pubsubwrite := "https://www.googleapis.com/auth/pubsub"

	script := generateStartupScript(instanceData.patchName, generateInstanceName)
	instance := &compute.Instance{
		Name:        generateInstanceName,
		Description: someDescription,
		MachineType: prefix + "/zones/" + instanceData.zone + "/machineTypes/n1-standard-1",
		Disks: []*compute.AttachedDisk{
			{
				AutoDelete: true,
				Boot:       true,
				Type:       "PERSISTENT",
				InitializeParams: &compute.AttachedDiskInitializeParams{
					DiskName:    generatedUniqueDiskName(),
					SourceImage: imageURL,
				},
			},
		},
		NetworkInterfaces: []*compute.NetworkInterface{
			{
				AccessConfigs: []*compute.AccessConfig{
					{
						Type: "ONE_TO_ONE_NAT",
						Name: "External NAT",
					},
				},
				Network: prefix + "/global/networks/default",
			},
		},
		ServiceAccounts: []*compute.ServiceAccount{
			{
				Email: "eleanor-sa@eleanor-270008.iam.gserviceaccount.com",
				Scopes: []string{
					compute.DevstorageFullControlScope,
					compute.ComputeScope,
					pubsubwrite,
				},
			},
		},
		Metadata: &compute.Metadata{
			Items: []*compute.MetadataItems{
				{
					Key:   "startup-script",
					Value: &script,
				},
			},
		},
	}

	_, err := service.Instances.Insert(instanceData.projectId, instanceData.zone, instance).Do()
	if err != nil {
		log.Fatal(err)
	}
}

func makeInstance(patchName string) {
	ctx := context.Background()
	computeService, err := compute.NewService(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var instanceData = InstanceData{patchName, "eleanor-270008", "devbench", "us-west1-c"}

	createInstance(computeService, instanceData)

	resultingInstances, err := computeService.Instances.List(instanceData.projectId, instanceData.zone).Do()

	for _, instance := range resultingInstances.Items {
		fmt.Printf("cloud: gcp, name: %v, instance id: %v, instance selfLink: %v, machine type: %v, launch time: %v\n",
			instance.Name,
			instance.Id,
			instance.SelfLink,
			instance.MachineType,
			instance.CreationTimestamp)
	}
}

func destroyInstance(instanceName string) {
	ctx := context.Background()
	computeService, err := compute.NewService(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//Instances *InstancesService
	_, err2 := computeService.Instances.Delete("eleanor-270008", "us-west1-c", instanceName).Do()
	if err2 != nil {
		log.Fatal(err)
	}

}

func (gc *GCPVirtualMachineProvider) CreateVirtualMachine(vmInstance VMInstance) error {
	makeInstance(vmInstance.DevbenchName)
	return nil
}

func (gc *GCPVirtualMachineProvider) DestroyVirtualMachine(vmName string) error {
	destroyInstance(vmName)
	return nil
}
