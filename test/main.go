package main

import (
	"github.com/PixDevopsSre/helloWorld/handle"
	"github.com/PixDevopsSre/helloWorld/pkg"
)

func init() {
	pkg.InitLogger()
}

func main() {
	handle.S3Handle()
}
