package main

import (
	"fmt"
	"logger/util/logger"
	"os"
	"path"
	"time"
)

func main() {
	// 当前时间
	now := time.Now()
	// 文件根目录
	baseDir, _ := os.Getwd()
	// 日志路径
	filename := path.Join(baseDir, "runtime", "log", now.Format("200601"), now.Format("02")+".log")
	_, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Error(fmt.Sprintf("%v", err))
	}

	//logger.Log()
	logger.Trace("几乎任何东西")
	logger.Info("重要信息")
	logger.Warning("警告")
	logger.Error("错误")
}
