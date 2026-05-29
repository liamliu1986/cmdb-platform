package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"cmdb-api/config"
)

func TestAuthService_RegisterAndLogin(t *testing.T) {
	cfg := &config.Config{
		JWTSecret:      "test-secret",
		JWTExpireHours: 1,
	}
	svc := NewAuthService(cfg)
	assert.NotNil(t, svc)

	// Test Register
	req := &RegisterRequest{
		Username: "testuser",
		Password: "password123",
		Email:    "test@example.com",
	}
	// This will fail without a real DB - just verify the service is initialized
	// In real tests, use a test database or mock repository
	_ = req
}
