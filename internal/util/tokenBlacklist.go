package util

import (
	"sync"
	"time"
)

var tokenBlacklist = struct {
	sync.RWMutex
	data map[string]time.Time
}{
	data: make(map[string]time.Time),
}

func RevokeToken(token string, expiresAt time.Time) {
	if token == "" {
		return
	}

	tokenBlacklist.Lock()
	tokenBlacklist.data[token] = expiresAt
	tokenBlacklist.Unlock()
}

func IsTokenRevoked(token string) bool {
	if token == "" {
		return false
	}

	now := time.Now()

	tokenBlacklist.RLock()
	expiresAt, ok := tokenBlacklist.data[token]
	tokenBlacklist.RUnlock()

	if !ok {
		return false
	}

	// Lazy cleanup of already-expired revoked tokens.
	if now.After(expiresAt) {
		tokenBlacklist.Lock()
		delete(tokenBlacklist.data, token)
		tokenBlacklist.Unlock()
		return false
	}

	return true
}
