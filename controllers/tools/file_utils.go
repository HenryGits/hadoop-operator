/**
 @author: ZHC
 @date: 2021-09-08 14:09:51
 @description:
**/
package tools

import (
	"k8s.io/klog/v2"
	"os"
	"path/filepath"
	"regexp"
)

// Find 在指定的目录下查找符合条件的文件
func Find(directory, pattern string) (files []string) {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		klog.Errorf("compile regex error: %v", err)
	}
	if err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err == nil && regex.MatchString(info.Name()) {
			files = append(files, path)
		}
		return nil
	}); err != nil {
		klog.Errorf("recursive file error: %v", err)
		return nil
	}
	return
}

// EnvVar 返回环境变量或者默认值
func EnvVar(name, defaultValue string) string {
	env, found := os.LookupEnv(name)
	if found {
		return env
	}
	return defaultValue
}
