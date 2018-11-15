// +build integration

package integration

import (
	"context"
	"log"
	"os"
	"testing"

	codeship "github.com/codeship/codeship-go"
)

const (
	organizationName = "codeship"
	organizationUUID = "1c150f00-e93d-0133-b53e-76bef8d7b14f"

	proProjectName = "codeship/codeship-tool-examples"
	proProjectUUID = "5eda0420-40c6-0133-ef9c-0e8a33e740fc"

	basicProjectName = "codeship/merrygoround"
	basicProjectUUID = "20a13930-c925-0134-952e-3a0fd8dae151"
)

var (
	org *codeship.Organization
)

func TestMain(m *testing.M) {
	log.SetFlags(0)

	if org != nil {
		os.Exit(m.Run())
		return
	}

	user := os.Getenv("CODESHIP_USER")
	if user == "" {
		log.Fatal("CODESHIP_USER env var required")
	}

	password := os.Getenv("CODESHIP_PASSWORD")
	if password == "" {
		log.Fatal("CODESHIP_PASSWORD env var required")
	}

	client, err := codeship.New(codeship.NewBasicAuth(user, password))
	if err != nil {
		log.Fatal(err)
	}

	org, err = client.Organization(context.Background(), organizationName)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}
