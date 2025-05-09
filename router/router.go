package router

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"net/http"
)

func SetupRoutes(r *gin.Engine) {

	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/ping"}, //配置跳过 /ping 的日志
	}), gin.Recovery())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/test", func(c *gin.Context) {
		// 从 Gin 的 request 中获取上下文
		tracer := otel.Tracer("hello_peng")
		_, span := tracer.Start(c.Request.Context(), "HelloTest")
		defer span.End()

		// 创建带 traceparent header 的 HTTP 请求
		req, err := http.NewRequestWithContext(c, "GET", "http://test-peng-cloud-bridge.pixocial.com/test", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 用 otelhttp 自动注入 traceparent
		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		// 正常响应
		c.JSON(http.StatusOK, gin.H{"message": "ok"})

	})
}
