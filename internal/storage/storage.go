package storage

import (
	"github.com/d1agnozzz/url-shortener/internal/aliaser"
	"github.com/d1agnozzz/url-shortener/internal/types"
	"github.com/d1agnozzz/url-shortener/internal/urlsanitizer"
)

type Storage interface {
	CreateURLMapping(url urlsanitizer.SanitizedURL) (types.URLMapping, error)
	GetByAlias(alias aliaser.Alias) (types.URLMapping, error)
}
