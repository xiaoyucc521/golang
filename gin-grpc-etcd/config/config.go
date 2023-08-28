package config

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

var Conf *Config

type Config struct {
	Project  *Project             `yaml:"project"`
	Server   *Server              `yaml:"server"`
	Etcd     *Etcd                `yaml:"etcd"`
	Services map[string]*Services `yaml:"services"`
}

type Project struct {
	Name string `yaml:"name"`
}

type Server struct {
	Name    string `yaml:"name"`
	Host    string `yaml:"host"`
	Version string `yaml:"version"`
	Weight  int    `yaml:"weight"`
}

type Etcd struct {
	Addr []string `yaml:"addr"`
}

type Services struct {
	Name        string `yaml:"name"`
	LoadBalance bool   `yaml:"loadBalance"`
}

// Init 初始化配置
func Init() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/config")
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			panic("找不到配置文件")
		}
	}
	if err := viper.Unmarshal(&Conf); err != nil {
		panic(err)
	}
}
