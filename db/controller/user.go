package controller

import (
	"db/repository/db/dao"
	req "db/request"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func List(context *gin.Context) {
	list := dao.NewUserDao(context).UserList()
	context.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"data":    list,
	})
}

func Detail(context *gin.Context) {
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

func Add(context *gin.Context) {
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

func Update(context *gin.Context) {
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

func Delete(context *gin.Context) {
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
