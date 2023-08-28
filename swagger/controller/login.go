package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"swagger/repository/db/dao"
	req "swagger/request"
	"swagger/util"
)

type LoginController struct{}

func NewLoginController() *LoginController {
	return &LoginController{}
}

// Login 登录
// @Summary 	登录
// @Tags 		auth
// @Produce 	json
// @Param 		username formData string true "用户名"
// @Param 		password formData string true "密码"
// @Router 		/api/v1/login [post]
func (login *LoginController) Login(context *gin.Context) {
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
// @Summary      注册
// @Description  注册
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        username formData string true "用户名"
// @Param        password formData string true "密码"
// @Param        mobile formData string true "手机号"
// @Router       /api/v1/register [post]
func (login *LoginController) Register(context *gin.Context) {
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
