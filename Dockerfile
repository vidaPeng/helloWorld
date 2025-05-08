# ───────── build stage ─────────
FROM golang:1.24 AS builder

# 1️⃣ 安装 otel CLI
RUN curl -fsSL https://cdn.jsdelivr.net/gh/alibaba/opentelemetry-go-auto-instrumentation@main/install.sh | bash

# 2️⃣ 加入 PATH（顶格写，去掉中文注释或放到独立行）
ENV PATH="/usr/local/bin:$PATH"

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .

# 3️⃣ otel 配置（可选）
RUN otel set -verbose -log=/dev/stdout

# 4️⃣ 用 otel 构建
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    otel go build -o /server ./main.go

# ───────── runtime stage ─────────
FROM alpine:latest
WORKDIR /app
COPY --from=builder /server /app/server

# 5️⃣ 标准 OTel 变量
ENV OTEL_EXPORTER_OTLP_ENDPOINT=http://opentelemetry-collector.observable.svc:4317 \
    OTEL_TRACES_EXPORTER=otlp \
    OTEL_SERVICE_NAME=my-go-service \
    OTEL_EXPORTER_OTLP_INSECURE=true

ENTRYPOINT ["/app/server"]