/**
 @author: ZHC
 @date: 2021-09-09 09:18:46
 @description:
**/
package tools

import (
	"encoding/json"
	"fmt"
	"github.com/HenryGits/hadoop-operator/apis/hadoop/v1"
	"github.com/fatih/structs"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"path"
	"reflect"
	"strings"
	"text/template"
)

func (p *Parser) buildFunctionMap() template.FuncMap {
	return template.FuncMap{
		"toYaml":                     MustToYaml,
		"toJson":                     ToJson,
		"snipe":                      Snipe,
		"hitch":                      Hitch,
		"filter":                     Filter,
		"has":                        Has,
		"antiFilter":                 AntiFilter,
		"include":                    p.Include,
		"generateName":               Name,
		"hexenc":                     HexEncode,
		"hexdec":                     HexDecode,
		"md5":                        MD5,
		"parseInt":                   ParseInt,
		"parseInitial":               ParseInitial,
		"parseEnv":                   ParseEnv,
		"parseHandler":               ParseHandler,
		"parseLifecycle":             ParseLifecycle,
		"parseProbe":                 ParseProbe,
		"parseService":               ParseService,
		"parsePodPort":               ParsePodPort,
		"parseNodeSelector":          ParseNodeSelector,
		"parseExporter":              ParseExporter,
		"parseHostAliases":           ParseHostAliases,
		"parseResources":             ParseResources,
		"parsePersistentVolumeClaim": ParsePersistentVolumeClaim,
		"parseConfigMap":             ParseConfigMap,
		"parseVolumes":               ParseVolumes,
		"parseVolumeMounts":          ParseVolumeMounts,
		"parseDestinationRule":       ParseDestinationRule,
		"parseArguments":             ParseArguments,
		"parseUpdateStrategy":        ParseUpdateStrategy,
		"parsePodManagementPolicy":   ParsePodManagementPolicy,
		"parseKind":                  ParseKind,
		"parseCommand":               ParseCommand,
		"parseImage":                 ParseImage,
		"parseAutoscaler":            ParseAutoscaler,
		"parseReplicas":              ParseReplicas,
		"compassHelper":              KubernetesConfHelper,
	}
}

// MustToYaml 将结构体转化为Yaml格式的字符串
func MustToYaml(object interface{}) string {
	bytess, err := yaml.Marshal(object)
	if err != nil {
		klog.ErrorS(err, "Failed to marshal struct")
		return ""
	}
	klog.V(8).Infoln(string(bytess), "method", "toYaml")
	return string(bytess)
}

// ToJson 将结构体转化为JSON格式的字符串
func ToJson(object interface{}) string {
	bytess, err := json.Marshal(object)
	if err != nil {
		klog.ErrorS(err, "Failed to marshal struct")
		return ""
	}
	klog.V(8).Infoln(string(bytess), "method", "toJson")
	return string(bytess)
}

// Snipe 根据路径获取结构体字段值
func Snipe(object interface{}, path string) interface{} {
	data, found, err := NestedField(structs.Map(object), strings.Split(path, ".")...)
	if err != nil {
		klog.ErrorS(err, "snipe error")
		return nil
	}
	if found {
		return data
	}
	return nil
}

// Hitch 根据路径设置结构体字段值
func Hitch(object interface{}, path string, value interface{}) (result map[string]interface{}) {
	if object == nil {
		result = make(map[string]interface{})
	} else {
		result = structs.Map(object)
	}
	if i := strings.Index(path, "."); i == -1 {
		result[path] = value
	} else {
		result[path[:i]] = Hitch(nil, path[i+1:], value)
	}
	return result
}

// ParseInitial 解析initial
func ParseInitial(parameters interface{}) string {
	tpl := `- name: initial
  image: busybox:1.31.1
  command:
    - /bin/sh
    - -c
    - |
      {{ . }}
  securityContext:
    runAsUser: 0`

	res, _ := defaultParser.ParseString(tpl, parameters)
	return res
}

// ParseEnv 解析环境变量
func ParseEnv(parameters interface{}) string {
	tpl := `{{- with . -}}
env:
{{- range . }}
  - name: {{ .Name | quote }}
    value: {{ .Value | quote }}
{{- end }}
{{- end }}`

	res, _ := defaultParser.ParseString(tpl, parameters)
	return res
}

// ParsePersistentVolumeClaim 解析存储卷
func ParsePersistentVolumeClaim(volumes []*v1.Volume, replicas int64, lastName, id string) (pvcs []string) {
	if len(volumes) < 1 {
		return
	}

	tmpl := `
{{- range $index := until (int (default "1" .Replicas)) }}
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: {{ $.Volume.ID }}-` + lastName + `-{{ $index }}
  annotations:
    operator.dameng.com/title: {{ $.Volume.Title | hexenc | quote }}
    operator.dameng.com/describe: {{ $.Volume.Describe | hexenc | quote }}
  labels:
    operator.dameng.com/id: ` + id + `
spec:
  accessModes:
    - {{ default "ReadWriteOnce" $.Volume.AccessMode }}
  resources:
    requests:
      storage: "{{ $.Volume.Capacity }}Gi"
  storageClassName: {{ $.Volume.StorageClass }}
{{- end }}
`

	for _, volume := range volumes {
		if volume.Type == v1.PersistentVolumeClaim {
			parameters := map[string]interface{}{
				"Replicas": replicas,
				"Volume":   volume,
			}
			pvc, _ := defaultParser.ParseString(tmpl, parameters)
			pvcs = append(pvcs, strings.TrimSpace(pvc))
		}
	}
	return pvcs
}

// ParseService 解析Service对象
func ParseService(service *v1.Service, ip string, id string) string {
	if service == nil {
		return ""
	}
	getType := func(ip string) string {
		if ip == "" {
			return "NodePort"
		}
		return "LoadBalancer"
	}
	tmpl := `
kind: Service
apiVersion: v1
metadata:
  labels:
    operator.dameng.com/id: ` + id + `
  name: {{ .Name }}
spec:
  type: ` + getType(ip) + `
  ports:
{{- range .Ports }}
    - name: {{ .Protocol | lower }}-{{ .Port }}
      targetPort: {{ .Port }}
      port: {{ .Port }}
{{- end }}
  selector:
    operator.dameng.com/id: ` + id + `
`

	result, _ := defaultParser.ParseString(tmpl, service)
	return strings.TrimSpace(result)
}

func ParsePodPort(service *v1.Service) string {
	if service == nil {
		return ""
	}
	if len(service.Ports) < 1 {
		return ""
	}

	tmpl := `
ports:
{{- range . }}
  - name: p-{{ .Port }}
    containerPort: {{ .Port }}
    protocol: TCP
{{- end}}
`
	result, _ := defaultParser.ParseString(tmpl, service.Ports)
	return strings.TrimSpace(result)
}

// ParseExporter 解析Prometheus注解
func ParseExporter(exporter *v1.Exporter) string {
	if exporter == nil {
		return ""
	}

	tmpl := `
prometheus.io/scrape: "true"
prometheus.io/path: {{ default "/metrics" .Path | quote }}
prometheus.io/port: {{ default "9100" .Port | quote }}
`

	result, _ := defaultParser.ParseString(tmpl, exporter)
	return strings.TrimSpace(result)
}

// ParseNodeSelector 解析NodeSelector
func ParseNodeSelector(labels []*v1.Label) string {
	if len(labels) < 1 {
		return ""
	}

	tmpl := `
nodeSelector:
{{- range .}}
  {{ .Name }}: {{ .Value | quote }}
{{- end }}
`

	result, _ := defaultParser.ParseString(tmpl, labels)
	return strings.TrimSpace(result)
}

// ParseHostAliases 解析HostAlias
func ParseHostAliases(hostAliases []*v1.HostAlias) string {
	if len(hostAliases) < 1 {
		return ""
	}

	tmpl := `
hostAliases:
{{- range . }}
  - ip: {{ .IP | quote }}
    hostnames:
{{- range .Hostnames }}
      - {{ .Name | quote }}
{{- end }}
{{- end }}
`

	result, _ := defaultParser.ParseString(tmpl, hostAliases)
	return strings.TrimSpace(result)
}

// ParseResources 解析Resource
func ParseResources(r *v1.Resources) string {
	if r == nil {
		return ""
	}

	tmpl := `
resources:
  limits:
{{- range $i, $v := .Limits}}
{{- if eq $i "memory"}}
    memory: {{ $v | quote  }}
{{- end }}
{{- if eq $i "cpu"}}
    cpu: {{ $v | quote  }}
{{- end}}
{{- end}}
  requests:
{{- range $ri, $rv := .Requests}}
{{- if eq $ri "memory"}}
    memory: {{ $rv | quote  }}
{{- end }}
{{- if eq $ri "cpu"}}
    cpu: {{ $rv | quote  }}
{{- end}}
{{- end}}
`

	result, _ := defaultParser.ParseString(tmpl, r)
	return strings.TrimSpace(result)
}

// ParseConfigMap 解析配置
func ParseConfigMap(configs []*v1.Configuration, lastName, id string) string {
	if len(configs) < 1 {
		return ""
	}

	tmpl := `
apiVersion: v1
kind: ConfigMap
metadata:
  name: ` + lastName + `
  labels:
    operator.dameng.com/id: ` + id + `
data:
{{- range . }}
  {{ .MountPoint | md5 }}: |
    {{- nindent 4 .Content }}
{{- end }}
`

	result, _ := defaultParser.ParseString(tmpl, configs)
	return strings.TrimSpace(result)
}

// ParseVolumeMounts 解析挂载点
func ParseVolumeMounts(volumes []*v1.Volume, configs []*v1.Configuration, logs []*v1.Log) string {
	if len(volumes) < 1 && len(configs) < 1 && len(logs) < 1 {
		return ""
	}

	parameters := map[string]interface{}{
		"Volumes": volumes,
		"Configs": configConverter(configs),
		"Logs":    logs,
	}

	tmpl := `
volumeMounts:
{{- range .Volumes }}
  - name: {{ .ID }}
    mountPath: {{ .MountPoint }}
{{- end }}
{{- range $path, $files := .Configs }}
  - name: {{ $path | md5 }}
    mountPath: {{ $path }}
{{- end }}
{{- range .Logs }}
  - name: {{ .ID }}
    mountPath: {{ .Directory }}
    json: {{ toJson . }}
{{- end }}
`

	result, _ := defaultParser.ParseString(tmpl, parameters)
	return strings.TrimSpace(result)
}

// ParseVolumes 解析卷
func ParseVolumes(volumes []*v1.Volume, configs []*v1.Configuration, logs []*v1.Log, lastName, id string) string {
	if len(volumes) < 1 && len(configs) < 1 && len(logs) < 1 {
		return ""
	}

	parameters := map[string]interface{}{
		"Volumes": volumes,
		"Configs": configConverter(configs),
		"Logs":    logs,
	}

	tmpl := `
volumes:
{{- range $path, $files := .Configs }}
  - name: {{ $path | md5 }}
    configMap:
      name: ` + lastName + `
      items:
{{- range $files }}
        - key: {{ cat $path "/" . | nospace | md5 }}
          path: {{ . }}
{{- end }}
{{- end }}
{{- range .Volumes }}
{{- if eq .Type "EmptyDir" }}
  - name: {{ .ID }}
    emptyDir:
      medium: Memory
      sizeLimit: {{ .Capacity }}Gi
{{- end }}
{{- if eq .Type "HostPath" }}
  - name: {{ .ID }}
    hostPath:
      path: "{{ .Location }}"
      type: "DirectoryOrCreate"
{{- end }}
{{- end }}
{{- range .Logs }}
  - name: {{ .ID }}
    hostPath:
      path: "{{ default "/logs" (env "LOG_DIRECTORY_LOCATION") }}/` + lastName + "-" + id + `{{ .Directory }}"
      type: "DirectoryOrCreate"
{{- end }}
`

	result, _ := defaultParser.ParseString(tmpl, parameters)
	return strings.TrimSpace(result)
}

// ParseDestinationRule 解析目标规则
func ParseDestinationRule(service *v1.Service, connectionPool *v1.ConnectionPool, outlierDetection *v1.OutlierDetection, loadBalancer *v1.LoadBalancer, lastName, id, title string) string {
	if service == nil {
		return ""
	}

	parameters := map[string]interface{}{
		"Service":          service,
		"ConnectionPool":   connectionPool,
		"OutlierDetection": outlierDetection,
		"LoadBalancer":     loadBalancer,
	}

	tmpl := `
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: {{ .Service.Name }}
  labels:
    operator.dameng.com/id: ` + id + `
  annotations:
    operator.dameng.com/title: "` + HexEncode(title) + `"
spec:
  host: {{ .Service.Name }}
{{- if or .ConnectionPool .OutlierDetection .LoadBalance }}
  trafficPolicy:
{{- with .ConnectionPool }}
    connectionPool:
{{- if or .MaxConnections .ConnectTimeout }}
      tcp:
{{- with .MaxConnections }}
        maxConnections: {{ . }}
{{- end }}
{{- with .ConnectTimeout }}
        connectTimeout: {{ . }}ms
{{- end }}
{{- end }}
{{- if or .Http1MaxPendingRequests .Http2MaxRequests .MaxRequestsPerConnection }}
      http:
{{- with .Http1MaxPendingRequests }}
        http1MaxPendingRequests: {{ . }}
{{- end }}
{{- with .Http2MaxRequests }}
        http2MaxRequests: {{ . }}
{{- end }}
{{- with .MaxRequestsPerConnection }}
        maxRequestsPerConnection: {{ . }}
{{- end }}
{{- with .IdleTimeout }}
        idleTimeout: {{ . }}s
{{- end }}
{{- end }}
{{- end }}
{{- with .OutlierDetection }}
    outlierDetection:
{{- if eq .Type "consecutiveGatewayErrors" }}
      consecutiveGatewayErrors: {{ .Consecutive }}
{{- end }}
{{- if eq .Type "consecutive5xxErrors" }}
      consecutive5xxErrors: {{ .Consecutive }}
{{- end }}
{{- with .Interval }}
      interval: {{ . }}s
{{- end }}
{{- with .BaseEjectionTime }}
      baseEjectionTime: {{ . }}s
{{- end }}
{{- with .MaxEjectionPercent }}
      maxEjectionPercent: {{ . }}
{{- end }}
{{- with .MinHealthPercent }}
      minHealthPercent: {{ . }}
{{- end }}
{{- end }}
{{- with .LoadBalancer }}
    loadBalancer:
{{- if eq .Type "simple" }}
      simple: {{ .Policy }}
{{- else }}
      consistentHash:
{{- if .UseSourceIp }}
        useSourceIp: true
{{- end }}
{{- if .Header }}
        httpHeaderName: {{ .Header | quote }}
{{- end }}
{{- with .Cookie }}
        httpCookie:
          name: {{ default "istio-cookie-hash" .Name }}
          path: {{ default "/" .Path }}
          ttl: {{ default "3600" .TTL }}s
{{- end }}
{{- with .HttpQueryParameterName }}
        httpQueryParameterName: {{ . | quote }}
{{- end }}
{{- end }}
{{- end }}
{{- end }}
`

	result, _ := defaultParser.ParseString(tmpl, parameters)
	return strings.TrimSpace(result)
}

func ParseArguments(arguments []*v1.Argument) string {
	if len(arguments) < 1 {
		return ""
	}

	tmpl := `
args:
{{- range . }}
  - {{ .Value | quote }}
{{- end }}
`

	result, _ := defaultParser.ParseString(tmpl, arguments)
	return strings.TrimSpace(result)
}

// ParseUpdateStrategy 解析升级策略
func ParseUpdateStrategy(volumes []*v1.Volume, updateStrategy *v1.UpdateStrategy, podManagementPolicy *v1.PodManagementPolicy, uniqueness bool) string {
	if updateStrategy == nil {
		return ""
	}

	field := "strategy"
	strategy := *updateStrategy

	kind := parseKind(volumes, podManagementPolicy, uniqueness)
	if kind == "StatefulSet" {
		field = "updateStrategy"
	}

	if *updateStrategy == v1.RecreateStrategy && field == "updateStrategy" {
		strategy = v1.OnDeleteStrategy
	}

	tmpl := `
` + field + `:
  type: ` + string(strategy) + `
`
	return strings.TrimSpace(tmpl)
}

// ParsePodManagementPolicy 明确是否需要podManagementPolicy字段
func ParsePodManagementPolicy(podManagementPolicy *v1.PodManagementPolicy) string {
	if podManagementPolicy != nil && *podManagementPolicy == v1.OrderedReadyPolicy {
		return "podManagementPolicy: OrderedReady"
	}
	return ""
}

// ParseKind 判定部署类型
func ParseKind(volumes []*v1.Volume, podManagementPolicy *v1.PodManagementPolicy, uniqueness bool) string {
	return "kind: " + parseKind(volumes, podManagementPolicy, uniqueness)
}

// ParseCommand 解析command
func ParseCommand(command string) string {
	if command == "" {
		return ""
	}
	return `command:
  - ` + command
}

// ParseHandler 解析lifecycle和probe中的handler
func ParseHandler(handler *v1.Handler) string {
	if handler == nil {
		return ""
	}

	tmpl := `{{- if eq .Action "Exec" -}}
exec:
  command:
    - /bin/sh
    - -c
    - {{ .Command }}
{{- else if eq .Action "HTTPGet" -}}
httpGet:
{{- with .Scheme }}
  scheme: {{ . }}
{{- end }}
{{- with .Host }}
  host: {{ . }}
{{- end }}
  port: {{ default 8080 .Port }}
{{- with .Path }}
  path: {{ . }}
{{- end }}
{{- with .Headers }}
  httpHeaders:
{{- range . }}
    - name: {{ .Name | quote }}
      value: {{ .Value | quote }}
{{- end }}
{{- end }}
{{- else -}}
tcpSocket:
  port: {{ default 8080 .Port }}
{{- with .Host }}
  host: {{ . }}
{{- end }}
{{- end }}`

	result, _ := defaultParser.ParseString(tmpl, handler)
	return strings.TrimSpace(result)
}

// ParseLifecycle 解析lifecycle
func ParseLifecycle(terminator *v1.Terminator) string {
	if terminator == nil {
		return ""
	}

	tmpl := `
lifecycle:
  {{- parseHandler .Handler | nindent 2 }}
`

	result, _ := defaultParser.ParseString(tmpl, terminator)
	return strings.TrimSpace(result)
}

// ParseProbe 解析探针
func ParseProbe(readinessProbe, livenessProbe, startupProbe *v1.Probe) string {
	if readinessProbe == nil && livenessProbe == nil && startupProbe == nil {
		return ""
	}

	parameters := map[string]interface{}{
		"ReadinessProbe": readinessProbe,
		"LivenessProbe":  livenessProbe,
		"StartupProbe":   startupProbe,
	}

	subTmpl := `{{- with .InitialDelaySeconds }}
  initialDelaySeconds: {{ . }}
  {{- end }}
  {{- with .TimeoutSeconds }}
  timeoutSeconds: {{ . }}
  {{- end }}
  {{- with .PeriodSeconds }}
  periodSeconds: {{ . }}
  {{- end }}
  {{- with .SuccessThreshold }}
  successThreshold: {{ . }}
  {{- end }}
  {{- with .FailureThreshold }}
  failureThreshold: {{ . }}
  {{- end -}}
{{ parseHandler .Handler | nindent 2 }}`

	tmpl := `
{{- with .ReadinessProbe }}
startupProbe:
` + subTmpl + `
{{- end }}
{{- with .LivenessProbe }}
livenessProbe:
` + subTmpl + `
{{- end }}
{{- with .StartupProbe }}
readinessProbe:
` + subTmpl + `
{{- end }}
`

	result, _ := defaultParser.ParseString(tmpl, parameters)
	return strings.TrimSpace(result)
}

// ParseEnvironments 解析环境变量
func ParseEnvironments(environments []*v1.Environment) string {
	tmpl := `
{{- with . -}}
env:
{{- range . }}
  - name: {{ .Name | quote }}
    value: {{ .Value | quote }}
{{- end }}
{{- end }}`

	result, _ := defaultParser.ParseString(tmpl, environments)
	return strings.TrimSpace(result)
}

// ParseTerminationGracePeriodSeconds 解析terminationGracePeriodSeconds
func ParseTerminationGracePeriodSeconds(terminator *v1.Terminator) string {
	if terminator == nil {
		return ""
	}
	tmpl := `terminationGracePeriodSeconds: {{ default 30 .Grace }}`
	result, _ := defaultParser.ParseString(tmpl, terminator)
	return strings.TrimSpace(result)
}

// ParseImage 解析镜像
func ParseImage(image *v1.Image) string {
	if image == nil {
		return ""
	}

	tmpl := `
image: {{ .Registry }}/{{ .Repository }}:{{ .Tag }}
imagePullPolicy: IfNotPresent
`

	result, _ := defaultParser.ParseString(tmpl, image)
	return strings.TrimSpace(result)
}

// ParseAutoscaler 解析HPA
func ParseAutoscaler(autoscaler *v1.Autoscaler, kind, lastName, id string) string {
	if autoscaler == nil {
		return ""
	}

	tmpl := `
apiVersion: autoscaling/v2beta2
kind: HorizontalPodAutoscaler
metadata:
  name: ` + lastName + `
  labels:
    operator.dameng.com/id: ` + id + `
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    ` + kind + `
    name: ` + lastName + `
  minReplicas: {{ default "1" .MinReplicas }}
  maxReplicas: {{ default "10" .MaxReplicas }}
  metrics:
{{- range .Metrics }}
{{- if eq .Name "CPUUtilization" }}
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: {{ .Value }}
{{- end }}
{{- if eq .Name "MemoryUtilization" }}
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: {{ .Value }}
{{- end }}
{{- end }}
`

	result, _ := defaultParser.ParseString(tmpl, autoscaler)
	return strings.TrimSpace(result)
}

func ParseWorkingDirectory(workingDirectory string) string {
	if workingDirectory == "" {
		return ""
	}
	return fmt.Sprintf("workingDir: %s", workingDirectory)
}

func ParseReplicas(power *v1.Power, replicas int64) string {
	if power == nil || *power == v1.PowerOn {
		return fmt.Sprintf("replicas: %v", func() int64 {
			if replicas == 0 {
				return 1
			}
			return replicas
		}())
	}
	return "replicas: 0"
}

// ParsePodLabels 生成Pod使用的Label
func ParsePodLabels(labels []*v1.Label, id string) string {
	if len(labels) < 1 || id == "" {
		return ""
	}

	tmpl := `
{{- range . }}
operator.dameng.com/tenant-{{ .Name | hexenc }}: {{ .Value | hexenc | quote }}
{{- end }}
`

	result, _ := defaultParser.ParseString(tmpl, labels)
	return strings.TrimSpace(result)
}

// ParseVolumeClaimTemplates 解析volume模版
func ParseVolumeClaimTemplates(volumes []*v1.Volume, id string) string {
	var vcts []*v1.Volume
	for _, volume := range volumes {
		if volume.Type == v1.PersistentVolumeClaim {
			vcts = append(vcts, volume)
		}
	}
	if len(vcts) < 1 {
		return ""
	}

	tmpl := `
volumeClaimTemplates:
{{- range . }}
  - metadata:
      name: {{ .ID }}
      annotations:
        operator.dameng.com/title: {{ .Title | hexenc | quote }}
        operator.dameng.com/describe: {{ .Describe | hexenc | quote }}
      labels:
        operator.dameng.com/id: ` + id + `
    spec:
      accessModes:
        - {{ default "ReadWriteOnce" .AccessMode }}
      storageClassName: {{ .StorageClass }}
      resources:
        requests:
          storage: 1Gi
{{- end }}
`

	result, _ := defaultParser.ParseString(tmpl, vcts)
	return strings.TrimSpace(result)
}

// ParseImagePullSecrets todo: 尚未完成前端设计
func ParseImagePullSecrets() string {
	return ""
}

// ParseAffinity todo: 尚未完成前端设计
func ParseAffinity() string {
	return ""
}

// ParseHeadlessService 解析HeadlessService
func ParseHeadlessService(service *v1.Service, volumes []*v1.Volume, podManagementPolicy *v1.PodManagementPolicy, uniqueness bool, lastName, id string) string {
	kind := parseKind(volumes, podManagementPolicy, uniqueness)
	if kind != "StatefulSet" {
		return ""
	}

	tmpl := `
kind: Service
apiVersion: v1
metadata:
  labels:
    operator.dameng.com/id: ` + id + `
  name: ` + lastName + `
spec:
  clusterIP: None
{{- with . }}
  ports:
{{- range .Ports }}
    - name: {{ .Protocol | lower }}-{{ .Port }}
      targetPort: {{ .Port }}
      port: {{ .Port }}
{{- end }}
{{- end }}
  selector:
    operator.dameng.com/id: ` + id + `
`

	result, _ := defaultParser.ParseString(tmpl, service)
	return strings.TrimSpace(result)
}

func ParseServiceName(volumes []*v1.Volume, podManagementPolicy *v1.PodManagementPolicy, uniqueness bool, lastName string) string {
	kind := parseKind(volumes, podManagementPolicy, uniqueness)
	if kind == "StatefulSet" {
		return fmt.Sprintf("serviceName: %s", lastName)
	}
	return ""
}

// configConverter 配置文件格式转化逻辑
func configConverter(configs []*v1.Configuration) map[string][]string {
	files := make(map[string][]string)
	for _, config := range configs {
		name := path.Base(config.MountPoint)
		dir := path.Dir(config.MountPoint)
		if names, ok := files[dir]; ok {
			files[dir] = append(names, name)
		} else {
			files[dir] = []string{name}
		}
	}
	return files
}

// parseKind 判定部署类型
func parseKind(volumes []*v1.Volume, podManagementPolicy *v1.PodManagementPolicy, uniqueness bool) string {
	if (podManagementPolicy != nil && *podManagementPolicy == v1.OrderedReadyPolicy) || uniqueness {
		return "StatefulSet"
	}
	if len(volumes) > 0 {
		return "StatefulSet"
	}
	return "Deployment"
}

type KubernetesConf struct {
	PersistentVolumeClaims        []string
	ConfigMap                     string
	Kind                          string
	Labels                        string
	Exporter                      string
	UpdateStrategy                string
	PodManagementPolicy           string
	ServiceName                   string
	NodeSelector                  string
	HostAliases                   string
	Replicas                      string
	ImagePullSecrets              string
	Affinity                      string
	Image                         string
	Command                       string
	Arguments                     string
	Ports                         string
	Environments                  string
	Resources                     string
	Lifecycle                     string
	Probe                         string
	WorkingDirectory              string
	VolumeMounts                  string
	Volumes                       string
	VolumeClaimTemplates          string
	TerminationGracePeriodSeconds string
	Service                       string
	HeadlessService               string
	Autoscaler                    string
	DestinationRule               string
}

// CompassHelper Compass模版助手
func KubernetesConfHelper(object interface{}) (conf KubernetesConf) {
	entire := reflect.ValueOf(object)
	metadata, ok := entire.FieldByName("ObjectMeta").Interface().(metav1.ObjectMeta)

	lastName := metadata.GetName()
	klog.V(4).InfoS("template helper", "lastName", lastName)

	fields := entire.FieldByName("Spec")

	id := fields.FieldByName("ID").String()
	klog.V(4).InfoS(fmt.Sprintf("template helper: %s", lastName), "id", id)

	title := fields.FieldByName("Title").String()
	klog.V(4).InfoS(fmt.Sprintf("template helper: %s", lastName), "title", title)

	configs, ok := fields.FieldByName("Configurations").Interface().([]*v1.Configuration)
	klog.V(4).InfoS(fmt.Sprintf("template helper: %s", lastName), "ok", ok, "configs", configs)

	volumes, ok := fields.FieldByName("Volumes").Interface().([]*v1.Volume)
	klog.V(4).InfoS(fmt.Sprintf("template helper: %s", lastName), "ok", ok, "volumes", volumes)

	exporter, ok := fields.FieldByName("Exporter").Interface().(*v1.Exporter)
	klog.V(4).InfoS(fmt.Sprintf("template helper: %s", lastName), "ok", ok, "exporter", exporter)

	replicas := fields.FieldByName("Replicas").Int()
	klog.V(4).InfoS(fmt.Sprintf("template helper: %s", lastName), "ok", ok, "replicas", replicas)

	power := fields.FieldByName("Power").Interface().(*v1.Power)
	klog.V(4).InfoS(fmt.Sprintf("template helper: %s", lastName), "ok", ok, "power", power)

	command := fields.FieldByName("Command").String()
	klog.V(4).InfoS(fmt.Sprintf("template helper: %s", lastName), "ok", ok, "command", command)

	arguments, ok := fields.FieldByName("Arguments").Interface().([]*v1.Argument)
	klog.V(4).InfoS(fmt.Sprintf("template helper: %s", lastName), "ok", ok, "arguments", arguments)

	environments, ok := fields.FieldByName("Environments").Interface().([]*v1.Environment)
	klog.V(4).InfoS(fmt.Sprintf("template helper: %s", lastName), "ok", ok, "environments", environments)

	terminator, ok := fields.FieldByName("Terminator").Interface().(*v1.Terminator)
	klog.V(4).InfoS(fmt.Sprintf("template helper: %s", lastName), "ok", ok, "terminator", terminator)

	logs, ok := fields.FieldByName("Logs").Interface().([]*v1.Log)
	klog.V(4).InfoS(fmt.Sprintf("template helper: %s", lastName), "ok", ok, "logs", logs)

	service, ok := fields.FieldByName("Service").Interface().(*v1.Service)
	klog.V(4).InfoS(fmt.Sprintf("template helper: %s", lastName), "ok", ok, "service", service)

	updateStrategy := fields.FieldByName("UpdateStrategy").Interface().(*v1.UpdateStrategy)
	klog.V(4).InfoS(fmt.Sprintf("template helper: %s", lastName), "ok", ok, "updateStrategy", updateStrategy)

	podManagementPolicy := fields.FieldByName("PodManagementPolicy").Interface().(*v1.PodManagementPolicy)
	klog.V(4).InfoS(fmt.Sprintf("template helper: %s", lastName), "podManagementPolicy", podManagementPolicy)

	//memory, ok := fields.FieldByName("Memory").Interface().(*v12.Resource)
	//klog.V(4).InfoS(fmt.Sprintf("template helper: %s", lastName), "ok", ok, "memory", memory)
	//
	//cpu, ok := fields.FieldByName("CPU").Interface().(*v12.Resource)
	//klog.V(4).InfoS(fmt.Sprintf("template helper: %s", lastName), "ok", ok, "cpu", cpu)

	nodeSelector, ok := fields.FieldByName("NodeSelector").Interface().([]*v1.Label)
	klog.V(4).InfoS(fmt.Sprintf("template helper: %s", lastName), "ok", ok, "nodeSelector", nodeSelector)

	hostAliases, ok := fields.FieldByName("HostAliases").Interface().([]*v1.HostAlias)
	klog.V(4).InfoS(fmt.Sprintf("template helper: %s", lastName), "ok", ok, "hostAliases", hostAliases)

	loadBalancer, ok := fields.FieldByName("LoadBalancer").Interface().(*v1.LoadBalancer)
	klog.V(4).InfoS(fmt.Sprintf("template helper: %s", lastName), "ok", ok, "loadBalancer", loadBalancer)

	connectionPool, ok := fields.FieldByName("ConnectionPool").Interface().(*v1.ConnectionPool)
	klog.V(4).InfoS(fmt.Sprintf("template helper: %s", lastName), "ok", ok, "connectionPool", connectionPool)

	outlierDetection, ok := fields.FieldByName("OutlierDetection").Interface().(*v1.OutlierDetection)
	klog.V(4).InfoS(fmt.Sprintf("template helper: %s", lastName), "ok", ok, "outlierDetection", outlierDetection)

	autoscaler, ok := fields.FieldByName("Autoscaler").Interface().(*v1.Autoscaler)
	klog.V(4).InfoS(fmt.Sprintf("template helper: %s", lastName), "ok", ok, "autoscaler", autoscaler)

	readinessProbe, ok := fields.FieldByName("ReadinessProbe").Interface().(*v1.Probe)
	klog.V(4).InfoS(fmt.Sprintf("template helper: %s", lastName), "ok", ok, "readinessProbe", readinessProbe)

	livenessProbe, ok := fields.FieldByName("LivenessProbe").Interface().(*v1.Probe)
	klog.V(4).InfoS(fmt.Sprintf("template helper: %s", lastName), "ok", ok, "livenessProbe", livenessProbe)

	startupProbe, ok := fields.FieldByName("StartupProbe").Interface().(*v1.Probe)
	klog.V(4).InfoS(fmt.Sprintf("template helper: %s", lastName), "ok", ok, "startupProbe", startupProbe)

	ip := fields.FieldByName("IP").String()
	klog.V(4).InfoS(fmt.Sprintf("template helper: %s", lastName), "ip", ip)

	image, ok := fields.FieldByName("Image").Interface().(*v1.Image)
	klog.V(4).InfoS(fmt.Sprintf("template helper: %s", lastName), "ok", ok, "image", image)

	workingDirectory := fields.FieldByName("WorkingDirectory").String()
	klog.V(4).InfoS(fmt.Sprintf("template helper: %s", lastName), "workingDirectory", workingDirectory)

	labels, ok := fields.FieldByName("Labels").Interface().([]*v1.Label)
	klog.V(4).InfoS(fmt.Sprintf("template helper: %s", lastName), "ok", ok, "labels", labels)

	rs, ok := fields.FieldByName("Resources").Interface().(*v1.Resources)
	klog.V(4).InfoS(fmt.Sprintf("template helper: %s", lastName), "ok", ok, "Resources", rs)

	conf.ConfigMap = ParseConfigMap(configs, lastName, id)
	klog.V(5).InfoS(fmt.Sprintf("template helper: %s", lastName), "ConfigMap", conf.ConfigMap)

	conf.PersistentVolumeClaims = ParsePersistentVolumeClaim(volumes, replicas, lastName, id)
	klog.V(5).InfoS(fmt.Sprintf("template helper: %s", lastName), "PersistentVolumeClaims", conf.PersistentVolumeClaims)

	conf.Kind = ParseKind(volumes, podManagementPolicy, false)
	klog.V(5).InfoS(fmt.Sprintf("template helper: %s", lastName), "Kind", conf.Kind)

	conf.Labels = ParsePodLabels(labels, id)
	klog.V(5).InfoS(fmt.Sprintf("template helper: %s", lastName), "Labels", conf.Labels)

	conf.Exporter = ParseExporter(exporter)
	klog.V(5).InfoS(fmt.Sprintf("template helper: %s", lastName), "Exporter", conf.Exporter)

	conf.UpdateStrategy = ParseUpdateStrategy(volumes, updateStrategy, podManagementPolicy, false)
	klog.V(5).InfoS(fmt.Sprintf("template helper: %s", lastName), "UpdateStrategy", conf.UpdateStrategy)

	conf.PodManagementPolicy = ParsePodManagementPolicy(podManagementPolicy)
	klog.V(5).InfoS(fmt.Sprintf("template helper: %s", lastName), "PodManagementPolicy", conf.PodManagementPolicy)

	conf.NodeSelector = ParseNodeSelector(nodeSelector)
	klog.V(5).InfoS(fmt.Sprintf("template helper: %s", lastName), "NodeSelector", conf.NodeSelector)

	conf.Replicas = ParseReplicas(power, replicas)
	klog.V(5).InfoS(fmt.Sprintf("template helper: %s", lastName), "Replicas", conf.Replicas)

	// todo
	conf.Affinity = ParseAffinity()
	klog.V(5).InfoS(fmt.Sprintf("template helper: %s", lastName), "Affinity", conf.Affinity)

	// todo
	conf.ImagePullSecrets = ParseImagePullSecrets()
	klog.V(5).InfoS(fmt.Sprintf("template helper: %s", lastName), "ImagePullSecrets", conf.ImagePullSecrets)

	conf.Image = ParseImage(image)
	klog.V(5).InfoS(fmt.Sprintf("template helper: %s", lastName), "Image", conf.Image)

	conf.Command = ParseCommand(command)
	klog.V(5).InfoS(fmt.Sprintf("template helper: %s", lastName), "Command", conf.Command)

	conf.Arguments = ParseArguments(arguments)
	klog.V(5).InfoS(fmt.Sprintf("template helper: %s", lastName), "Arguments", conf.Arguments)

	conf.Ports = ParsePodPort(service)
	klog.V(5).InfoS(fmt.Sprintf("template helper: %s", lastName), "Ports", conf.Ports)

	conf.WorkingDirectory = ParseWorkingDirectory(workingDirectory)
	klog.V(5).InfoS(fmt.Sprintf("template helper: %s", lastName), "WorkingDirectory", conf.WorkingDirectory)

	conf.HostAliases = ParseHostAliases(hostAliases)
	klog.V(5).InfoS(fmt.Sprintf("template helper: %s", lastName), "HostAliases", conf.HostAliases)

	conf.Environments = ParseEnvironments(environments)
	klog.V(5).InfoS(fmt.Sprintf("template helper: %s", lastName), "Environments", conf.Environments)

	conf.Resources = ParseResources(rs)
	klog.V(5).InfoS(fmt.Sprintf("template helper: %s", lastName), "Resources", conf.Resources)

	conf.Lifecycle = ParseLifecycle(terminator)
	klog.V(5).InfoS(fmt.Sprintf("template helper: %s", lastName), "Lifecycle", conf.Lifecycle)

	conf.Probe = ParseProbe(readinessProbe, livenessProbe, startupProbe)
	klog.V(5).InfoS(fmt.Sprintf("template helper: %s", lastName), "Probe", conf.Probe)

	conf.VolumeMounts = ParseVolumeMounts(volumes, configs, logs)
	klog.V(5).InfoS(fmt.Sprintf("template helper: %s", lastName), "VolumeMounts", conf.VolumeMounts)

	conf.Volumes = ParseVolumes(volumes, configs, logs, lastName, id)
	klog.V(5).InfoS(fmt.Sprintf("template helper: %s", lastName), "Volumes", conf.Volumes)

	conf.VolumeClaimTemplates = ParseVolumeClaimTemplates(volumes, id)
	klog.V(5).InfoS(fmt.Sprintf("template helper: %s", lastName), "VolumeClaimTemplates", conf.VolumeClaimTemplates)

	conf.TerminationGracePeriodSeconds = ParseTerminationGracePeriodSeconds(terminator)
	klog.V(5).InfoS(fmt.Sprintf("template helper: %s", lastName), "TerminationGracePeriodSeconds", conf.TerminationGracePeriodSeconds)

	conf.Service = ParseService(service, ip, id)
	klog.V(5).InfoS(fmt.Sprintf("template helper: %s", lastName), "Service", conf.Service)

	conf.HeadlessService = ParseHeadlessService(service, volumes, podManagementPolicy, false, lastName, id)
	klog.V(5).InfoS(fmt.Sprintf("template helper: %s", lastName), "HeadlessService", conf.HeadlessService)

	conf.ServiceName = ParseServiceName(volumes, podManagementPolicy, false, lastName)
	klog.V(5).InfoS(fmt.Sprintf("template helper: %s", lastName), "ServiceName", conf.ServiceName)

	conf.Autoscaler = ParseAutoscaler(autoscaler, conf.Kind, lastName, id)
	klog.V(5).InfoS(fmt.Sprintf("template helper: %s", lastName), "Autoscaler", conf.Autoscaler)

	conf.DestinationRule = ParseDestinationRule(service, connectionPool, outlierDetection, loadBalancer, lastName, id, title)
	klog.V(5).InfoS(fmt.Sprintf("template helper: %s", lastName), "DestinationRule", conf.DestinationRule)

	return conf
}
