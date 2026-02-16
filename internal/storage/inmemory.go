package storage

import (
	"context"
	"sync"
	"time"

	"github.com/d1agnozzz/url-shortener/internal/types"
	"github.com/jackc/pgx/v5"
)

type inMemStorage struct {
	mu          sync.RWMutex
	urlMappings map[string]types.URLMapping
	idGen       int64
}

func NewInMemoryStorage() Storage {
	res := make(map[string]types.URLMapping)
	return &inMemStorage{
		urlMappings: res,
	}
}

func (s *inMemStorage) InsertURLMapping(ctx context.Context, urlMapping types.URLMapping) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	urlMapping.Id = s.idGen
	urlMapping.CreatedAt = time.Now()

	// check  by alias
	_, exists := s.urlMappings[urlMapping.Alias]

	if exists {
		return InsertDuplicateError(urlMapping)
	}

	// check duplicate by url
	// TODO: fix linear scan
	for _, v := range s.urlMappings {
		if v.Url == urlMapping.Url {
			return InsertDuplicateError(urlMapping)
		}
	}

	s.urlMappings[urlMapping.Alias] = urlMapping
	s.idGen++

	return nil
}

func (s *inMemStorage) GetByAlias(ctx context.Context, alias string) (*types.URLMapping, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	url, exists := s.urlMappings[alias]

	if !exists {
		return nil, pgx.ErrNoRows
	}

	return &url, nil
}
