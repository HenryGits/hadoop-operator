/**
 @author: ZHC
 @date: 2021-09-09 09:25:06
 @description:
**/
package tools

import (
	"k8s.io/klog/v2"
	"reflect"
)

// Has 字符串数组是否包含给定的值
func Has(array []string, value string) bool {
	for _, item := range array {
		if item == value {
			return true
		}
	}
	return false
}

// Remove 从字符串数组中移除指定元素
func Remove(source []string, value string) (destination []string) {
	for _, item := range source {
		if item != value {
			destination = append(destination, item)
		}
	}
	return destination
}

// Filter 数组元素过滤器（保留符合条件的元素）
func Filter(array interface{}, path string, value interface{}) (result []interface{}) {
	if obj := reflect.ValueOf(array); obj.Kind() == reflect.Slice {
		for i := 0; i < obj.Len(); i++ {
			data := Snipe(obj.Index(i).Interface(), path)
			dv := reflect.ValueOf(data)
			if dv.Kind() == reflect.ValueOf(value).Kind() {
				switch dv.Kind() {
				case reflect.String:
					if dv.String() == value.(string) {
						result = append(result, obj.Index(i).Interface())
					}
				case reflect.Bool:
					if dv.Bool() == value.(bool) {
						result = append(result, obj.Index(i).Interface())
					}
				case reflect.Int:
					if dv.Int() == value.(int64) {
						result = append(result, obj.Index(i).Interface())
					}
				default:
					klog.V(4).Infof("not supported type: %v, %v, %v", array, path, dv.Kind())
				}
			} else {
				klog.V(4).Infof("different type: %v, %v, %v, %v", array, path, data, value)
			}
		}
	} else {
		klog.Errorf("error type: %v", obj.Kind())
	}
	klog.V(4).Infof("filter result: %v, %v, %v, %v", array, path, value, result)
	return result
}

// AntiFilter 数组元素过滤器（保留不符合条件的元素）
func AntiFilter(array interface{}, path string, value interface{}) (result []interface{}) {
	if obj := reflect.ValueOf(array); obj.Kind() == reflect.Slice {
		for i := 0; i < obj.Len(); i++ {
			data := Snipe(obj.Index(i).Interface(), path)
			dv := reflect.ValueOf(data)
			if dv.Kind() == reflect.ValueOf(value).Kind() {
				switch dv.Kind() {
				case reflect.String:
					if dv.String() != value.(string) {
						result = append(result, obj.Index(i).Interface())
					}
				case reflect.Bool:
					if dv.Bool() != value.(bool) {
						result = append(result, obj.Index(i).Interface())
					}
				case reflect.Int:
					if dv.Int() != value.(int64) {
						result = append(result, obj.Index(i).Interface())
					}
				default:
					klog.V(4).Infof("filter not supported type: %v, %v, %v", array, path, dv.Kind())
				}
			} else {
				klog.V(4).Infof("filter different type: %v, %v, %v, %v", array, path, data, value)
			}
		}
	} else {
		klog.Errorf("filter error type: %v", obj.Kind())
	}
	klog.V(4).Infof("filter result: %v, %v, %v, %v", array, path, value, result)
	return result
}
