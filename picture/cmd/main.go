package main

import (
	// 标准库
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"picture/internal/service"
	"syscall"
	"time"
	// 三方库
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	// 内部库
	"picture/config"
)

func main() {
	// 初始化配置
	config.Init()

	engine := gin.Default()

	// 测试接口
	engine.GET("ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, "success")
	})

	// swagger 文档
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	upload := service.NewUpload()
	// 图片上传
	engine.POST("/upload", upload.Upload)

	server := &http.Server{
		Addr:           config.Conf.Server.Host,
		Handler:        engine,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		// 服务连接
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ERROR listen: %s\n", err)
		}
	}()

	// 优雅的关闭
	var exitCh = make(chan int, 1)
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
	log.Println("成功关闭", <-exitCh)
}
