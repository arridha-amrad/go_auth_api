package services_test

import (
	"my-go-api/internal/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordService_HashAndVerify(t *testing.T) {
	passwordService := services.NewPasswordService()

	t.Run("It should hash password without error", func(t *testing.T) {
		hash, err := passwordService.Hash("secret123")
		assert.NoError(t, err)
		assert.NotEmpty(t, hash)
	})

	t.Run("It should verify the correct password", func(t *testing.T) {
		hash, _ := passwordService.Hash("secret123")
		err := passwordService.Verify(hash, "secret123")
		assert.NoError(t, err)
	})

	t.Run("It should fail verification with wrong password", func(t *testing.T) {
		hash, _ := passwordService.Hash("secret123")
		err := passwordService.Verify(hash, "wrongpassword")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "hashedPassword")
	})
}
