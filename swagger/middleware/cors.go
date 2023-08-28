package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 请求方法
		method := context.Request.Method
		// 该字段是必须的。它的值要么是请求时Origin字段的值，要么是一个*，表示接受任意域名的请求
		context.Header("Access-Control-Allow-Origin", "*")
		// 如果浏览器请求包括Access-Control-Request-Headers字段，则Access-Control-Allow-Headers字段是必需的。它也是一个逗号分隔的字符串，表明服务器支持的所有头信息字段，不限于浏览器在"预检"中请求的字段。
		context.Header("Access-Control-Allow-Headers", "*")
		// 该字段必需，它的值是逗号分隔的一个字符串，表明服务器支持的所有跨域请求的方法。注意，返回的是所有支持的方法，而不单是浏览器请求的那个方法。这是为了避免多次"预检"请求。
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			context.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		context.Next()
	}
}
