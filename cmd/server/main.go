package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	userv1 "jyb-resource-mgr/api/user/v1"
	"jyb-resource-mgr/internal/handler"
	"jyb-resource-mgr/internal/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// service 只负责CRUD，不关心http和grpc的调用方式
	svc := service.NewUserService()

	// 启动grpc服务，用于处理gRPC请求
	go func() {
		// 监听9090端口
		lis, err := net.Listen("tcp", ":9090")
		if err != nil {
			log.Fatalf("gRPC listen failed: %v", err)
		}
		// 创建grpc服务器实例
		s := grpc.NewServer()
		// handler.NewGRPCServer(svc) 将把业务逻辑svc包装成gRPC handler对象
		// 将handler对象注册到grpc服务器实例s
		userv1.RegisterUserServiceServer(s, handler.NewGRPCServer(svc))
		// 开启反射，grpcurl才能发现服务
		reflection.Register(s)
		log.Println("[gRPC] listening on :9090")
		// 启动grpc服务器实例s
		if err := s.Serve(lis); err != nil {
			log.Fatalf("[gRPC] serve failed: %v", err)
		}
	}()

	// 启动http服务
	go func() {
		// 创建http服务器实例，用于处理HTTP请求
		httpServer := handler.NewHTTPServer(svc)
		// 创建http路由复用器，根据请求路径调用不同的handler
		mux := http.NewServeMux()
		// 注册http路由
		httpServer.RegisterRoutes(mux)
		log.Println("[HTTP] listening on :8080")
		// 启动http服务器实例mux，监听8080端口
		if err := http.ListenAndServe(":8080", mux); err != nil {
			log.Fatalf("[HTTP] serve failed: %v", err)
		}
	}()

	// 创建一个信号通道，用于接收系统信号
	quit := make(chan os.Signal, 1)
	// 收到SIGINT或SIGTERM信号时，将其发送到channel
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// 阻塞等待信号
	<-quit
	log.Println("shutting down...")
}
