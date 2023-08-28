package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gin-grpc-etcd/client/rpc"
	"gin-grpc-etcd/idl/pb"
)

func Ping(ctx *gin.Context) {
	var req pb.Empty

	ping, err := rpc.UserPing(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    0,
			"message": "ping RPC服务调用错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": ping.Message,
	})
}

func Login(ctx *gin.Context) {
	var req pb.LoginRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    0,
			"message": "绑定参数错误",
		})
		return
	}

	login, err := rpc.UserLogin(ctx, &req)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    0,
			"message": "login RPC服务调用错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": login.Message,
	})
}
