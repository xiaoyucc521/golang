package router

import (
	"db/controller"
	"db/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewRouter(router *gin.Engine) *gin.Engine {

	// 注册路由
	v1 := router.Group("/api/v1")
	{
		// 测试接口
		v1.GET("ping", func(context *gin.Context) {
			context.JSON(http.StatusOK, "success")
		})

		v1.POST("register", controller.Register)
		v1.POST("login", controller.Login)

		// 需要登录验证
		authed := v1.Group("/user")
		// 注册中间件
		authed.Use(middleware.JWT())
		{
			authed.GET("/", controller.List)
			authed.GET("/:id", controller.Detail)
			authed.POST("/add", controller.Add)
			authed.PUT("/:id", controller.Update)
			authed.DELETE("/:id", controller.Delete)
		}
	}

	return router
}
