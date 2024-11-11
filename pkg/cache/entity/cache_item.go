package entity

import (
	"time"

	"github.com/Arculus-Holdings-L-L-C/gin-cache/pkg/define"
)

type CacheItem struct {
	Rsp      define.Response
	CreateAt time.Time
	ExpireAt time.Time
	Hits     uint64
}
