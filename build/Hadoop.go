/**
 @author: ZHC
 @date: 2021-09-16 17:12:22
 @description: Hadoop服务启动类
**/
package main

import (
	"context"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

var (
	log   *zap.Logger
	level zapcore.Level

	nnDir     = flag.String("nnDir", "/usr/local/hadoop/nn", "NameNode Directory")
	nnService = flag.Bool("NameNode", false, "是否启动NameNode服务")
	dnService = flag.Bool("DataNode", false, "是否启动DataNode服务")
)

func main() {
	log = Zap()
	ctx := context.Background()
	log.Info("NameNode dirPath: ", zap.String("msg", *nnDir))

	if *nnDir != "" && *nnService {
		path, _ := filepath.Abs(*nnDir)
		// 判断NameNode是否有效
		files, err := ioutil.ReadDir(path)
		if err != nil {
			log.Error("NameNode路径不存在!", zap.Error(err))
			os.Exit(1)
		}
		// 判断NameNode是否已被初始化过
		if files == nil {
			log.Info("Formatting NameNode name directory: ", zap.String("msg", *nnDir))
			err := exec.CommandContext(ctx, "$HADOOP_HOME/bin/hdfs namenode -format").Start()
			if err != nil {
				log.Error("NameNode初始化失败!", zap.Error(err))
				os.Exit(1)
			}
		}
		err = exec.CommandContext(ctx, "$HADOOP_HOME/bin/hdfs --daemon start namenode").Start()
		if err != nil {
			log.Error("NameNode启动失败!", zap.Error(err))
			os.Exit(1)
		}
	}

	if dnService != nil {
		err := exec.CommandContext(ctx, "$HADOOP_HOME/bin/hdfs --daemon start datanode").Start()
		if err != nil {
			log.Error("DataNode启动失败!", zap.Error(err))
			os.Exit(1)
		}
	}

	fmt.Println(`
 ██      ██          ██  ██            ██      ██                ██                          
░██     ░██         ░██ ░██           ░██     ░██               ░██                   ██████ 
░██     ░██  █████  ░██ ░██  ██████   ░██     ░██  ██████       ░██  ██████   ██████ ░██░░░██
░██████████ ██░░░██ ░██ ░██ ██░░░░██  ░██████████ ░░░░░░██   ██████ ██░░░░██ ██░░░░██░██  ░██
░██░░░░░░██░███████ ░██ ░██░██   ░██  ░██░░░░░░██  ███████  ██░░░██░██   ░██░██   ░██░██████ 
░██     ░██░██░░░░  ░██ ░██░██   ░██  ░██     ░██ ██░░░░██ ░██  ░██░██   ░██░██   ░██░██░░░  
░██     ░██░░██████ ███ ███░░██████   ░██     ░██░░████████░░██████░░██████ ░░██████ ░██     
░░      ░░  ░░░░░░ ░░░ ░░░  ░░░░░░    ░░      ░░  ░░░░░░░░  ░░░░░░  ░░░░░░   ░░░░░░  ░░
	`)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		fmt.Println("terminating: context cancelled")
	case <-sigterm:
		fmt.Println("terminating: via signal")
	}
}

// Zap 初始日志zap
func Zap() (logger *zap.Logger) {
	level = zap.InfoLevel
	logger = zap.New(getEncoderCore())
	logger = logger.WithOptions(zap.AddCaller())
	return logger
}

// getEncoderConfig 获取zapcore.EncoderConfig
func getEncoderConfig() (config zapcore.EncoderConfig) {
	config = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	config.EncodeLevel = zapcore.LowercaseLevelEncoder
	return config
}

// getEncoder 获取zapcore.Encoder
func getEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(getEncoderConfig())
}

// getEncoderCore 获取Encoder的zapcore.Core
func getEncoderCore() (core zapcore.Core) {
	writer := zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout))
	return zapcore.NewCore(getEncoder(), writer, level)
}

// 自定义日志输出时间格式
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("[Hadoop] " + "2006-01-02 15:04:05.000"))
}
