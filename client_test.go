package vmWareGo

import (
	"fmt"
	"testing"
)

func TestFindIP(t *testing.T) {
	vmC, err := NewClient(ClientParams{
		URL:      "https://HOST",
		Insecure: true,
		User:     "USER",
		Password: "PASS",
		Ctx:      context.Background(),
	})
	if err != nil {
		t.Fatal(err)
	}

	vms, err := vmC.VmsFilter("base-w2k16-*")
	if err != nil {
		t.Fatal(err)
	}

	type Tags struct {
		Template string `json:"Template"`
	}
	var tags Tags

	err = vmC.VmRetrieve(vms, []string{"summary"})
	if err != nil {
		t.Fatal(err)
	}

	for _, vm := range vms {
		err = vmC.VmCustomFields(vm, &tags)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("%s || Template => %s\n", vm.Summary.Config.Name, tags.Template)
	}

	vmC.Close()
}
