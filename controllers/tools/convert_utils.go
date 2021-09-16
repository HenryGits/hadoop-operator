/**
 @author: ZHC
 @date: 2021-09-09 09:23:59
 @description:
**/
package tools

import (
	"k8s.io/klog/v2"
	"strconv"
)

func ParseInt(source string, base, bitSize int) int64 {
	if result, err := strconv.ParseInt(source, base, bitSize); err != nil {
		klog.ErrorS(err, "parse int error", "source", source, "base", base, "bitSize", bitSize)
		return 0
	} else {
		return result
	}
}

// CloneMap clones a map and applies patches without overwrites
func CloneMap(ms ...map[string]string) map[string]string {
	a := map[string]string{}
	for _, p := range ms {
		for k, v := range p {
			if _, ok := a[k]; ok {
				continue
			}
			a[k] = v
		}
	}
	return a
}
