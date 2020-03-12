package vmWareGo

import (
	"context"
	"log"
	"testing"
)

var cParams ClientParams

func TestFindIP(t *testing.T) {
	cParams = ClientParams{
		URL:      "URL",
		Insecure: true,
		User:     "USER",
		Password: "PASS",
		Ctx:      context.Background(),
	}

	vmC, err := NewClient(cParams)
	if err != nil {
		log.Println(err)
	}

	_ = PrettyPrint(vmC.Client.ServiceContent.About)
}
