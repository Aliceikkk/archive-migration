package main

import (
	"datarp/config"
	"datarp/database"
	"datarp/handlers"
	"datarp/logger"
	"fmt"
	"net/http"
)

func main() {
	// 初始化日志
	logger.InitLogger()
	defer logger.CloseLogger()

	// 加载配置
	if err := config.LoadConfig(); err != nil {
		logger.Error("加载配置失败: " + err.Error())
		return
	}

	// 初始化数据库连接
	if err := database.InitDB(); err != nil {
		logger.Error("初始化数据库失败: " + err.Error())
		return
	}
	conf := config.GetConfig()
	logger.Info(fmt.Sprintf("成功连接到数据库 %s:%d/%s", conf.MySQL.Host, conf.MySQL.Port, conf.MySQL.Database))

	// 设置路由
	http.Handle("/", http.FileServer(http.Dir("http")))
	http.HandleFunc("/migrate", handlers.HandleMigrate)

	// 启动服务器
	addr := fmt.Sprintf(":%d", conf.Server.Port)
	logger.Info("服务器启动在" + addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		logger.Error("服务器启动失败: " + err.Error())
	}
}
