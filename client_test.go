package vmWareGo

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestFindIP(t *testing.T) {

	var defTransport *http.Transport

	if tr, ok := http.DefaultTransport.(*http.Transport); ok {
		defTransport = &http.Transport{
			Proxy:                 tr.Proxy,
			DialContext:           tr.DialContext,
			MaxIdleConns:          tr.MaxIdleConns,
			IdleConnTimeout:       tr.IdleConnTimeout,
			TLSHandshakeTimeout:   tr.TLSHandshakeTimeout,
			ExpectContinueTimeout: tr.ExpectContinueTimeout,
		}
	}

	// custom transport
	if defTransport != nil {
		defTransport.TLSHandshakeTimeout = time.Minute * 5
		defTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	vmC, err := NewClient(ClientParams{
		URL:      "https://HOST",
		Insecure: true,
		User:     "USER",
		Password: "PASS",
		Ctx:      context.Background(),
	}, defTransport)
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
