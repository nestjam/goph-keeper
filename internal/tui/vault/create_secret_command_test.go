package vault

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSecretCommand(t *testing.T) {
	sut := newCreateSecretCommand()

	got := sut.execute()

	_, ok := got.(createSecretRequestedMsg)
	assert.True(t, ok)
}
