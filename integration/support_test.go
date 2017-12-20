package integration

import (
	"context"
	"os"

	codeship "github.com/codeship/codeship-go"
)

const (
	organizationName = "codeship"
	organizationUUID = "1c150f00-e93d-0133-b53e-76bef8d7b14f"

	proProjectName = "codeship/codeship-go"
	proProjectUUID = "c38f3280-792b-0135-21bb-4e0cf8ff365b"

	basicProjectName = "codeship/shipscope"
	basicProjectUUID = "688357c0-6652-0135-f3b7-1268853457c2"
)

var (
	org *codeship.Organization
)

func setup() {
	if org != nil {
		return
	}

	user := os.Getenv("CODESHIP_USER")
	if user == "" {
		panic("CODESHIP_USER required")
	}

	password := os.Getenv("CODESHIP_PASSWORD")
	if password == "" {
		panic("CODESHIP_PASSWORD required")
	}

	client, err := codeship.New(codeship.NewBasicAuth(user, password))
	if err != nil {
		panic(err)
	}

	org, err = client.Organization(context.Background(), organizationName)
	if err != nil {
		panic(err)
	}
}
