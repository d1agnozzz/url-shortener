package storage

import (
	"fmt"
	"github.com/d1agnozzz/url-shortener/internal/aliaser"
	"github.com/d1agnozzz/url-shortener/internal/types"
	"github.com/d1agnozzz/url-shortener/internal/urlsanitizer"
	"sync"
	"time"
)

type inMemStorage struct {
	mu          sync.RWMutex
	urlMappings map[aliaser.Alias]types.URLMapping
	idGen       int64
	config      StorageConfig
}

func NewInMemoryStorage(config StorageConfig) inMemStorage {
	res := make(map[aliaser.Alias]types.URLMapping)
	return inMemStorage{
		urlMappings: res,
		config:      config,
	}
}

func (s *inMemStorage) CreateURLMapping(url urlsanitizer.SanitizedURL, alias aliaser.Alias) (types.URLMapping, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	alias := s.config.aliaser.GenerateByStr(url.String())

	mapping, exists := s.urlMappings[alias]

	if exists {
		return mapping, nil
	}

	new_mapping := types.URLMapping{
		Id:        s.idGen,
		Url:       url.String(),
		Alias:     alias.String(),
		CreatedAt: time.Now().UTC(),
	}

	s.urlMappings[alias] = new_mapping

	return new_mapping, nil
}

func (s *inMemStorage) GetByAlias(alias aliaser.Alias) (*types.URLMapping, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	url, exists := s.urlMappings[alias]

	if !exists {
		return nil, fmt.Errorf("url not found by alias '%s'", alias)
	}

	return &url, nil

}
