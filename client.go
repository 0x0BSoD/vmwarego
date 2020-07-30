package vmWareGo

import (
	"context"
	"net/http"
	"net/url"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/soap"
)

type ClientParams struct {
	URL       string
	User      string
	Password  string
	VCSToPass []string
	Insecure  bool
	Ctx       context.Context
}

type Vmware struct {
	Client *govmomi.Client
	Ctx    context.Context
}

// NewClient creates a govmomi.Client
func NewClient(clientParams ClientParams, customTransport *http.Transport) (Vmware, error) {

	// Parse URL from string
	_, err := url.ParseRequestURI(clientParams.URL)
	if err != nil {
		return Vmware{}, err
	}

	u, err := soap.ParseURL(clientParams.URL)
	if err != nil {
		return Vmware{}, err
	}

	u.User = url.UserPassword(clientParams.User, clientParams.Password)

	var vm Vmware

	// Connect and log in to ESX or vCenter
	c, err := govmomi.NewClient(clientParams.Ctx, u, clientParams.Insecure)
	if err != nil {
		return Vmware{}, err
	}

	if customTransport != nil {
		c.Client.Transport = customTransport
	}

	vm.Client = c
	vm.Ctx = clientParams.Ctx

	return vm, nil
}

func (vm *Vmware) Close() {
	vm.Client.CloseIdleConnections()
}
