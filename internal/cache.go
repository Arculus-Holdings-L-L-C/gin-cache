package internal

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Arculus-Holdings-L-L-C/gin-cache/pkg"
	"github.com/Arculus-Holdings-L-L-C/gin-cache/pkg/define"
	"github.com/gin-gonic/gin"
)

var bodyBytesKey = "bodyIO"

type Cache interface {
	Load(ctx context.Context, key string) (http.Header, string)
	Set(ctx context.Context, key string, hdr http.Header, data string, timeout time.Duration)
	DoEvict(ctx context.Context, keys []string)
}

type CacheHandler struct {
	Cache      Cache
	OnCacheHit define.CacheHitHook // 命中缓存钩子 优先级低
}

func (cache *CacheHandler) Load(ctx context.Context, key string) (http.Header, string) {
	return cache.Cache.Load(ctx, key)
}

func (cache *CacheHandler) Set(ctx context.Context, key string, hdr http.Header, data string, timeout time.Duration) {
	cache.Cache.Set(ctx, key, hdr, data, timeout)
}

func (cache *CacheHandler) DoEvict(ctx context.Context, keys []string) {
	cache.Cache.DoEvict(ctx, keys)
}

func New(c Cache, onCacheHit ...func(c *gin.Context, cacheValue string)) *CacheHandler {
	return &CacheHandler{c, onCacheHit}
}

// Handler for startup
func (cache *CacheHandler) Handler(caching define.Caching, next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {

		doEvict := len(caching.Evict) > 0
		ctx := context.Background()

		var key, cacheString string
		var cacheHeader http.Header

		if c.Request.Body != nil {
			body, err := io.ReadAll(c.Request.Body)
			if err != nil {
				body = []byte("")
			}
			c.Set(bodyBytesKey, body)
			c.Request.Body = io.NopCloser(bytes.NewReader(body))
		}

		// pointer 指向 writer, 重写 c.writer
		c.Writer = &pkg.ResponseBodyWriter{
			Body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}

		key = cache.getCacheKey(caching.Cacheable, c)
		if key != "" {
			cacheHeader, cacheString = cache.Cache.Load(ctx, key)
		}

		if cacheString == "" {
			refreshBodyData(c)
			next(c)
			if c.Writer.Status() != http.StatusOK {
				return
			}
			refreshBodyData(c)
		} else {
			cache.doCacheHit(c, caching, cacheHeader, cacheString)
		}

		if doEvict {
			cache.doCacheEvict(ctx, c, caching.Evict...)
		}

		if _, cacheString = cache.Cache.Load(ctx, key); cacheString == "" {
			s := c.Writer.(*pkg.ResponseBodyWriter).Body.String()
			cacheHeader = c.Writer.Header()
			cache.Cache.Set(ctx, key, cacheHeader, s, caching.Cacheable.CacheTime)
		}

	}
}

func (cache *CacheHandler) getCacheKey(cacheable define.Cacheable, c *gin.Context) string {
	return strings.ToLower(cacheable.GenKey(c))
}

func (cache *CacheHandler) doCacheEvict(ctx context.Context, c *gin.Context, cacheEvicts ...define.CacheEvict) {
	keys := make([]string, 0)
	for _, evict := range cacheEvicts {
		s := evict(c)
		if s != "" {
			keys = append(keys, strings.ToLower(s))
		}
	}

	if len(keys) > 0 {
		cache.Cache.DoEvict(ctx, keys)
	}
}

func (cache *CacheHandler) doCacheHit(ctx *gin.Context, caching define.Caching, cacheHeader http.Header, cacheValue string) {

	if caching.Cacheable.OnCacheHit != nil {
		caching.Cacheable.OnCacheHit[0](ctx, cacheValue)
		ctx.Abort()
		return
	}

	if len(cache.OnCacheHit) > 0 {
		cache.OnCacheHit[0](ctx, cacheValue)
		ctx.Abort()
		return
	}

	for key, values := range cacheHeader {
		for _, value := range values {
			ctx.Writer.Header().Add(key, value)
		}
	}
	if ctx.Writer.Header().Get("Content-Type") == "" {
		ctx.Writer.Header().Set("Content-Type", "application/json; Charset=utf-8")
	}
	ctx.String(http.StatusOK, cacheValue)
	ctx.Abort()
}

func refreshBodyData(c *gin.Context) {
	if c.Request.Body != nil {
		bodyStr, exists := c.Get(bodyBytesKey)
		if exists {
			c.Request.Body = io.NopCloser(bytes.NewReader(bodyStr.([]byte)))
		}
	}
}
