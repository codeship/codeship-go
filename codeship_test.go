package codeship

import (
	"os"
	"testing"
)

func TestAuthenticate(t *testing.T) {
	t.SkipNow()
	testSetup()
	username := os.Getenv("CODESHIP_USERNAME")
	password := os.Getenv("CODESHIP_PASSWORD")
	orgName := os.Getenv("CODESHIP_ORGNAME")
	client, err := New(username, password, orgName)
	if err != nil {
		t.Error("New returned error:", err)
	}
	if client.Authentication.AccessToken == "" {
		t.Error("client.Authentication nil after return", client.Authentication)
		t.FailNow()
	}
}

func TestNewFailsWithoutUsernameOrPassword(t *testing.T) {
	_ = os.Setenv("CODESHIP_USERNAME", "")
	_ = os.Setenv("CODESHIP_PASSWORD", "")
	_, err := New("", "", "")
	if err == nil {
		t.Error("Constructor did not throw error when username and password were not available")
		t.FailNow()
	}
}

func TestAuthenticateFailesWithInvalidUsernameAndPassword(t *testing.T) {
	t.SkipNow()
	_, err := New("invalid", "invalid", "")
	if err == nil {
		t.Error("Authentication did not throw error when username and password were invalid")
		t.FailNow()
	}
}
