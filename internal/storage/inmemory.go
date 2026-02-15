package storage

import (
	"fmt"
	"github.com/d1agnozzz/url-shortener/internal/types"
	"github.com/d1agnozzz/url-shortener/internal/urlsanitizer"
	"sync"
	"time"
)

type inMemStorage struct {
	mu          sync.RWMutex
	urlMappings map[string]types.URLMapping
	idGen       int64
	config      StorageConfig
}

func NewInMemoryStorage(config StorageConfig) inMemStorage {
	res := make(map[string]types.URLMapping)
	return inMemStorage{
		urlMappings: res,
		config:      config,
	}
}

func (s *inMemStorage) CreateURLMapping(url urlsanitizer.SanitizedURL) (*types.URLMapping, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for attempt := 0; attempt <= s.config.maxCollisions; attempt++ {
		salted := fmt.Sprintf("%s::%d", url.String(), attempt)
		alias := s.config.aliaser.GenerateByStr(salted)

		mapping, exists := s.urlMappings[alias.String()]

		if !exists {
			newMapping := types.URLMapping{
				Id:        s.idGen,
				Url:       url.String(),
				Alias:     alias.String(),
				CreatedAt: time.Now().UTC(),
			}
			s.urlMappings[alias.String()] = newMapping
			s.idGen++
			return &newMapping, nil
		}

		if mapping.Url == url.String() {
			return &mapping, nil
		}
	}

	return nil, fmt.Errorf("too much collisions, reject")

}

func (s *inMemStorage) GetByAlias(alias string) (*types.URLMapping, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	url, exists := s.urlMappings[alias]

	if !exists {
		return nil, fmt.Errorf("url not found by alias '%s'", alias)
	}

	return &url, nil

}
