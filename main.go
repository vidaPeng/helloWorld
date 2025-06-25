package main

import (
	"github.com/PixDevopsSre/helloWorld/pkg"
	"github.com/PixDevopsSre/helloWorld/router"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// 初始化日志组件
	pkg.InitLogger()
	defer pkg.Sync()

	pkg.InitTracer()

	r := gin.New()
	router.SetupRoutes(r)
	if err := r.Run(":80"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
