package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"
)

// getLogFile 获取主日志文件名
func getLogFile() (*os.File, error) {
	// 当前时间
	now := time.Now()
	// 文件根目录
	baseDir, _ := os.Getwd()
	// 日志路径
	filepath := path.Join(baseDir, "runtime", "log", now.Format("200601"))

	// 判断路径是否存在，不存在创建
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(filepath, 0755); err != nil {
			return nil, err
		}
	}

	// log 文件
	filename := path.Join(filepath, now.Format("02")+".log")
	// 写入文件
	src, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return src, nil
}

func Trace(format string) {
	logger := log.New(os.Stdout, "[TRACE] ", log.Ldate|log.Ltime|log.Lshortfile)
	_ = logger.Output(2, format)
}

func Info(format string) {
	logger := log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	_ = logger.Output(2, format)
}

func Warning(format string) {
	logger := log.New(os.Stdout, "[WAINING] ", log.Ldate|log.Ltime|log.Lshortfile)
	_ = logger.Output(2, format)
}

func Error(format string) {
	file, err := getLogFile()
	if err != nil {
		fmt.Println(err)
	}
	logger := log.New(io.MultiWriter(file, os.Stderr), "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	_ = logger.Output(2, format)
}

func Debug(format string) {
	file, err := getLogFile()
	if err != nil {
		fmt.Println(err)
	}
	logger := log.New(io.MultiWriter(file, os.Stderr), "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
	_ = logger.Output(2, format)
}
