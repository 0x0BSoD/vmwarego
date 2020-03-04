package vmWareGo

import (
	"fmt"
	"log"
	"testing"
)

var cParams ClientParams

func TestFindIP(t *testing.T) {

	cParams = ClientParams{
		URL: "URL",
		Insecure: true,
		User: "USER",
		Password: "PASS",
	}


	vms, err := AllVMs(cParams)
	if err != nil {
		log.Println(err)
	}

	for _, i := range vms {
		fmt.Println(i.Config.Annotation)
	}
}
