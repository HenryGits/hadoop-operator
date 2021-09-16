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
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
)

var (
	log *zap.Logger

	nnDir = flag.String("nnDir", "/usr/local/hadoop/nn", "NameNode Directory")
)

func main() {
	log = Zap()
	ctx := context.Background()
	log.Info("NameNode dirPath: ", zap.String("msg", *nnDir))

	if *nnDir != "" {
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
