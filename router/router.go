package router

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/PixDevopsSre/helloWorld/pkg"
	"github.com/PixDevopsSre/helloWorld/proto"
	"github.com/gin-gonic/gin"
	"github.com/google/go-querystring/query"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	tracer = otel.Tracer("hello_peng")

	api_key = "ca7b7e2a75f385eee979643510e38dfc57e454100c69eaa5cea334d7fc240651"
)

func generateSignature(timestamp, apiKey string) string {
	keyMd5 := md5.Sum([]byte(apiKey))
	keyMd5Hex := strings.ToLower(hex.EncodeToString(keyMd5[:]))

	signMd5 := md5.Sum([]byte(timestamp + keyMd5Hex))
	signMd5Hex := strings.ToLower(hex.EncodeToString(signMd5[:]))
	return signMd5Hex
}

func SetupRoutes(r *gin.Engine) {
	r.POST("/getList", func(c *gin.Context) {
		// 2. 像定义 JSON 一样定义你的请求结构体，但使用 `url` tag
		type requestForm struct {
			ID        string `url:"id"`
			ApiToken  string `url:"api_token"`
			Timestamp string `url:"timestamp"`
		}

		// 3. 正常地填充结构体
		form := requestForm{
			ID: "7",
		}

		// --- 认证逻辑和之前一样 ---
		apiKey := "ca7b7e2a75f385eee979643510e38dfc57e454100c69eaa5cea334d7fc240651"
		ts := time.Now().Unix()
		form.Timestamp = strconv.FormatInt(ts, 10)
		form.ApiToken = generateSignature(form.Timestamp, apiKey)

		// 4. ✨ 最优雅的一步：使用库自动将结构体转换为 url.Values
		values, err := query.Values(form)
		if err != nil {
			// 一般这里不会出错，除非 tag 定义有问题
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode form: " + err.Error()})
			return
		}

		// 5. 将 url.Values 编码成字符串
		encodedData := values.Encode()

		// 后续代码完全不变
		body, err := httpRequest(
			c.Request.Context(),
			"localhost:51613",
			"/v1/workflow/get_workflow_id",
			http.MethodPost,
			strings.NewReader(encodedData),
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": string(body)})
	})

	r.GET("/getGrpc", func(c *gin.Context) {
		ctx, span := tracer.Start(c.Request.Context(), "getGrpc")
		defer span.End()

		// 这里模拟一个 gRPC 的调用
		conn, err := grpc.NewClient(
			"10.220.62.114:9000", // gRPC 服务地址
			// ✨ 使用 otelgrpc 自动注入 traceparent
			grpc.WithStatsHandler(otelgrpc.NewClientHandler()), // 设置 StatsHandler)
			grpc.WithTransportCredentials(insecure.NewCredentials()))

		if err != nil {
			pkg.InfoTrace(ctx, err.Error())
			span.RecordError(err)
			span.SetStatus(codes.Error, "grpc.NewClient failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		defer func() {
			if err := conn.Close(); err != nil {
				pkg.InfoTrace(ctx, err.Error())
				span.RecordError(err)
				span.SetStatus(codes.Error, "conn.Close failed")
			}
		}()

		client := proto.NewGreeterClient(conn)
		resp, err := client.SayHello(ctx, &proto.HelloRequest{
			Name: "hello",
		})
		if err != nil {
			pkg.InfoTrace(ctx, err.Error())
			span.RecordError(err)
			span.SetStatus(codes.Error, "client.SayHello failed")
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": resp.GetMessage()})
		if err != nil {
			return
		}
	})

}

var client = http.Client{
	// ✨ 用 otelhttp 自动注入 traceparent , 这一步是重点
	Transport: otelhttp.NewTransport(http.DefaultTransport),
	Timeout:   5 * time.Second, // 设置超时时间
}

func httpRequest(ctx context.Context, host, path, method string, data io.Reader) ([]byte, error) {
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
	req, err := http.NewRequestWithContext(ctx, method, url, data)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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
