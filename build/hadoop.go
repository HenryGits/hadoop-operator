/**
 @author: ZHC
 @date: 2021-09-16 17:12:22
 @description: Hadoop服务启动类
**/
package main

import "os"

func main() {

	nnDir := os.Getenv("HDFS_NAME_NODE_DIR")
	if nnDir != "" {
		// 判断路径不为空

	}

}
