package main

import (
	"context"
	"flag"
	"fmt"
	"geek/global"
	"geek/internal/dbeye/router"
	"geek/internal/dbeye/stroe/mysql"
	"geek/pkg/log"
	"geek/setting"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const APP_NAME = "dbeye"
const SERVICE_NAME = "dbeye"

func setupSetting(config string) {
	if config == "" {
		log.Logger.Panic("配置文件不能为空")
	}
	sett, err := setting.NewSetting(config)
	if err != nil {
		log.Logger.Panic("setting set new configs error", zap.Error(err))
	}
	err = sett.ReadSection("APP", &global.APPSetting)
	if err != nil {
		log.Logger.Panic("setting read app section error", zap.Error(err))
	}
	err = sett.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		log.Logger.Panic("setting read kafka section error", zap.Error(err))
	}
	err = sett.ReadSection("MySQL", &global.MySQLSetting)
	if err != nil {
		log.Logger.Panic("setting read bulk http section error", zap.Error(err))
	}
	err = sett.ReadSection("JWT", &global.JWTSetting)

}

func setupAny(config string) {

	//初始化全局配置信息
	setupSetting(config)
	//初始化Log配置
	log.SetupLog(SERVICE_NAME, APP_NAME)
	mysql.SetupModel()

}

func main() {
	var (
		config string
		err    error
	)
	flag.StringVar(&config, "configs", "", "请制定配置文件路径")
	flag.Parse()
	//初始化配置
	setupAny(config)

	runMsg := fmt.Sprintf("config is %s", config)
	log.Logger.Info(runMsg)

	//gin
	routerHandler := router.NewRouter()
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", global.ServerSetting.HTTPPort),
		Handler:        routerHandler,
		ReadTimeout:    time.Duration(global.ServerSetting.ReadTimeout),
		WriteTimeout:   time.Duration(global.ServerSetting.WriteTimeout),
		MaxHeaderBytes: 1 << 20,
	}
	runMsg = fmt.Sprintf("running app success port %d,config is %s", global.ServerSetting.HTTPPort, config)
	log.Logger.Info(runMsg)

	//优雅起停服务
	go func() {
		// 将服务在 goroutine 中启动
		err = s.ListenAndServe()
		if err != nil {
			log.Logger.Error("run app error exited", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit // 阻塞等待接收 channel 数据
	log.Logger.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 5s 缓冲时间处理已有请求
	defer cancel()
	if err := s.Shutdown(ctx); err != nil { // 调用 net/http 包提供的优雅关闭函数：Shutdown
		log.Logger.Error("Server Shutdown:", zap.Error(err))
	}
	log.Logger.Info("Server exiting")
}
