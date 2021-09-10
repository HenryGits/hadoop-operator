/**
 @author: ZHC
 @date: 2021-09-09 09:23:10
 @description:
**/
package tools

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"k8s.io/klog/v2"
)

// HexEncode 字符串转Hash
func HexEncode(source string) string {
	return hex.EncodeToString([]byte(source))
}

// HexDecode Hash转字符串
func HexDecode(source string) string {
	if result, err := hex.DecodeString(source); err != nil {
		klog.ErrorS(err, "hex decode error", "source", source)
		return ""
	} else {
		return string(result)
	}
}

// MD5 计算md5值
func MD5(source string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(source)))
}
