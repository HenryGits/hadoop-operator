/**
 @author: ZHC
 @date: 2021-09-08 14:05:34
 @description:
**/
package tools

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"k8s.io/klog/v2"
	"strings"
	"text/template"
)

type Parser struct {
	Directory string
	Pattern   string
}

// ParseTemplate 解析文件模版
func (p *Parser) ParseTemplate(name string, parameters interface{}) (string, error) {
	templates, err := template.New("GoToolsFileTextTemplate").Funcs(sprig.TxtFuncMap()).Funcs(p.buildFunctionMap()).ParseFiles(Find(p.Directory, p.Pattern)...)
	if err != nil {
		klog.ErrorS(err, "Failed to parse template", "name", name)
		return "", err
	}
	buffer := new(bytes.Buffer)
	if err := templates.ExecuteTemplate(buffer, name, parameters); err != nil {
		klog.ErrorS(err, "Failed to execute template", "name", name)
		return "", err
	}
	klog.V(8).Infoln(buffer.String())
	return buffer.String(), nil
}

// ParseString 解析字符串模版
func (p *Parser) ParseString(tmpl string, parameters interface{}) (string, error) {
	templates, err := template.New("GoToolsStringTextTemplate").Funcs(sprig.TxtFuncMap()).Funcs(p.buildFunctionMap()).Parse(tmpl)
	if err != nil {
		klog.ErrorS(err, "Failed to parse template", "template", tmpl)
		return "", err
	}
	buffer := new(bytes.Buffer)
	if err := templates.Execute(buffer, parameters); err != nil {
		klog.ErrorS(err, "Failed to execute template", "template", tmpl)
		return "", err
	}
	klog.V(8).Infoln(buffer.String())
	return buffer.String(), nil
}

// Include 解析模版，返回解析结果
func (p *Parser) Include(name string, parameters interface{}) (result string) {
	result, err := p.ParseTemplate(name, parameters)
	if err != nil {
		klog.ErrorS(err, "template include error", "name", name)
		return ""
	}
	return result
}

// NestedField returns a reference to a nested field.
// Returns false if value is not found and an error if unable
// to traverse obj.
//
// Note: fields passed to this function are treated as keys within the passed
// object; no array/slice syntax is supported.
// Reference: k8s.io/apimachinery@v0.20.2/pkg/apis/meta/v1/unstructured/helpers.go:53
func NestedField(obj map[string]interface{}, fields ...string) (interface{}, bool, error) {
	var val interface{} = obj

	for i, field := range fields {
		if val == nil {
			return nil, false, nil
		}
		if m, ok := val.(map[string]interface{}); ok {
			val, ok = m[field]
			if !ok {
				return nil, false, nil
			}
		} else {
			return nil, false, fmt.Errorf("%v accessor error: %v is of the type %T, expected map[string]interface{}", jsonPath(fields[:i+1]), val, val)
		}
	}
	return val, true, nil
}

func jsonPath(fields []string) string {
	return "." + strings.Join(fields, ".")
}

// 初始化需渲染的默认模板路径
var defaultParser = &Parser{EnvVar("GT_TEMPLATE_PATH", "/etc/operator/templates"), "\\.gotmpl$"}

// ParseTemplate 解析文件模版函数
func ParseTemplate(name string, parameters interface{}) (string, error) {
	return defaultParser.ParseTemplate(name, parameters)
}

// ParseString 解析字符串模版函数
func ParseString(tmpl string, parameters interface{}) (string, error) {
	return defaultParser.ParseString(tmpl, parameters)
}
