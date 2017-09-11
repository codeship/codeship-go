package codeship

import (
	"os"
	"testing"
)

func TestAuthenticate(t *testing.T) {
	testSetup()
	username := os.Getenv("CODESHIP_USERNAME")
	password := os.Getenv("CODESHIP_PASSWORD")
	client, err := New(username, password, "")
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
	_, err := New("invalid", "invalid", "")
	if err == nil {
		t.Error("Authentication did not throw error when username and password were invalid")
		t.FailNow()
	}
}
