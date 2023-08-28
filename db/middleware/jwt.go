package middleware

import (
	"db/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func JWT() gin.HandlerFunc {
	return func(context *gin.Context) {
		token := context.GetHeader("Authorization")
		if token == "" {
			context.JSON(http.StatusOK, gin.H{
				"status":  0,
				"message": "请登录",
			})
			context.Abort()
			return
		}
		claims, err := util.ParseToken(token)
		if err != nil {
			context.JSON(http.StatusOK, gin.H{
				"status":  0,
				"message": "token错误",
			})
		}
		context.Set("userId", claims.UserID)
		context.Set("username", claims.Username)
		// 处理请求
		context.Next()
	}
}
