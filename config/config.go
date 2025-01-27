package config

import (
	"datarp/logger"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type Config struct {
	MySQL struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
		Database string `json:"database"`
	} `json:"mysql"`
	Server struct {
		Port int `json:"port"`
	} `json:"server"`
}

var (
	config *Config
	once   sync.Once
)

func LoadConfig() error {
	var err error
	once.Do(func() {
		config = &Config{}
		file, err := os.Open("config.json")
		if err != nil {
			return
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		err = decoder.Decode(config)
		if err == nil {
			logger.Info(fmt.Sprintf("成功加载配置文件，数据库地址: %s:%d", config.MySQL.Host, config.MySQL.Port))
		}
	})
	return err
}

func GetConfig() *Config {
	return config
} 