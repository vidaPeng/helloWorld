# 第一阶段：构建 Go 应用
FROM golang:1.24 AS builder

WORKDIR /app

# 复制 go.mod 和 go.sum 先做依赖缓存（提升构建效率）
COPY go.mod ./
RUN go mod download

# 再复制项目代码
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /server .main.go

# 第二阶段：使用更小的基础镜像运行程序
FROM alpine:latest

# 可选：使用非 root 用户更安全（记得开放文件或端口权限）
# RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# 可选：创建配置文件目录（如果你的程序依赖 conf）
RUN mkdir -p /app/conf

# 从构建阶段复制可执行文件
COPY --from=builder /server .

# 可选：复制静态文件或配置文件（如果有）
# COPY ./conf /app/conf

# 可选：切换到非 root 用户运行（如上所设）
# USER appuser

# 启动入口
ENTRYPOINT ["/app/server"]