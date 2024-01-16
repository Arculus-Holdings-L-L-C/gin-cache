package memcache

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Arculus-Holdings-L-L-C/gin-cache/internal/entity"
)

// memoryHandler is private
type memoryHandler struct {
	cacheStore sync.Map
}

// NewMemoryHandler do new memory startup object
func NewMemoryHandler() *memoryHandler {
	memoryHandler := &memoryHandler{
		cacheStore: sync.Map{},
	}

	timer := time.NewTicker(time.Second * 30)

	go func() {
		for {
			<-timer.C
			memoryHandler.cacheStore.Range(func(key, value interface{}) bool {
				item := value.(entity.CacheItem)
				if item.ExpireAt.UnixNano() < time.Now().UnixNano() {
					memoryHandler.cacheStore.Delete(key)
				}
				return true
			})
		}
	}()

	return memoryHandler
}

func (m *memoryHandler) Load(_ context.Context, key string) (http.Header, string) {
	load, ok := m.cacheStore.Load(key)
	if ok {
		item := load.(entity.CacheItem)
		if item.ExpireAt.UnixNano() < time.Now().UnixNano() {
			m.cacheStore.Delete(key)
			return nil, ""
		}
		return item.Header, item.Value
	}
	return nil, ""
}

func (m *memoryHandler) Set(_ context.Context, key string, hdr http.Header, data string, timeout time.Duration) {
	if timeout > 0 {
		m.cacheStore.Store(key, entity.CacheItem{
			Value:    data,
			Header:   hdr,
			CreateAt: time.Now(),
			ExpireAt: time.Now().Add(timeout),
		})
	} else {
		m.cacheStore.Store(key, entity.CacheItem{
			Value:    data,
			Header:   hdr,
			CreateAt: time.Now(),
			ExpireAt: time.Now().Add(time.Hour * 1000000),
		})
	}

}

func (m *memoryHandler) DoEvict(_ context.Context, keys []string) {
	var evictKeys []string
	for _, key := range keys {
		isEndingStar := key[len(key)-1:]
		m.cacheStore.Range(func(keyInMap, _ interface{}) bool {
			// match *
			if isEndingStar == "*" {
				if strings.Contains(keyInMap.(string), strings.ReplaceAll(key, "*", "")) {
					evictKeys = append(evictKeys, keyInMap.(string))
				}
			} else {
				if keyInMap == key {
					evictKeys = append(evictKeys, key)
				}
			}
			return true
		})
	}

	for _, key := range evictKeys {
		m.cacheStore.Delete(key)
	}
}
