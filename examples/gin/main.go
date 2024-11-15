package main

import (
	"fmt"

	"github.com/Arculus-Holdings-L-L-C/gin-cache/cmd/startup"
	"github.com/Arculus-Holdings-L-L-C/gin-cache/pkg/define"
	"github.com/gin-gonic/gin"
)

func main() {
	cache, _ := startup.MemCache()
	r := gin.Default()

	r.GET("/ping", cache.Func(
		define.Caching{
			Cacheable: define.Cacheable{
				// params["id"] 是请求数据, 来自于query 或者 post data, 例如: `/?id=1`, 缓存将会生成为: `anson:id:1`
				GenKey: func(c *gin.Context) string {
					return fmt.Sprintf("anson:id:%s", c.Query("id"))
				},
			},
		},
		func(c *gin.Context) {
			query, _ := c.GetQuery("id")

			c.JSON(200, gin.H{
				"message": query, // 返回数据将会被缓存
			})
		},
	))

	r.GET("/pings", func(c *gin.Context) {
		query, _ := c.GetQuery("id")

		c.JSON(200, gin.H{
			"message": query, // 返回数据将会被缓存
		})
	})

	r.Run(":80")
}
