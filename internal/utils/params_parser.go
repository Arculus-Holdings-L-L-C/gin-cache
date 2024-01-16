package utils

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func ParameterParser(c *gin.Context) map[string]interface{} {
	m := make(map[string]interface{})
	split := strings.Split(c.FullPath(), `/`)
	params := strings.Split(c.Request.URL.Path, `/`)
	for i, preKey := range split {
		if strings.Contains(preKey, ":") {
			key := strings.ReplaceAll(preKey, ":", "")
			m[key] = params[i]
		}
	}
	if c.Request.Method == http.MethodGet {
		queryParams := c.Request.URL.Query()
		for key, val := range queryParams {
			m[key] = strings.Join(val, ",")
		}
	} else if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut {
		postMap := make(map[string]interface{})
		err := c.ShouldBindBodyWith(&postMap, binding.JSON)
		if err == nil {
			for key, val := range postMap {
				m[key] = val
			}
		}
	}
	return m
}
