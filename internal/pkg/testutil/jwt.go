package testutil

import (
	"time"

	"github.com/novriyantoAli/cn-wallet/internal/config"
	"github.com/novriyantoAli/cn-wallet/internal/pkg/jwt"
)

// CreateTestJWTManager creates a JWT manager for testing
func CreateTestJWTManager() *jwt.JWTManager {
	cfg := &config.Config{
		JWT: config.JWTConfig{
			SecretKey: "test-secret-key-for-testing-only",
			Expiry:    24 * time.Hour,
		},
	}
	return jwt.NewJWTManager(cfg)
}

// CreateValidJWTToken creates a valid JWT token for testing
func CreateValidJWTToken(manager *jwt.JWTManager, userID uint) string {
	token, err := manager.GenerateToken(userID, "test@example.com", "user")
	if err != nil {
		panic("failed to generate test token: " + err.Error())
	}
	return token
}
