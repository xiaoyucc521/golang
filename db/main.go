package main

import (
	"context"
	"db/config"
	"db/repository/db/dao"
	router2 "db/router"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var exitCh = make(chan int)

func main() {
	// 加载配置
	config.InitConfig()
	// 数据库初始化
	dao.InitDB()

	// 使用默认路由
	router := gin.Default()

	// 注册路由
	router = router2.NewRouter(router)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		// 服务连接
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	go gracefulExit(server)
	<-exitCh
	log.Println("成功关闭")
}

// gracefulExit 优雅的关闭
func gracefulExit(server *http.Server) {
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
}
