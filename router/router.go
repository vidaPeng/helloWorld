package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"net/http"
)

func SetupRoutes(r *gin.Engine) {
	var tracer = otel.Tracer("hello_peng")
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/ping"}, //配置跳过 /ping 的日志
	}), gin.Recovery())
	r.GET("/ping", func(c *gin.Context) {
		ctx := otel.GetTextMapPropagator().Extract(
			// 从 header 里面自动提取 trace 链路相关数据
			c.Request.Context(),
			propagation.HeaderCarrier(c.Request.Header),
		)

		_, span := tracer.Start(ctx, "ping")
		defer span.End()

		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/test", func(c *gin.Context) {
		// 从 Gin 的 request 中获取上下文
		ctx, span := tracer.Start(context.Background(), "HelloTest")
		defer span.End()

		// 创建带 traceparent header 的 HTTP 请求
		req, err := http.NewRequestWithContext(ctx, "GET", "http://test-oci-hello-peng.pixocial.com/ping", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 用 otelhttp 自动注入 traceparent
		client := http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		}
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
