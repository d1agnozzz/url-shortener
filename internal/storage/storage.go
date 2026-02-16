package storage

import (
	"context"
	"fmt"

	"github.com/d1agnozzz/url-shortener/internal/types"
)

type Storage interface {
	InsertURLMapping(ctx context.Context, urlMapping types.URLMapping) error
	GetByAlias(ctx context.Context, alias string) (*types.URLMapping, error)
}

type InsertDuplicateError types.URLMapping

func (s InsertDuplicateError) String() string {
	return fmt.Sprintf("%v; %v; %v; %v;", s.Id, s.Url, s.Alias, s.CreatedAt)
}

func (s InsertDuplicateError) Error() string {
	return "duplicate insertion: " + s.String()
}

type NotFound string

func (s NotFound) Error() string {
	return "not found by: " + string(s)
}
