package router

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"net/http"
	"swagger/controller"
	"swagger/middleware"
)

func NewRouter(router *gin.Engine) *gin.Engine {

	router.Use(middleware.Cors())
	// swagger 文档
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// 注册路由
	v1 := router.Group("/api/v1")
	{
		// 测试接口
		v1.GET("ping", func(context *gin.Context) {
			context.JSON(http.StatusOK, "success")
		})

		NewLoginController := controller.NewLoginController()
		v1.POST("login", NewLoginController.Login)
		v1.POST("register", NewLoginController.Register)

		// 实例化 UserController 实体
		UserController := controller.NewUserController()
		// 需要登录验证
		authed := v1.Group("/user")
		// 注册中间件
		authed.Use(middleware.JWT())
		{
			authed.GET("/", UserController.List)
			authed.GET("/:id", UserController.Detail)
			authed.POST("/add", UserController.Add)
			authed.PUT("/:id", UserController.Update)
			authed.DELETE("/:id", UserController.Delete)
		}
	}

	return router
}
