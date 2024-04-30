package model

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	sut := &Secret{
		ID:   uuid.New(),
		IV:   []byte("init vector"),
		Data: []byte("data"),
	}

	got := sut.Copy()

	assert.True(t, sut != got)
	assert.Equal(t, sut, got)
}
