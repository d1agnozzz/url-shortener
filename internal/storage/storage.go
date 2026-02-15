package storage

import (
	"github.com/d1agnozzz/url-shortener/internal/aliasgenerator"
	"github.com/d1agnozzz/url-shortener/internal/types"
	"github.com/d1agnozzz/url-shortener/internal/urlsanitizer"
)

type Storage interface {
	CreateURLMapping(url urlsanitizer.SanitizedURL) (types.URLMapping, error)
	GetByAlias(alias aliasgenerator.Alias) (types.URLMapping, error)
}
