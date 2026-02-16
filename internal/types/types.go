package types

import (
	"time"
)

type APIError struct {
	Error string `json:"api_error"`
}

type AliasResponse struct {
	Alias string `json:"alias"`
}

type UrlResponse struct {
	Url string `json:"url"`
}

type URLMapping struct {
	Id        int64     `json:"id"`
	Url       string    `json:"url"`
	Alias     string    `json:"alias"`
	CreatedAt time.Time `json:"created_at"`
}

type PostURLRequest struct {
	Url string `json:"url"`
}
