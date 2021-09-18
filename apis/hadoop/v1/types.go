/**
 @author: ZHC
 @date: 2021-09-09 10:33:25
 @description:
**/
package v1

import (
	corev1 "k8s.io/api/core/v1"
)

type UpdateStrategy string

const (
	RecreateStrategy UpdateStrategy = "Recreate"
	RollingStrategy  UpdateStrategy = "RollingUpdate"
	OnDeleteStrategy UpdateStrategy = "OnDelete"
)

type PodManagementPolicy string

const (
	OrderedReadyPolicy PodManagementPolicy = "OrderedReady"
	ParallelPolicy     PodManagementPolicy = "Parallel"
)

type Power string

const (
	PowerOn  Power = "on"
	PowerOff Power = "off"
)

type Label struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type VolumeType string

const (
	HostPath              VolumeType = "HostPath"
	EmptyDir              VolumeType = "EmptyDir"
	PersistentVolumeClaim VolumeType = "PersistentVolumeClaim"
)

type VolumeCreateType string

const (
	Automatic VolumeCreateType = "automatic"
	Manual    VolumeCreateType = "manual"
)

type Volume struct {
	ID           string                            `json:"id"`
	Type         VolumeType                        `json:"type,omitempty"`
	Name         string                            `json:"name,omitempty"`
	Title        string                            `json:"title,omitempty"`
	Describe     string                            `json:"describe,omitempty"`
	Namespace    string                            `json:"namespace,omitempty"`
	MountPoint   string                            `json:"mountPoint,omitempty"`
	StorageClass string                            `json:"storageClass,omitempty"`
	Capacity     string                            `json:"capacity,omitempty"`
	AccessMode   corev1.PersistentVolumeAccessMode `json:"accessMode,omitempty"`
	Location     string                            `json:"location,omitempty"`
	CreateType   VolumeCreateType                  `json:"createType,omitempty"`
}

type Port struct {
	Port     string `json:"port,omitempty"`
	Protocol string `json:"protocol,omitempty"`
}

type Service struct {
	Name  string  `json:"name,omitempty"`
	Ports []*Port `json:"ports,omitempty"`
}

type Exporter struct {
	Path string `json:"path,omitempty"`
	Port int32  `json:"port,omitempty"`
}

type Hostname struct {
	Name string `json:"name,omitempty"`
}

type HostAlias struct {
	ID        string      `json:"id,omitempty"`
	IP        string      `json:"ip,omitempty"`
	Hostnames []*Hostname `json:"hostnames,omitempty"`
}

type Resources struct {
	// +optional
	// +kubebuilder:default:={cpu: "100m", memory: "256Mi"}
	Requests map[corev1.ResourceName]string `json:"requests,omitempty"`
	// +optional
	// +kubebuilder:default:={cpu: "500m", memory: "1Gi"}
	Limits map[corev1.ResourceName]string `json:"limits,omitempty"`
}

type Configuration struct {
	ID          string `json:"id,omitempty"`
	MountPoint  string `json:"mountPoint,omitempty"`
	Content     string `json:"content,omitempty"`
	NeedRestart bool   `json:"needRestart,omitempty"`
}

type Log struct {
	ID        string `json:"id,omitempty"`
	Directory string `json:"directory,omitempty"`
	Pattern   string `json:"pattern,omitempty"`
}

type ConnectionPool struct {
	MaxConnections           int64 `json:"maxConnections,omitempty"`
	ConnectTimeout           int64 `json:"connectTimeout,omitempty"`
	Http1MaxPendingRequests  int64 `json:"http1MaxPendingRequests,omitempty"`
	Http2MaxRequests         int64 `json:"http2MaxRequests,omitempty"`
	MaxRequestsPerConnection int64 `json:"maxRequestsPerConnection,omitempty"`
	IdleTimeout              int64 `json:"idleTimeout,omitempty"`
}

type OutlierDetection struct {
	Type               string `json:"type,omitempty"`
	Consecutive        int64  `json:"consecutive,omitempty"`
	Interval           int64  `json:"interval,omitempty"`
	BaseEjectionTime   int64  `json:"baseEjectionTime,omitempty"`
	MaxEjectionPercent int64  `json:"maxEjectionPercent,omitempty"`
	MinHealthPercent   int64  `json:"minHealthPercent,omitempty"`
}

type Cookie struct {
	Name string `json:"name,omitempty"`
	Path string `json:"path,omitempty"`
	TTL  int64  `json:"ttl,omitempty"`
}

type LoadBalancer struct {
	Type                   string  `json:"type,omitempty"`
	Policy                 string  `json:"policy,omitempty"`
	Header                 string  `json:"header,omitempty"`
	Cookie                 *Cookie `json:"cookie,omitempty"`
	UseSourceIp            bool    `json:"useSourceIp,omitempty"`
	HttpQueryParameterName string  `json:"httpQueryParameterName,omitempty"`
}

type Argument struct {
	ID    string `json:"id,omitempty"`
	Value string `json:"value,omitempty"`
}

type Image struct {
	Registry   string `json:"registry,omitempty"`
	Repository string `json:"repository,omitempty"`
	Tag        string `json:"tag,omitempty"`
	Digest     string `json:"digest,omitempty"`
}

type ActionType string

const (
	HTTPGet   ActionType = "HTTPGet"
	Exec      ActionType = "Exec"
	TCPSocket ActionType = "TCPSocket"
)

type Header struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type Handler struct {
	Action  ActionType `json:"action,omitempty"`
	Scheme  string     `json:"scheme,omitempty"`
	Host    string     `json:"host,omitempty"`
	Port    int32      `json:"port,omitempty"`
	Path    string     `json:"path,omitempty"`
	Headers []*Header  `json:"headers,omitempty"`
	Command string     `json:"command,omitempty"`
}

type Terminator struct {
	Grace   int32    `json:"grace,omitempty"`
	Handler *Handler `json:"handler,omitempty"`
}

type ProbeType string

const (
	ReadinessProbe ProbeType = "readiness"
	LivenessProbe  ProbeType = "liveness"
	StartupProbe   ProbeType = "startup"
)

type Probe struct {
	Type                ProbeType `json:"type,omitempty"`
	InitialDelaySeconds int32     `json:"initialDelaySeconds,omitempty"`
	TimeoutSeconds      int32     `json:"timeoutSeconds,omitempty"`
	PeriodSeconds       int32     `json:"periodSeconds,omitempty"`
	SuccessThreshold    int32     `json:"successThreshold,omitempty"`
	FailureThreshold    int32     `json:"failureThreshold,omitempty"`
	Handler             *Handler  `json:"handler,omitempty"`
}

type Environment struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type Metric struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Value int32  `json:"value,omitempty"`
}

type Autoscaler struct {
	MinReplicas int32     `json:"minReplicas,omitempty"`
	MaxReplicas int32     `json:"maxReplicas,omitempty"`
	Metrics     []*Metric `json:"metrics,omitempty"`
}

type Registry struct {
	ID       string `json:"id,omitempty"`
	Title    string `json:"title,omitempty"`
	URL      string `json:"url,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type Capacity struct {
	CPU     string `json:"cpu,omitempty"`
	Memory  string `json:"memory,omitempty"`
	Storage string `json:"storage,omitempty"`
}

type Threshold struct {
	CPU     int32 `json:"cpu,omitempty"`
	Memory  int32 `json:"memory,omitempty"`
	Storage int32 `json:"storage,omitempty"`
}

type Zone struct {
	ID        string     `json:"id,omitempty"`
	Name      string     `json:"name,omitempty"`
	Title     string     `json:"title,omitempty"`
	Describe  string     `json:"describe,omitempty"`
	Capacity  *Capacity  `json:"capacity,omitempty"`
	Default   *Capacity  `json:"default,omitempty"`
	Threshold *Threshold `json:"threshold,omitempty"`
	Usage     *Capacity  `json:"usage,omitempty"`
}

type Phase string

const (
	//调合
	Reconciling Phase = "Reconciling"
	Ready       Phase = "Ready"
	ShutDown    Phase = "ShutDown"
	Running     Phase = "Running"
	Deleting    Phase = "Deleting"
)
