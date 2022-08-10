package main

import (
	"context"
	"flag"
	"fmt"
	"geek/global"
	sc "geek/internal/sync"
	"geek/pkg/log"
	"geek/setting"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func setupAny() {

	//初始化Log配置
	serviceName := "geek"
	appName := global.APPSetting.Name
	log.SetupLog(serviceName, appName)

}

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
	err = sett.ReadSection("Kafka", &global.KafkaSetting)
	if err != nil {
		err = fmt.Errorf("setting read kafka section error:%w", err)
		panic(err)
	}

}

func main() {
	var (
		config string
		cpuNum int
	)
	flag.StringVar(&config, "config", "", "请制定配置文件路径")
	flag.IntVar(&cpuNum, "cpuNum", 1, "请输入cpu核心数量，如果不输入默认4核心")
	flag.Parse()
	//先读取配置
	setupSetting(config)
	//再初始化
	setupAny()
	log.Logger.Info("主程序启动，程序环境", zap.String("config", config), zap.Int("cpuNum", cpuNum))
	exit := make(chan struct{}, cpuNum)
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(cpuNum)
	for i := 0; i < cpuNum; i++ {
		go func() {
			defer wg.Done()
			sc.SyncData(ctx, exit)
		}()
	}
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {

	case <-sigterm:
		log.Logger.Info("主程序terminating: via signal")
	case <-ctx.Done():
		log.Logger.Info("主程序terminating: context cancelled")
	}
	cancel()
	wg.Wait()
	for i := 0; i < cpuNum; i++ {
		<-exit
	}
	close(exit)
	log.Logger.Info("主程序退出")
}
