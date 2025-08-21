package main

import (
	"fmt"
	"main/config"
	"main/routers"
)

func main() {
	config.InitConfig()
	fmt.Printf("正在加载配置文件：%v\n")
	router := routers.SetupRouter()
	port := config.AppConfig.App.BackendPort
	if port == "" {
		port = "8080"
	}
	router.Run("0.0.0.0:" + port)
	fmt.Printf("服务器已启动，监听端口：%s\n", port)

}
