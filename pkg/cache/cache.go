package cache

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/Arculus-Holdings-L-L-C/gin-cache/pkg"
	"github.com/Arculus-Holdings-L-L-C/gin-cache/pkg/define"
	"github.com/gin-gonic/gin"
)

var bodyBytesKey = "bodyIO"

type Cache interface {
	Load(ctx context.Context, key string) define.Response
	Set(ctx context.Context, key string, rsp define.Response, timeout time.Duration)
	DoEvict(ctx context.Context, keys []string)
}

type CacheHandler struct {
	Cache      Cache
	OnCacheHit []define.OnCacheHit // 命中缓存钩子 优先级低
}

func (cache *CacheHandler) Load(ctx context.Context, key string) define.Response {
	return cache.Cache.Load(ctx, key)
}

func (cache *CacheHandler) Set(ctx context.Context, key string, rsp define.Response, timeout time.Duration) {
	cache.Cache.Set(ctx, key, rsp, timeout)
}

func (cache *CacheHandler) DoEvict(ctx context.Context, keys []string) {
	cache.Cache.DoEvict(ctx, keys)
}

func New(c Cache, onCacheHit ...define.OnCacheHit) *CacheHandler {
	return &CacheHandler{c, onCacheHit}
}

// Handler for startup
func (cache *CacheHandler) Handler(config define.Caching, next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		doEvict := len(config.Evict) > 0
		ctx := context.Background()
		var key string
		var cacheRsp define.Response

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

		key = cache.getCacheKey(config.Cacheable, c)
		if key != "" {
			cacheRsp = cache.Cache.Load(ctx, key)
		}

		if cacheRsp.Body == "" {
			refreshBodyData(c)
			next(c)
			if code := c.Writer.Status(); code != http.StatusOK &&
				!slices.Contains(config.CacheErrorCodes, code) {
				return
			}
			refreshBodyData(c)
		} else {
			cache.doCacheHit(c, config, cacheRsp)
		}

		if doEvict {
			cache.doCacheEvict(ctx, c, config.Evict...)
		}

		if cacheRsp = cache.Cache.Load(ctx, key); cacheRsp.Body == "" {
			cacheRsp.Body = c.Writer.(*pkg.ResponseBodyWriter).Body.String()
			cacheRsp.Header = c.Writer.Header()
			cacheRsp.Status = c.Writer.Status()
			cache.Cache.Set(ctx, key, cacheRsp, config.Cacheable.CacheTime)
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

func (cache *CacheHandler) doCacheHit(c *gin.Context, config define.Caching, r define.Response) {
	if config.Cacheable.OnCacheHit != nil {
		config.Cacheable.OnCacheHit[0](c, r)
		c.Abort()
		return
	}

	if len(cache.OnCacheHit) > 0 {
		cache.OnCacheHit[0](c, r)
		c.Abort()
		return
	}

	for key, values := range r.Header {
		for _, value := range values {
			c.Writer.Header().Add(key, value)
		}
	}
	if c.Writer.Header().Get("Content-Type") == "" {
		c.Writer.Header().Set("Content-Type", "application/json; Charset=utf-8")
	}
	c.String(r.Status, r.Body)
	c.Abort()
}

func refreshBodyData(c *gin.Context) {
	if c.Request.Body != nil {
		bodyStr, exists := c.Get(bodyBytesKey)
		if exists {
			c.Request.Body = io.NopCloser(bytes.NewReader(bodyStr.([]byte)))
		}
	}
}
