package controller

import (
	"db/repository/db/dao"
	req "db/request"
	"db/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Login 登录
func Login(context *gin.Context) {
	var request req.LoginRequest
	if err := context.ShouldBind(&request); err != nil {
		context.JSON(http.StatusOK, gin.H{
			"message": "参数错误",
			"result":  err.Error(),
		})
		return
	}

	info, err := dao.NewUserDao(context).GetUserInfo(request)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"message": "用户不存在",
			"result":  err.Error(),
		})
		return
	}

	token, err := util.GenerateToken(int64(info.ID), info.Username)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"message": "系统错误，请稍后再试",
			"result":  err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"token":   token,
	})
}

// Register 注册
func Register(context *gin.Context) {
	var request req.RegisterRequest
	if err := context.ShouldBind(&request); err != nil {
		context.JSON(http.StatusOK, gin.H{
			"message": "参数错误",
			"result":  err.Error(),
		})
		return
	}

	err := dao.NewUserDao(context).CreateUser(request)

	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"message": "注册失败",
			"error":   err,
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "注册成功",
	})
	return
}
