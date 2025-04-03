package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
)

type TokenService struct{}

func NewTokenService() *TokenService {
	return &TokenService{}
}

// GenerateToken generates a random token using base32 encoding (no padding)
func (s *TokenService) GenerateToken(ctx context.Context) (string, error) {
	// Generate 20 random bytes (160 bits)
	bytes := make([]byte, 20)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Encode using base32 (lowercase, no padding)
	encoder := base32.NewEncoding("abcdefghijklmnopqrstuvwxyz234567").WithPadding(base32.NoPadding)
	token := encoder.EncodeToString(bytes)

	return token, nil
}

// HashToken hashes a token using SHA-256 and encodes it as a lowercase hex string
func (s *TokenService) HashToken(ctx context.Context, token string) (string, error) {
	// Create a SHA-256 hash
	hash := sha256.Sum256([]byte(token))

	// Encode as lowercase hex
	return hex.EncodeToString(hash[:]), nil
}
