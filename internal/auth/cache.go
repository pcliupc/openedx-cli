package auth

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/openedx/cli/internal/config"
)

// expiryMargin is the safety margin applied before the actual token expiry.
// A cached token is considered expired if fewer than this many seconds remain.
const expiryMargin = 60 * time.Second

// cachedToken holds a token together with its absolute expiry time.
type cachedToken struct {
	token     *Token
	expiresAt time.Time
}

// isExpired returns true if the cached token has expired or will expire within
// the safety margin.
func (ct *cachedToken) isExpired(now time.Time) bool {
	return now.After(ct.expiresAt.Add(-expiryMargin))
}

// CachingTokenProvider wraps an inner TokenProvider with an in-memory cache.
// Tokens are stored per profile key and returned if they have not yet expired
// (with a configurable safety margin). Expired or missing tokens trigger a
// refresh via the inner provider.
//
// All operations are protected by a sync.Mutex for thread safety.
type CachingTokenProvider struct {
	inner   TokenProvider
	mu      sync.Mutex
	tokens  map[string]*cachedToken
	clockFn func() time.Time
}

// NewCachingTokenProvider creates a new caching wrapper around the given
// TokenProvider. The clockFn parameter allows injecting a custom clock for
// testing; if nil, time.Now is used.
func NewCachingTokenProvider(inner TokenProvider, clockFn func() time.Time) *CachingTokenProvider {
	if clockFn == nil {
		clockFn = time.Now
	}
	return &CachingTokenProvider{
		inner:   inner,
		tokens:  make(map[string]*cachedToken),
		clockFn: clockFn,
	}
}

// profileKey generates a unique cache key from the profile. It uses the
// TokenURL and ClientIDEnv fields, which together identify a unique credential
// set.
func profileKey(profile config.Profile) string {
	return fmt.Sprintf("%s|%s", profile.TokenURL, profile.ClientIDEnv)
}

// Token returns a cached token for the given profile if it is still valid, or
// fetches a new one via the inner provider and caches it.
func (c *CachingTokenProvider) Token(ctx context.Context, profile config.Profile) (*Token, error) {
	key := profileKey(profile)

	c.mu.Lock()
	defer c.mu.Unlock()

	now := c.clockFn()

	if cached, ok := c.tokens[key]; ok && !cached.isExpired(now) {
		return cached.token, nil
	}

	token, err := c.inner.Token(ctx, profile)
	if err != nil {
		return nil, err
	}

	c.tokens[key] = &cachedToken{
		token:     token,
		expiresAt: now.Add(time.Duration(token.ExpiresIn) * time.Second),
	}

	return token, nil
}
