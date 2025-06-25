package main

import (
	"github.com/PixDevopsSre/helloWorld/pkg"
	"github.com/PixDevopsSre/helloWorld/router"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"log"
)

func main() {
	// 初始化日志组件
	pkg.InitLogger()
	defer pkg.Sync()

	// 初始化 OpenTelemetry 跟踪器
	pkg.InitTracer()

	r := gin.New()
	r.Use(otelgin.Middleware("http-helloPeng"))
	router.SetupRoutes(r)
	if err := r.Run(":80"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
