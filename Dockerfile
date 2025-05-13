# ───────── build stage ─────────
FROM golang:1.24 AS builder

# 1️⃣ 手动安装 otel CLI（跳过 sudo）
RUN curl -L \
      -o /usr/local/bin/otel \
      https://github.com/alibaba/opentelemetry-go-auto-instrumentation/releases/latest/download/otel-linux-amd64 \
 && chmod +x /usr/local/bin/otel

ENV PATH="/usr/local/bin:${PATH}"

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .

# 3️⃣ 用 otel 编译，自动插桩
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    otel go build -o /server ./main.go

# ───────── runtime stage ─────────
FROM alpine:latest
WORKDIR /app
COPY --from=builder /server /app/server

# 4️⃣ 设置 OTel 导出端点
ENV OTEL_EXPORTER_OTLP_ENDPOINT=http://opentelemetry-collector.observable.svc:4318 \
    OTEL_TRACES_EXPORTER=otlp \
    OTEL_SERVICE_NAME=helloWorld \
    OTEL_GO_AUTO_DISABLE_METRICS=true \
    OTEL_EXPORTER_OTLP_INSECURE=true

ENTRYPOINT ["/app/server"]