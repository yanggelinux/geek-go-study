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
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const APP_NAME = "dbeye"
const SERVICE_NAME = "dbeye"

func setupSetting(config string) {
	if config == "" {
		panic("配置文件不能为空")
	}
	sett, err := setting.NewSetting(config)
	if err != nil {
		err = fmt.Errorf("setting set new configs error:%w", err)
		panic(err)
	}
	err = sett.ReadSection("APP", &global.APPSetting)
	if err != nil {
		err = fmt.Errorf("setting read app section error:%w", err)
		panic(err)
	}
	err = sett.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		err = fmt.Errorf("setting read server section error:%w", err)
		panic(err)
	}
	err = sett.ReadSection("MySQL", &global.MySQLSetting)
	if err != nil {
		err = fmt.Errorf("setting read mysql section error:%w", err)
		panic(err)
	}
	err = sett.ReadSection("JWT", &global.JWTSetting)
	if err != nil {
		err = fmt.Errorf("setting read jwt section error:%w", err)
		panic(err)
	}

}

func setupAny(config string) {

	//初始化全局配置信息
	setupSetting(config)
	//初始化Log配置
	log.SetupLog(SERVICE_NAME, APP_NAME)
	//初始化mysql
	mysql.SetupModel()

}

func main() {
	var (
		config string
		err    error
	)
	flag.StringVar(&config, "config", "", "请制定配置文件路径")
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
	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error {
		err := s.ListenAndServe()
		return err
	})

	g.Go(func() error {
		select {
		case <-ctx.Done():
			log.Logger.Info("errgroup exit...")
		}

		timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		log.Logger.Info("shutting down server...")
		return s.Shutdown(timeoutCtx)
	})

	g.Go(func() error {
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case sig := <-quit:
			return errors.Errorf("get os signal: %v", sig)
		}
	})
	log.Logger.Info("server start...")
	err = g.Wait()
	if err != nil {
		log.Logger.Error("http server exit...:", zap.Error(err))
	}
}
