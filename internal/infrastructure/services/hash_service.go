package services

import (
	"context"
	"errors"

	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type ArgonParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

type ArgonHashService struct {
	params ArgonParams
}

func NewArgonHashService() *ArgonHashService {
	return &ArgonHashService{
		params: ArgonParams{
			Memory:      64 * 1024, // 64MB
			Iterations:  3,
			Parallelism: 2,
			SaltLength:  16,
			KeyLength:   32,
		},
	}
}

func (s *ArgonHashService) HashPassword(ctx context.Context, password string) (string, error) {
	// Generate a random salt
	salt := make([]byte, s.params.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// Hash the password
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		s.params.Iterations,
		s.params.Memory,
		s.params.Parallelism,
		s.params.KeyLength,
	)

	// Format the hash with parameters for future verification
	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	// Format: $argon2id$v=19$m=65536,t=3,p=2$<salt>$<hash>
	fullHash := fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		s.params.Memory,
		s.params.Iterations,
		s.params.Parallelism,
		encodedSalt,
		encodedHash,
	)

	return fullHash, nil
}

func (s *ArgonHashService) VerifyPassword(ctx context.Context, hashedPassword, password string) (bool, error) {
	// Parse the hashed password
	parts := strings.Split(hashedPassword, "$")
	if len(parts) != 6 {
		return false, errors.New("invalid hash format")
	}

	if parts[1] != "argon2id" {
		return false, errors.New("unsupported hash algorithm")
	}

	// Parse the parameters
	var params ArgonParams
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &params.Memory, &params.Iterations, &params.Parallelism)
	if err != nil {
		return false, err
	}

	// Extract salt
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	// Extract hash
	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}
	params.KeyLength = uint32(len(hash))
	params.SaltLength = uint32(len(salt))

	// Compute hash with the same parameters
	computedHash := argon2.IDKey(
		[]byte(password),
		salt,
		params.Iterations,
		params.Memory,
		params.Parallelism,
		params.KeyLength,
	)

	// Time-constant comparison to avoid timing attacks
	return compareHashAndPassword(hash, computedHash), nil
}

// compareHashAndPassword compares two hashes in a time-constant manner
func compareHashAndPassword(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	var result byte
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}

	return result == 0
}
