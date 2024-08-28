package memcache

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/Arculus-Holdings-L-L-C/gin-cache/internal/entity"
	"github.com/Arculus-Holdings-L-L-C/gin-cache/pkg/define"
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

func (m *memoryHandler) Load(_ context.Context, key string) define.Response {
	load, ok := m.cacheStore.Load(key)
	if ok {
		item := load.(entity.CacheItem)
		if item.ExpireAt.UnixNano() < time.Now().UnixNano() {
			m.cacheStore.Delete(key)
			return define.Response{}
		}
		return item.Rsp
	}
	return define.Response{}
}

func (m *memoryHandler) Set(_ context.Context, key string, rsp define.Response, timeout time.Duration) {
	if timeout > 0 {
		m.cacheStore.Store(key, entity.CacheItem{
			Rsp:      rsp,
			CreateAt: time.Now(),
			ExpireAt: time.Now().Add(timeout),
		})
	} else {
		m.cacheStore.Store(key, entity.CacheItem{
			Rsp:      rsp,
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
