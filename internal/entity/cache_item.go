package entity

import (
	"net/http"
	"time"
)

type CacheItem struct {
	Value    string
	Header   http.Header
	CreateAt time.Time
	ExpireAt time.Time
	Hits     uint64
}
