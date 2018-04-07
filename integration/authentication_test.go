// +build integration

package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthenticate(t *testing.T) {
	assert.Equal(t, organizationName, org.Name)
	assert.NotEmpty(t, org.UUID)
	assert.NotEmpty(t, org.Scopes)
}
