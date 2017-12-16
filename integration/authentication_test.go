package integration

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	codeship "github.com/codeship/codeship-go"
)

const (
	organizationName = "codeship"
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

	client, err := codeship.New(user, password)
	if err != nil {
		panic(err)
	}

	org, err = client.Scope(context.Background(), organizationName)
	if err != nil {
		panic(err)
	}
}

func TestAuthenticate(t *testing.T) {
	setup()

	assert.Equal(t, "codeship", org.Name)
	assert.NotEmpty(t, org.UUID)
	assert.NotEmpty(t, org.Scopes)
}
