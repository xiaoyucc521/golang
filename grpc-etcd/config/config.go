package config

import (
	// 标准库
	"os"

	"github.com/spf13/viper"
)

var Conf *Config

type Config struct {
	Project  *Project            `yaml:"project"`
	Server   *Server             `yaml:"server"`
	Etcd     *Etcd               `yaml:"etcd"`
	Services map[string]*Service `yaml:"services"`
}

// Project 项目配置
type Project struct {
	Name string `yaml:"name"`
}

// Server 服务配置
type Server struct {
	Name    string `yaml:"name"`
	Host    string `yaml:"host"`
	Version string `yaml:"version"`
	Weight  int    `yaml:"weight"`
}

// Etcd 注册中心配置
type Etcd struct {
	Address []string `yaml:"address"`
}

type Service struct {
	Name        string `yaml:"name"`
	LoadBalance bool   `yaml:"loadBalance"`
}

// InitConfig 初始化配置
func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&Conf)
	if err != nil {
		panic(err)
	}
}
