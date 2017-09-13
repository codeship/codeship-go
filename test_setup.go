package codeship

import "os"

const CodeshipPassword = ""
const CodeshipUsername = ""
const CodeshipOrgName = ""

func testSetup() {
	username := os.Getenv("CODESHIP_USERNAME")
	if username == "" {
		os.Setenv("CODESHIP_USERNAME", CodeshipUsername)
	}
	password := os.Getenv("CODESHIP_PASSWORD")
	if password == "" {
		os.Setenv("CODESHIP_PASSWORD", CodeshipPassword)
	}
	orgName := os.Getenv("CODESHIP_ORGNAME")
	if orgName == "" {
		os.Setenv("CODESHIP_ORGNAME", CodeshipOrgName)
	}
}
