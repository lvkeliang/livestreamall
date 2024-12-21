package config

import (
	"github.com/spf13/viper"
	"log"
)

type AppConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type StreamConfig struct {
	PullBaseURL string `mapstructure:"pull_base_url"`
	PushBaseURL string `mapstructure:"push_base_url"`
}

type Config struct {
	App    AppConfig    `mapstructure:"app"`
	Stream StreamConfig `mapstructure:"stream"`
}

var App AppConfig
var Stream StreamConfig

func LoadConfig() {
	viper.SetConfigName("config.yaml") // 配置文件名
	viper.SetConfigType("yaml")        // 配置文件类型
	viper.AddConfigPath("./config")    // 配置文件路径（当前目录）

	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	// 将配置赋值到全局变量
	App = cfg.App
	Stream = cfg.Stream
}
