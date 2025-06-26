package pkg

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func InitTracer() {
	ctx := context.Background()

	exp, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint("opentelemetry.trace.grpc.pix.com"),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("helloPeng"), // ← 用新 API 写法
		)),
	)
	otel.SetTracerProvider(tp)

	// ⭐️⭐️ 关键：告诉全局“用 W3C TraceContext + Baggage 解析/注入 header”
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{}, // 解析 traceparent / 注入 traceparent
		),
	)
}
