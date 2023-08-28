package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"gin-grpc-etcd/client/controller"
	"gin-grpc-etcd/client/rpc"
	"gin-grpc-etcd/config"
)

func main() {
	config.Init()
	rpc.Init()

	// 实例化 gin 框架
	engine := gin.Default()

	engine.GET("ping", controller.Ping)

	engine.POST("login", controller.Login)

	server := &http.Server{
		Addr:    ":8080",
		Handler: engine,
	}

	go func() {
		// 服务连接
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 监控系统型号，等待 ctrl + c 系统信号通知关闭
	exitCh := make(chan int, 1)
	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, os.Kill)

		sig := <-signalChan
		log.Printf("catch signal, %+v", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second) // 4秒后退出
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}
		log.Printf("server exiting")
		close(exitCh)
	}()
	log.Println(fmt.Sprintf("exit %v", <-exitCh))
}
