package config

import (
	// 标准库
	"os"
	// 三方库
	"github.com/spf13/viper"
)

var Conf *Config

type Config struct {
	Server *Server `yaml:"server"`
}

type Server struct {
	Name    string `yaml:"name"`
	Host    string `yaml:"host"`
	Version string `yaml:"version"`
	Weight  int    `yaml:"weight"`
}

func Init() {
	dir, _ := os.Getwd()
	// 文件名称
	viper.SetConfigName("config")
	// 文件后缀
	viper.SetConfigType("yml")
	viper.AddConfigPath(dir + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&Conf)
	if err != nil {
		panic(err)
	}
}
