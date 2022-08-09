package log

import (
	"fmt"
	"geek/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var Logger *zap.Logger

type Field = zapcore.Field

func SetupLog(serviceName string, appName string) {
	filename := fmt.Sprintf("%s/%s/%s.log", global.APPSetting.LogPath, serviceName, appName)
	hook := lumberjack.Logger{
		Filename:   filename, // 日志文件路径
		MaxSize:    300,      // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 15,       // 日志文件最多保存多少个备份
		MaxAge:     15,       // 文件最多保存多少天
		Compress:   false,    // 是否压缩
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "date",
		LevelKey:       "level",
		NameKey:        "name",
		CallerKey:      "line",
		MessageKey:     "message",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,                            // 小写编码器
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000"), // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,                         //
		EncodeCaller:   zapcore.ShortCallerEncoder,                             // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 设置日志级别

	atomicLevel := zap.NewAtomicLevel()
	level := global.APPSetting.LogLevel
	if level == "DEBUG" {
		atomicLevel.SetLevel(zap.DebugLevel)
	} else if level == "INFO" {
		atomicLevel.SetLevel(zap.InfoLevel)
	} else if level == "WARN" {
		atomicLevel.SetLevel(zap.WarnLevel)
	} else {
		atomicLevel.SetLevel(zap.ErrorLevel)
	}

	if global.APPSetting.Development {
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),                                           // 编码器配置
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
			atomicLevel, // 日志级别
		)
		// 开启开发模式
		development := zap.Development()
		// 开启文件及行号
		caller := zap.AddCaller()
		// 设置初始化字段
		filed := zap.Fields(zap.String("serviceName", serviceName), zap.String("appName", appName))
		// 构造日志
		logger := zap.New(core, caller, development, filed)

		Logger = logger
	} else {
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig), // 编码器配置
			//zapcore.AddSync(os.Stdout)
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(&hook)), // 打印到控制台和文件
			atomicLevel, // 日志级别
		)
		// 开启文件及行号
		caller := zap.AddCaller()
		// 设置初始化字段
		filed := zap.Fields(zap.String("serviceName", serviceName), zap.String("appName", appName))
		// 构造日志
		logger := zap.New(core, filed, caller)
		Logger = logger
	}
}
