// Copyright 2023 Intrinsic Innovation LLC

package auth

import (
	"fmt"
	"net/http"
	"sync"
)

var (
	tsFactoryInit sync.Once
	tsFactory     *sharedTokenSourceFactory
)

func getSharedTokenSourceFactory() *sharedTokenSourceFactory {
	tsFactoryInit.Do(func() {
		tsFactory = &sharedTokenSourceFactory{
			c:   http.DefaultClient,
			tsm: make(map[string]*cachedTokenSource),
		}
	})
	return tsFactory
}

// sharedTokenSourceFactory is a factory for cached token sources.
type sharedTokenSourceFactory struct {
	c     *http.Client
	tsm   map[string]*cachedTokenSource
	smMtx sync.Mutex
}

// cacheKey returns a cache key for the given token service address and API key.
// We re-use a token source if it is based on the same accounts token exchange service address and
// the same API key.
func cacheKey(fsAddr, apiKey string) string {
	return fmt.Sprintf("%s-%s", fsAddr, apiKey)
}

// LoadOrNew returns a cached token source for the given API key and token service address.
// It is thread-safe and can be used by multiple goroutines.
func (s *sharedTokenSourceFactory) LoadOrNew(fsAddr, apiKey string) (*cachedTokenSource, error) {
	k := cacheKey(fsAddr, apiKey)
	// lock the token source map
	s.smMtx.Lock()
	defer s.smMtx.Unlock()
	// lookup the token source for the given API key
	ts, ok := s.tsm[k]
	if ok {
		return ts, nil
	}
	// if there is no token source for the given address and API key, create one and store it
	tsc, err := NewTokensServiceClient(s.c, fsAddr)
	if err != nil {
		return nil, fmt.Errorf("cannot create token exchange: %w", err)
	}
	nts := cachedTokenSource{
		tp:               tsc,
		apiKey:           apiKey,
		minTokenLifetime: defaultMinTokenLifetime,
	}
	s.tsm[k] = &nts
	return &nts, nil
}
