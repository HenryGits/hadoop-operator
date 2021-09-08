/**
 @author: ZHC
 @date: 2021-09-08 14:02:31
 @description:
**/
package tools

import (
	"bytes"
	"io"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/klog/v2"
	"strings"
)

// ParseYaml 解析yaml文件内容，返回unstructured object
func ParseYaml(content string) (objects []*runtime.Object, err error) {
	decoder := yaml.NewYAMLOrJSONDecoder(strings.NewReader(content), 4096)
	for {
		extension := runtime.RawExtension{}

		err := decoder.Decode(&extension)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				klog.ErrorS(err, "yaml parse error")
				return nil, err
			}
		}

		object, err := runtime.Decode(unstructured.UnstructuredJSONScheme, bytes.TrimSpace(extension.Raw))
		if err != nil {
			klog.ErrorS(err, "runtime decode error")
			return nil, err
		}
		objects = append(objects, &object)
	}
	return objects, nil
}
