package main

import (
	"fmt"
	"main/config"
	"main/routers"
)

func main() {
	config.InitConfig()
	fmt.Printf("正在加载配置文件：%+v\n", config.AppConfig)
	router := routers.SetupRouter()
	port := config.AppConfig.App.BackendPort
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
	fmt.Println("已经完成加载")

}
