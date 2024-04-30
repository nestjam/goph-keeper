package model

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDataKeyCopy(t *testing.T) {
	want, err := NewDataKey()
	want.ID = uuid.New()
	want.EncryptedSize = 1024
	require.NoError(t, err)

	got := want.Copy()

	assert.Equal(t, want, got)
}
