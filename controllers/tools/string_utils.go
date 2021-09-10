/**
 @author: ZHC
 @date: 2021-09-07 16:47:40
 @description:
**/
package tools

import (
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"strings"
)

// Helper functions to check and remove string from a slice of strings.
func ContainsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// Name 生成随机名称
func Name() string {
	return strings.ReplaceAll(fmt.Sprintf("%s-%s-%s", strings.ToLower(randomdata.Alphanumeric(6)), strings.ToLower(randomdata.SillyName()), strings.ToLower(randomdata.Country(randomdata.TwoCharCountry))), " ", "-")
}
