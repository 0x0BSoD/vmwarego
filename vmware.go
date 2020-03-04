package vmWareGo

import (
	"context"
	"fmt"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/vim25/mo"
	"net/url"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
)

type ClientParams struct {
	URL       string
	User      string
	Password  string
	VCSToPass []string
	Insecure  bool
}

type Tags struct {
	Template string `json:"Template"`
	SWE      string `json:"swe,omitempty"`
}

// NewClient creates a govmomi.Client
func NewClient(ctx context.Context, c ClientParams) (*govmomi.Client, error) {

	// Parse URL from string
	u, err := soap.ParseURL(c.URL)
	if err != nil {
		return nil, err
	}

	u.User = url.UserPassword(c.User, c.Password)

	// Connect and log in to ESX or vCenter
	return govmomi.NewClient(ctx, u, c.Insecure)
}

func NameFilter(c *govmomi.Client, nameFilter string) ([]mo.VirtualMachine, error) {
	ctx := context.Background()
	var result []mo.VirtualMachine

	manager := view.NewManager(c.Client)

	containerView, err := manager.CreateContainerView(ctx, c.ServiceContent.RootFolder, nil, true)
	if err != nil {
		return nil, err
	}
	defer containerView.Destroy(ctx)

	var filter = make(map[string]types.AnyType)
	filter["name"] = nameFilter

	vmsRefs, err := containerView.Find(ctx, []string{"VirtualMachine"}, filter)
	if err != nil {
		return nil, err
	}
	if err := containerView.Destroy(ctx); err != nil {
		panic(err)
	}

	pc := property.DefaultCollector(c.Client)
	props := []string{"summary"}
	var vms []mo.VirtualMachine

	if len(vmsRefs) != 0 {
		err = pc.Retrieve(ctx, vmsRefs, props, &vms)
		if err != nil {
			panic(err)
		}
	}

	result = append(result, vms...)

	return result, nil
}

func VMInfo(name string, clientParameters ClientParams) (types.VirtualMachineSummary, error) {
	result := types.VirtualMachineSummary{}

	ctx := context.Background()

	c, err := NewClient(ctx, clientParameters)
	if err != nil {
		return result, err
	}

	vms, err := NameFilter(c, name)
	if err != nil {
		return result, err
	}
	if vms[0].Summary.Guest.ToolsStatus == "toolsNotRunning" {
		return result, fmt.Errorf("tools not running vm name: %s", name)
	}

	fmt.Printf("%T", vms[0].Summary)

	result = vms[0].Summary

	return result, nil
}

func AllVMs(clientParameters ClientParams) ([]types.VirtualMachineSummary, error) {
	ctx := context.Background()

	c, err := NewClient(ctx, clientParameters)
	if err != nil {
		return nil, err
	}

	manager := view.NewManager(c.Client)

	containerView, err := manager.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"VirtualMachine"}, true)
	if err != nil {
		return nil, err
	}
	defer containerView.Destroy(ctx)

	var vms []mo.VirtualMachine
	err = containerView.Retrieve(ctx, []string{"VirtualMachine"}, []string{"summary"}, &vms)
	if err != nil {
		return nil, err
	}

	var result []types.VirtualMachineSummary

	for _, vm := range vms {
		result = append(result, vm.Summary)
	}

	return result, nil
}
