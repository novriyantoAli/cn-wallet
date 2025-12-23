package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserLevel_IsValid(t *testing.T) {
	t.Run("should return true for valid levels", func(t *testing.T) {
		assert.True(t, UserLevelUser.IsValid())
		assert.True(t, UserLevelProvider.IsValid())
		assert.True(t, UserLevelReseller.IsValid())
		assert.True(t, UserLevelAdmin.IsValid())
	})

	t.Run("should return false for invalid levels", func(t *testing.T) {
		assert.False(t, UserLevel("invalid").IsValid())
		assert.False(t, UserLevel("agent").IsValid())
		assert.False(t, UserLevel("superuser").IsValid())
		assert.False(t, UserLevel("").IsValid())
	})
}

func TestUserLevel_String(t *testing.T) {
	t.Run("should return string representation", func(t *testing.T) {
		assert.Equal(t, "user", UserLevelUser.String())
		assert.Equal(t, "provider", UserLevelProvider.String())
		assert.Equal(t, "reseller", UserLevelReseller.String())
		assert.Equal(t, "admin", UserLevelAdmin.String())
	})
}

func TestUserLevelConstants(t *testing.T) {
	t.Run("should have correct constant values", func(t *testing.T) {
		assert.Equal(t, UserLevel("user"), UserLevelUser)
		assert.Equal(t, UserLevel("provider"), UserLevelProvider)
		assert.Equal(t, UserLevel("reseller"), UserLevelReseller)
		assert.Equal(t, UserLevel("admin"), UserLevelAdmin)
	})
}
