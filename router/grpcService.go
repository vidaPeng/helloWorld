package router

import (
	"context"
	"github.com/PixDevopsSre/helloWorld/pkg"
	"github.com/PixDevopsSre/helloWorld/proto"
	"go.opentelemetry.io/otel/attribute"
)

type GrpcService struct {
	proto.UnimplementedGreeterServer // ✅ 嵌入可选的默认实现
}

// SayHello 实现 GreeterServer 接口的 SayHello 方法
func (s *GrpcService) SayHello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloReply, error) {
	// 创建一个新的 span
	ctx, span := tracer.Start(ctx, "GrpcService.SayHello")
	defer span.End()

	// 记录 trace ID
	pkg.InfoTrace(ctx, "GrpcService.SayHello called with trace ID")

	// 模拟处理请求
	reply := &proto.HelloReply{
		Message: "Hello " + req.Name,
	}

	// 设置 span 的属性
	span.SetAttributes(attribute.String("response.message", reply.Message))

	return reply, nil
}
