package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"helloWorld/pkg"
	"io"
	"net/http"
)

func SetupRoutes(r *gin.Engine) {
	var tracer = otel.Tracer("hello_peng")
	r.GET("/ping", func(c *gin.Context) {
		//ctx := otel.GetTextMapPropagator().Extract(
		//	// 从 header 里面自动提取 trace 链路相关数据
		//	c.Request.Context(),
		//	propagation.HeaderCarrier(c.Request.Header),
		//)
		//
		//_, span := tracer.Start(ctx, "ping")
		//defer span.End()
		pkg.Info("ping")

		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/test1", func(c *gin.Context) {
		pkg.Info("test1")
		ctx, span := tracer.Start(context.Background(), "test1")
		defer span.End()
		pkg.InfoTrace(ctx, "test111111")

		request, err := http.NewRequest("GET", "http://test-oci-hello-peng.pixocial.com/test", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		request.Header.Set("pixcc_client", "123123123")

		// 使用 http.Client 发送请求
		client := &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		}
		resp, err := client.Do(request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()
		// 读取响应体
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// 正常响应
		c.JSON(http.StatusOK, gin.H{"message": string(body)})
	})

	r.GET("/test", func(c *gin.Context) {
		pkg.Info("test")
		// 从 Gin 的 request 中获取上下文
		//ctx, span := tracer.Start(
		//	context.Background(),
		//	"HelloTest",
		//	trace.WithAttributes(
		//		attribute.String("env", "dev"),
		//		attribute.Int64("version", 1),
		//		attribute.Bool("cache_hit", false),
		//	),
		//)
		//defer span.End()

		// 创建带 traceparent header 的 HTTP 请求

		req, err := http.NewRequestWithContext(context.Background(), "GET", "http://test-peng-cloud-bridge.pixocial.com/test", nil)
		//req, err := http.NewRequestWithContext(context.Background(), "GET", "http://test-oci-hello-peng.pixocial.com/ping", nil)
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

		// 读取响应体
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// 正常响应
		c.JSON(http.StatusOK, gin.H{"message": string(body)})
	})
}
