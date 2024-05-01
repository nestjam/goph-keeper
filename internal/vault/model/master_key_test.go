package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nestjam/goph-keeper/internal/utils"
)

func TestMasterKey_Seal(t *testing.T) {
	t.Run("seal data key", func(t *testing.T) {
		key, _ := utils.GenerateRandomAES256Key()
		sut := NewMasterKey(key)
		unsealed, _ := NewDataKey()

		_, err := sut.Seal(unsealed)

		require.NoError(t, err)
	})
	t.Run("seal with invalid key size", func(t *testing.T) {
		const keySize = 8
		key, _ := utils.GenerateRandom(keySize)
		sut := NewMasterKey(key)
		unsealed, _ := NewDataKey()

		_, err := sut.Seal(unsealed)

		require.Error(t, err)
	})
}

func TestMasterKey_Unseal(t *testing.T) {
	t.Run("unseal sealed data key", func(t *testing.T) {
		key, _ := utils.GenerateRandomAES256Key()
		sut := NewMasterKey(key)
		original, _ := NewDataKey()
		sealed, err := sut.Seal(original)
		require.NoError(t, err)

		unsealed, err := sut.Unseal(sealed)

		require.NoError(t, err)
		assert.Equal(t, original, unsealed)
	})
	t.Run("unseal invalid data", func(t *testing.T) {
		key, _ := utils.GenerateRandomAES256Key()
		sut := NewMasterKey(key)
		original, _ := NewDataKey()

		_, err := sut.Unseal(original)

		require.Error(t, err)
	})
}
