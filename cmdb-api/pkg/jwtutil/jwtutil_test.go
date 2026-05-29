package jwtutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAndParseToken(t *testing.T) {
	secret := "test-secret-key-for-unit-tests"
	token, err := GenerateToken(1, "admin", secret, 24)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := ParseToken(token, secret)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), claims.UserID)
	assert.Equal(t, "admin", claims.Username)
}

func TestParseInvalidToken(t *testing.T) {
	_, err := ParseToken("invalid.token.here", "secret")
	assert.Error(t, err)
}

func TestParseTokenWithWrongSecret(t *testing.T) {
	secret := "correct-secret"
	token, _ := GenerateToken(1, "admin", secret, 24)
	_, err := ParseToken(token, "wrong-secret")
	assert.Error(t, err)
}
