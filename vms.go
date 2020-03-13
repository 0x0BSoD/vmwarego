package vmWareGo

import (
	"errors"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

func (vm *Vmware) VmRetrieve(targetVm []mo.VirtualMachine, fields []string) error {
	manager := view.NewManager(vm.Client.Client)
	containerView, err := manager.CreateContainerView(vm.Ctx, vm.Client.ServiceContent.RootFolder, []string{"VirtualMachine"}, true)
	if err != nil {
		return err
	}
	defer containerView.Destroy(vm.Ctx)

	err = containerView.Retrieve(vm.Ctx, []string{"VirtualMachine"}, fields, &targetVm)
	if err != nil {
		return err
	}

	return nil
}

func (vm *Vmware) VmsAllSummary() ([]types.VirtualMachineSummary, error) {
	manager := view.NewManager(vm.Client.Client)

	containerView, err := manager.CreateContainerView(vm.Ctx, vm.Client.ServiceContent.RootFolder, []string{"VirtualMachine"}, true)
	if err != nil {
		return nil, err
	}
	defer containerView.Destroy(vm.Ctx)

	var vms []mo.VirtualMachine
	err = containerView.Retrieve(vm.Ctx, []string{"VirtualMachine"}, []string{"summary"}, &vms)
	if err != nil {
		return nil, err
	}

	var result []types.VirtualMachineSummary

	for _, vm := range vms {
		result = append(result, vm.Summary)
	}

	return result, nil
}

func (vm *Vmware) VmInfo(name string) (types.VirtualMachineSummary, error) {

	vms, err := vm.VmsFilter(name)
	if err != nil {
		return types.VirtualMachineSummary{}, err
	}

	if len(vms) > 1 {
		return types.VirtualMachineSummary{}, errors.New("get more than one vm")
	}

	return vms[0].Summary, nil
}

func (vm *Vmware) VmsFilter(nameFilter string) ([]mo.VirtualMachine, error) {
	var result []mo.VirtualMachine

	manager := view.NewManager(vm.Client.Client)

	containerView, err := manager.CreateContainerView(vm.Ctx, vm.Client.ServiceContent.RootFolder, nil, true)
	if err != nil {
		return nil, err
	}
	defer containerView.Destroy(vm.Ctx)

	var filter = make(map[string]types.AnyType)
	filter["name"] = nameFilter

	vmsRefs, err := containerView.Find(vm.Ctx, []string{"VirtualMachine"}, filter)
	if err != nil {
		return nil, err
	}
	if err := containerView.Destroy(vm.Ctx); err != nil {
		panic(err)
	}

	pc := property.DefaultCollector(vm.Client.Client)
	props := []string{"summary"}
	var vms []mo.VirtualMachine

	if len(vmsRefs) != 0 {
		err = pc.Retrieve(vm.Ctx, vmsRefs, props, &vms)
		if err != nil {
			panic(err)
		}
	}

	result = append(result, vms...)

	return result, nil
}
