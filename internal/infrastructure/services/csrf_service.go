package services

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"sync"
	"time"
)

type CSRFService struct {
	tokens     map[string]time.Time
	tokenMutex sync.RWMutex
	maxAge     time.Duration
}

func NewCSRFService(maxAge time.Duration) *CSRFService {
	service := &CSRFService{
		tokens: make(map[string]time.Time),
		maxAge: maxAge,
	}

	// Start a goroutine to clean up expired tokens
	go service.cleanupExpiredTokens()

	return service
}

func (s *CSRFService) cleanupExpiredTokens() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.tokenMutex.Lock()
		now := time.Now()
		for token, expiry := range s.tokens {
			if now.After(expiry) {
				delete(s.tokens, token)
			}
		}
		s.tokenMutex.Unlock()
	}
}

func (s *CSRFService) GenerateToken(r *http.Request) string {
	// Generate a random token
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}
	token := base64.StdEncoding.EncodeToString(bytes)

	// Store the token with its expiry time
	s.tokenMutex.Lock()
	s.tokens[token] = time.Now().Add(s.maxAge)
	s.tokenMutex.Unlock()

	return token
}

func (s *CSRFService) ValidateToken(r *http.Request, token string) bool {
	if token == "" {
		return false
	}

	s.tokenMutex.RLock()
	expiry, exists := s.tokens[token]
	s.tokenMutex.RUnlock()

	if !exists || time.Now().After(expiry) {
		return false
	}

	// Remove the token after use (one-time use)
	s.tokenMutex.Lock()
	delete(s.tokens, token)
	s.tokenMutex.Unlock()

	return true
}
