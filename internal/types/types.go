package types

import (
	"time"
)

type URLMapping struct {
	Id        int64     `json:"id"`
	Url       string    `json:"url"`
	Alias     string    `json:"alias"`
	CreatedAt time.Time `json:"created_at"`
}
