package main

import (
	"github.com/PixDevopsSre/helloWorld/pkg"
	"github.com/PixDevopsSre/helloWorld/proto"
	"github.com/PixDevopsSre/helloWorld/router"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof" // 关键：匿名导入pprof包
)

func main() {
	// 初始化日志组件
	pkg.InitLogger()
	defer pkg.Sync()

	// 初始化 OpenTelemetry 跟踪器
	pkg.InitTracer()

	// 1. 单独为 pprof 启动一个 HTTP 服务 (推荐)
	go func() {
		log.Println("Pprof server starting on localhost:6060...")
		// ListenAndServe 会使用 DefaultServeMux，pprof 路由已注册于此
		err := http.ListenAndServe("localhost:6060", nil)
		if err != nil {
			log.Printf("Pprof server failed: %v", err)
		}
	}()

	// 启用 Grpc 服务器
	go func() {
		listen, err := net.Listen("tcp", ":9000")
		if err != nil {
			panic(err)
		}
		service := grpc.NewServer(
			grpc.StatsHandler(otelgrpc.NewServerHandler()))

		proto.RegisterGreeterServer(service, &router.GrpcService{})

		err = service.Serve(listen)
		if err != nil {
			panic(err)
		}
	}()

	r := gin.New()
	r.Use(otelgin.Middleware("http-helloPeng"))
	router.SetupRoutes(r)
	if err := r.Run(":80"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
