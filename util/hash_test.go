package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashing(t *testing.T) {
	password := "mySecretPin"
	hash, err := HashPin(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.True(t, VerifyHash(password, hash))
}
