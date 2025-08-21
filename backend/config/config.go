package config

import (
	"log"
	"main/controllers"
	"os"
	"os/signal"
	"syscall"

	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name         string `mapstructure:"name" yaml:"name"`
		Version      string `mapstructure:"version" yaml:"version"`
		Author       string `mapstructure:"author" yaml:"author"`
		BackendPort  string `mapstructure:"backend-port" yaml:"backend-port"`
		FrontendPort string `mapstructure:"frontend-port" yaml:"frontend-port"`
	} `mapstructure:"app" yaml:"app"`

	Database struct {
		DSN          string `mapstructure:"dsn" yaml:"dsn"`
		MaxIdleConns int    `mapstructure:"max-idle-conns" yaml:"max-idle-conns"`
		MaxOpenConns int    `mapstructure:"max-open-conns" yaml:"max-open-conns"`
	} `mapstructure:"database" yaml:"database"`

	Redis struct {
		Addr     string `mapstructure:"addr" yaml:"addr"`
		Password string `mapstructure:"password" yaml:"password"`
		DB       int    `mapstructure:"db" yaml:"db"`
		PoolSize int    `mapstructure:"poolsize" yaml:"pool-size"`
	} `mapstructure:"redis" yaml:"redis"`
}

var (
	AppConfig      *Config
	ValidRelations = make(map[string]bool)
)

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("无法读取yaml配置文件: %v", err)
	}

	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	InitDB()
	InitRedis()
	// 启动定时任务
	cronManager := controllers.StartCronJobs()

	// 设置优雅关闭
	setupGracefulShutdown(cronManager)
}

func setupGracefulShutdown(cronManager *cron.Cron) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("收到关闭信号，正在优雅关闭...")

		// 停止定时任务
		controllers.StopCronJobs(cronManager)

		log.Println("应用已关闭")
		os.Exit(0)
	}()
}
