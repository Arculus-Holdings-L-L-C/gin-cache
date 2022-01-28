package gincache

import (
	"context"
	"strings"
	"sync"
	"time"
)

// Schedule entity
type Schedule struct {
	Key   string
	Timer *time.Timer
}

// memoryHandler is private
type memoryHandler struct {
	cacheStore sync.Map
	cacheTime  time.Duration
	pubSub     chan Schedule
	schedules  map[string]*time.Timer
}

var mux sync.Mutex

// NewMemoryHandler do new memory cache object
func NewMemoryHandler(cacheTime time.Duration) *memoryHandler {
	return &memoryHandler{
		cacheStore: sync.Map{},
		cacheTime:  cacheTime,
		pubSub:     make(chan Schedule),
		schedules:  make(map[string]*time.Timer),
	}
}

func (m *memoryHandler) LoadCache(_ context.Context, key string) string {
	load, ok := m.cacheStore.Load(key)
	if ok {
		return load.(string)
	}
	return ""
}

func (m *memoryHandler) SetCache(ctx context.Context, key string, data string) {
	mux.Lock()
	m.cacheStore.Store(key, data)
	// timeout
	schedule := Schedule{Key: key, Timer: time.NewTimer(m.cacheTime)}
	m.schedules[key] = schedule.Timer
	go func(s *Schedule) {
		select {
		case <-s.Timer.C:
			m.DoCacheEvict(ctx, []string{s.Key})
		}
	}(&schedule)
	defer mux.Unlock()
}

func (m *memoryHandler) DoCacheEvict(_ context.Context, keys []string) []string {
	mux.Lock()
	evictKeys := []string{}
	for _, key := range keys {
		isEndingStar := key[len(key)-1:]
		m.cacheStore.Range(func(keyInMap, value interface{}) bool {
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
		timer := m.schedules[key]
		if timer != nil {
			timer.Stop()
		}
		delete(m.schedules, key)
	}
	defer mux.Unlock()
	return evictKeys
}
