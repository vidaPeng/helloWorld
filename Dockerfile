## 第一阶段：构建 Go 应用
#FROM golang:1.24 AS builder
#
#WORKDIR /app
#
## 复制 go.mod 和 go.sum 先做依赖缓存（提升构建效率）
#COPY go.mod ./
#RUN go mod download
#
## 再复制项目代码
#COPY . .
#
## 编译主程序（main.go 在 cmd 目录下）
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /server ./main.go
#
## 第二阶段：使用更小的基础镜像运行程序
#FROM alpine:latest
#
## 可选：使用非 root 用户更安全（记得开放文件或端口权限）
## RUN addgroup -S appgroup && adduser -S appuser -G appgroup
#
#WORKDIR /app
#
## 可选：创建配置文件目录（如果你的程序依赖 conf）
#RUN mkdir -p /app/conf
#
## 从构建阶段复制可执行文件
#COPY --from=builder /server .
#
## 可选：复制静态文件或配置文件（如果有）
## COPY ./conf /app/conf
#
## 可选：切换到非 root 用户运行（如上所设）
## USER appuser
#
## 启动入口
#ENTRYPOINT ["/app/server"]

# ──────────────────────────────
# 第一阶段：编译并自动插桩
# ──────────────────────────────
FROM golang:1.24 AS builder

# 1️⃣ 安装 Alibaba 的 otel 自动插桩 CLI
RUN curl -fsSL https://cdn.jsdelivr.net/gh/alibaba/opentelemetry-go-auto-instrumentation@main/install.sh | bash
ENV PATH="/usr/local/bin:${PATH}"   # 确保 `otel` 在 PATH 中

WORKDIR /app

# 2️⃣ 依赖缓存
COPY go.mod ./
RUN go mod download

# 3️⃣ 拷贝源码
COPY . .

# 4️⃣ （可选）配置 otel 行为：打开详细日志、写到 stdout
RUN otel set -verbose -log=/dev/stdout   # 需要更安静可删掉

# 5️⃣ **用 otel 前缀替代 go build**（自动注入追踪逻辑）
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    otel go build -o /server ./main.go

# ──────────────────────────────
# 第二阶段：极简运行时镜像
# ──────────────────────────────
FROM alpine:latest

WORKDIR /app
COPY --from=builder /server /app/server

# 6️⃣ 设置标准 OpenTelemetry 环境变量
# 如果你的 Collector 用 HTTP + 明文端口，可再加 OTEL_EXPORTER_OTLP_INSECURE=true
ENV OTEL_EXPORTER_OTLP_ENDPOINT=http://opentelemetry-collector.observable.svc:4317 \
    OTEL_TRACES_EXPORTER=otlp \
    OTEL_SERVICE_NAME=my-go-service \
    OTEL_EXPORTER_OTLP_INSECURE=true

ENTRYPOINT ["/app/server"]