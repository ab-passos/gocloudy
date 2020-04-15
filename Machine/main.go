package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"

	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/micro/go-micro/v2"
)

type MySqlVM struct {
}

func (db *MySqlVM) MachineExists(VirtualMachineDetails VMDetails) bool {
	return true
}

type DevbenchImpl struct{}

func (dh *DevbenchImpl) Create(ctx context.Context, in *Name, out *empty.Empty) error {
	gcp := &GCPVirtualMachineProvider{}
	mysql := &MySqlVM{}
	vm := NewVirtualMachine(gcp, mysql)

	vmInstance := VMInstance{
		DevbenchName: in.String(),
		VirtualMachineDetails: VMDetails{
			MachineType:  "type",
			Os:           "os",
			DevbenchType: "type",
			Baseline:     "baseline"},
		StartupScript: "startup"}

	return vm.CreateVirtualMachine(vmInstance)
}

func (dh *DevbenchImpl) Delete(ctx context.Context, in *Name, out *empty.Empty) error {
	gcp := &GCPVirtualMachineProvider{}
	mysql := &MySqlVM{}
	vm := NewVirtualMachine(gcp, mysql)
	fmt.Println("Got message to delete ", in.String())
	var rgx = regexp.MustCompile(`\"(.*?)\"`)
	rs := rgx.FindStringSubmatch(in.String())

	return vm.DestroyVirtualMachine(rs[1])

}

func main() {

	fmt.Println("PROJECT_ID:", os.Getenv("PROJECT_ID"))

	service := micro.NewService(
		micro.Name("Machines"),
	)

	service.Init()

	RegisterDevbenchHandler(service.Server(), new(DevbenchImpl))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
