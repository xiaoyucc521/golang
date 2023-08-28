package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"swagger/repository/db/dao"
	req "swagger/request"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

// List 用户列表
// @Summary      用户列表
// @Description  用户列表
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "token"
// @Router       /api/v1/user [get]
func (user *UserController) List(context *gin.Context) {
	list := dao.NewUserDao(context).UserList()
	context.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"data":    list,
	})
}

// Detail 用户详情
// @Summary      用户详情
// @Description  用户详情
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "token"
// @Param        id path int true "用户ID"
// @Router       /api/v1/user/{id} [get]
func (user *UserController) Detail(context *gin.Context) {
	id, _ := strconv.ParseInt(context.Param("id"), 10, 64)
	info, err := dao.NewUserDao(context).GetUserInfoById(id)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"message": "系统错误，请稍后再试",
			"error":   err.Error(),
		})
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"data":    info,
	})
}

// Add 新增用户
// @Summary      新增用户
// @Description  新增用户
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "token"
// @Param        username formData string true "用户名"
// @Param        password formData string true "密码"
// @Param        mobile formData string true "手机号"
// @Router       /api/v1/user/add [post]
func (user *UserController) Add(context *gin.Context) {
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
			"message": "创建失败",
			"error":   err,
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "创建成功",
	})
	return
}

// Update 修改用户
// @Summary      修改用户
// @Description  修改用户
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "token"
// @Param        id path int true "用户ID"
// @Param        username query string true "用户名"
// @Param        password query string true "密码"
// @Param        mobile query string true "手机号"
// @Router       /api/v1/user/{id} [put]
func (user *UserController) Update(context *gin.Context) {
	id, _ := strconv.ParseInt(context.Param("id"), 10, 64)

	var request req.RegisterRequest
	if err := context.ShouldBind(&request); err != nil {
		context.JSON(http.StatusOK, gin.H{
			"message": "参数错误",
			"result":  err.Error(),
		})
		return
	}

	err := dao.NewUserDao(context).UpdateUser(id, request)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"message": "修改失败",
			"error":   err,
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "修改成功",
	})
	return
}

// Delete 删除用户
// @Summary      删除用户
// @Description  删除用户
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id path int true "用户ID"
// @Router       /api/v1/user/{id} [delete]
func (user *UserController) Delete(context *gin.Context) {
	id, _ := strconv.ParseInt(context.Param("id"), 10, 64)
	err := dao.NewUserDao(context).DeleteUserById(id)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"message": "删除失败",
			"error":   err,
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "删除成功",
	})
	return
}
