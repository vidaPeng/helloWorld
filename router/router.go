package router

import (
	"context"
	"fmt"
	"github.com/PixDevopsSre/helloWorld/pkg"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"io"
	"net/http"
	"time"
)

var tracer = otel.Tracer("hello_peng")

func SetupRoutes(r *gin.Engine) {
	r.GET("/getList", func(c *gin.Context) {
		ctx, span := tracer.Start(c.Request.Context(), "getList")
		defer span.End()

		pkg.InfoTrace(ctx, "这是一个会把 traceID 信息打印出来的日志组件")

		body, err := httpRequest(c.Request.Context(), "test-oci-hello-peng.pixocial.com", "/checkSqlList", "GET")
		if err != nil {
			pkg.InfoTrace(c, "httpRequest error")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// 正常响应
		c.JSON(http.StatusOK, gin.H{"message": string(body)})
	})

	r.GET("/checkSqlList", func(c *gin.Context) {
		ctx, span := tracer.Start(c.Request.Context(), "checkSqlList")
		defer span.End()

		// 这里模拟数据库的调用
		gormFunc(ctx)

		// 正常响应
		c.JSON(http.StatusOK, gin.H{"traceID": span.SpanContext().TraceID().String()})
	})

	r.GET("/getTimeout", func(c *gin.Context) {
		ctx, span := tracer.Start(c.Request.Context(), "getTimeout")
		defer span.End()

		body, err := httpRequest(ctx, "test-oci-hello-peng.pixocial.com", "/clientTimeout", "GET")
		if err != nil {
			pkg.InfoTrace(c, "httpRequest error")
			// 这里可以把 span 的状态设置为错误，这了可以对 opentelemetry 做一个封装，让使用更加方便
			span.RecordError(err)
			span.SetStatus(codes.Error, "httpRequest failed")

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// 正常响应
		c.JSON(http.StatusOK, gin.H{"message": string(body)})
	})

	r.GET("/clientTimeout", func(c *gin.Context) {
		_, span := tracer.Start(c.Request.Context(), "clientTimeout")
		defer span.End()

		// 这里模拟一个超时的操作
		time.Sleep(6 * time.Second)

		// 正常响应
		c.JSON(http.StatusOK, gin.H{"traceID": span.SpanContext().TraceID().String()})
	})

	r.GET("/getApisixTimeout", func(c *gin.Context) {
		ctx, span := tracer.Start(c.Request.Context(), "getApisixTimeout")
		defer span.End()

		body, err := httpRequest(ctx, "test-oci-hello-peng-timeout.pixocial.com", "/ApisixTimeout", "GET")
		if err != nil {
			pkg.InfoTrace(c, "httpRequest error")
			span.RecordError(err)
			span.SetStatus(codes.Error, "httpRequest failed")

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// 正常响应
		c.JSON(http.StatusOK, gin.H{"message": string(body)})
	})

	r.GET("/ApisixTimeout", func(c *gin.Context) {
		_, span := tracer.Start(c.Request.Context(), "ApisixTimeout")
		defer span.End()

		// 这里模拟一个调用 APISIX 的操作
		time.Sleep(6 * time.Second)

		// 正常响应
		c.JSON(http.StatusOK, gin.H{"traceID": span.SpanContext().TraceID().String()})
	})
}

func gormFunc(ctx context.Context) {
	_, gromSapn := tracer.Start(ctx, "gorm", trace.WithAttributes(
		attribute.String("Sql", "select * from user where id = 1"),
	))
	defer gromSapn.End()

	time.Sleep(1 * time.Second)
}

var client = http.Client{
	// ✨ 用 otelhttp 自动注入 traceparent , 这一步是重点
	Transport: otelhttp.NewTransport(http.DefaultTransport),
	Timeout:   5 * time.Second, // 设置超时时间
}

func httpRequest(ctx context.Context, host, path, method string) ([]byte, error) {
	// 从 Gin 的 request 中获取上下文
	ctx, span := tracer.Start(
		// 使用 tracer 创建一个新的 span，需要把 context 传入
		ctx,
		"httpRequest",
		// 添加一些属性到 span 中
		trace.WithAttributes(
			attribute.String("hello", "peng"),
			attribute.Int64("version", 1),
			attribute.Bool("cache_hit", false),
		),
	)
	defer span.End()

	url := fmt.Sprintf("http://%s%s", host, path)
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %s, body: %s", resp.Status, body)
	}
	return body, nil
}
