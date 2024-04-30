package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateRandom(t *testing.T) {
	t.Run("generate random", func(t *testing.T) {
		const wantSize = 8
		got, err := GenerateRandom(wantSize)

		require.NoError(t, err)
		assert.Equal(t, wantSize, len(got))
	})
	t.Run("zero size", func(t *testing.T) {
		const wantSize = 0
		got, err := GenerateRandom(wantSize)

		require.NoError(t, err)
		assert.Equal(t, wantSize, len(got))
	})
	t.Run("size is less zero", func(t *testing.T) {
		const size = -1

		assert.Panics(t, func() {
			_, _ = GenerateRandom(size)
		})
	})
}

func TestGenerateRandomAES256Key(t *testing.T) {
	const wantKeySize = 32
	got, err := GenerateRandomAES256Key()

	require.NoError(t, err)
	assert.Equal(t, wantKeySize, len(got))
}
