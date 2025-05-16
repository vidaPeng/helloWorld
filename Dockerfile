# ---------- build stage ----------
FROM golang:1.24 AS builder
WORKDIR /app

# 1. 先缓存依赖
COPY go.mod go.sum ./
RUN go mod download

# 2. 再复制其余代码
COPY . .

# 3. 编译
#   - CGO_DISABLED=0 生成静态二进制，适合扔到 distroless/alpine
#   - 如果 main.go 就在根目录，可直接 `go build -o /server .`
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /server ./main.go

# ---------- runtime stage ----------
FROM alpine:latest
WORKDIR /app

# 可选：更安全的非 root 运行
# RUN addgroup -S app && adduser -S app -G app
# USER app

COPY --from=builder /server /app/server
# 如果需要静态文件 / 配置文件，同样 COPY 过来
# COPY conf/ /app/conf/

ENTRYPOINT ["/app/server"]