/**
 @author: ZHC
 @date: 2021-09-16 17:12:22
 @description: Hadoop服务启动类
**/
package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

var (
	log   *zap.Logger
	level zapcore.Level

	nnDir       = flag.String("NameNodeDir", "/usr/local/hadoop/nn", "NameNode Directory")
	nnService   = flag.Bool("NameNode", false, "是否启动NameNode服务")
	dnService   = flag.Bool("DataNode", false, "是否启动DataNode服务")
	journalNode = flag.Bool("JournalNode", false, "是否启动JournalNode服务")
	rmService   = flag.Bool("ResourceManager", false, "是否启动ResourceManager服务")
	nmService   = flag.Bool("NodeManager", false, "是否启动NodeManager服务")
	hsService   = flag.Bool("HistoryServer", false, "是否启动HistoryServer服务")
)

func main() {
	log = Zap()
	ctx := context.Background()
	// Parse command line into the defined flags
	flag.Parse()

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

	if *nnService {
		path, _ := filepath.Abs(*nnDir)
		// 判断NameNode是否有效
		files, err := ioutil.ReadDir(path)
		if err != nil {
			log.Error("NameNode路径不存在!", zap.Error(err))
			os.Exit(1)
		}

		// 判断NameNode是否已被初始化过
		if len(files) <= 0 {
			cmd := exec.CommandContext(ctx, os.ExpandEnv("$HADOOP_HOME")+"/bin/hdfs", "namenode", "-format")
			_, err := cmd.CombinedOutput()
			if err != nil {
				log.Error("NameNode初始化失败!", zap.Error(err))
				os.Exit(1)
			}
			log.Info(fmt.Sprintf("=== NameNode Init Success === \n"))
		}
		err = CommandContext(ctx, os.ExpandEnv("$HADOOP_HOME")+"/bin/hdfs", "namenode")
		if err != nil {
			log.Error("NameNode启动失败!", zap.Error(err))
			os.Exit(1)
		}
	}

	if *dnService {
		err := CommandContext(ctx, os.ExpandEnv("$HADOOP_HOME")+"/bin/hdfs", "datanode")
		if err != nil {
			log.Error("DataNode启动失败!", zap.Error(err))
			os.Exit(1)
		}
	}

	if *journalNode {
		err := CommandContext(ctx, os.ExpandEnv("$HADOOP_HOME")+"/bin/hdfs", "journalnode")
		if err != nil {
			log.Error("JournalNode启动失败!", zap.Error(err))
			os.Exit(1)
		}
	}

	if *rmService {
		err := CommandContext(ctx, os.ExpandEnv("$HADOOP_HOME")+"/bin/yarn", "resourcemanager")
		if err != nil {
			log.Error("ResourceManager启动失败!", zap.Error(err))
			os.Exit(1)
		}
	}

	if *nmService {
		err := CommandContext(ctx, os.ExpandEnv("$HADOOP_HOME")+"/bin/yarn", "nodemanager")
		if err != nil {
			log.Error("NodeManager启动失败!", zap.Error(err))
			os.Exit(1)
		}
	}

	if *hsService {
		err := CommandContext(ctx, os.ExpandEnv("$HADOOP_HOME")+"/bin/mapred", "historyserver")
		if err != nil {
			log.Error("HistoryServer启动失败!", zap.Error(err))
			os.Exit(1)
		}
	}

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		fmt.Println("Bye: context cancelled")
	case <-sigterm:
		fmt.Println("Bye: signal cancelled")
	}
}

// CommandContext 执行shell实时输出日志
func CommandContext(ctx context.Context, name string, cmd ...string) error {
	c := exec.CommandContext(ctx, name, cmd...)
	stdout, err := c.StderrPipe()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		reader := bufio.NewReader(stdout)
		for {
			readString, err := reader.ReadString('\n')
			if err != nil || err == io.EOF {
				panic(err)
			}
			fmt.Print(readString)
		}
	}(&wg)
	err = c.Start()
	wg.Wait()
	return err
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
