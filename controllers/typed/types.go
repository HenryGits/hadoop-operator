///**
// @author: ZHC
// @date: 2021-09-09 10:32:27
// @description:
//**/
package typed

//
//import (
//	corev1 "k8s.io/api/core/v1"
//	netv1 "k8s.io/api/networking/v1"
//)
//
//// VolumeType 自定义operator中公共的struct
//type VolumeType string
//
//const (
//	Mounted               VolumeType = "Mounted"
//	PersistentVolumeClaim VolumeType = "PersistentVolumeClaim"
//	HostPath              VolumeType = "HostPath"
//	EmptyDir              VolumeType = "EmptyDir"
//)
//
//type Label struct {
//	Name  string `json:"name,omitempty"`
//	Value string `json:"value,omitempty"`
//}
//
//type Image struct {
//	HarborDomain    string            `json:"harborDomain,omitempty"`
//	Repository      string            `json:"repository,omitempty"`
//	Tag             string            `json:"tag,omitempty"`
//	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
//}
//
//type Resource struct {
//	Share   bool  `json:"share,omitempty"`
//	Request int64 `json:"request,omitempty"`
//	Limit   int64 `json:"limit,omitempty"`
//}
//
//type Parameter struct {
//	ID    string `json:"id,omitempty"`
//	Value string `json:"value,omitempty"`
//}
//
//type Environment struct {
//	Name  string `json:"name,omitempty"`
//	Value string `json:"value,omitempty"`
//}
//
//type Volume struct {
//	ID               string                            `json:"id,omitempty"`
//	Type             VolumeType                        `json:"type,omitempty"`
//	Mount            string                            `json:"mount,omitempty"`
//	AccessMode       corev1.PersistentVolumeAccessMode `json:"accessMode,omitempty"`
//	Capacity         string                            `json:"capacity,omitempty"`
//	StorageClassName string                            `json:"storageClassName,omitempty"`
//	Location         string                            `json:"location,omitempty"`
//	LocationType     corev1.HostPathType               `json:"locationType,omitempty"`
//}
//
//type Config struct {
//	ID      string `json:"id,omitempty"`
//	Mount   string `json:"mount,omitempty"`
//	Content string `json:"content,omitempty"`
//}
//
//type Log struct {
//	ID        string `json:"id,omitempty"`
//	Directory string `json:"directory,omitempty"`
//	Pattern   string `json:"pattern,omitempty"`
//}
//
//type Port struct {
//	LoadBalancer    `json:",inline"`
//	ID              string         `json:"id,omitempty"`
//	Protocol        string         `json:"protocol,omitempty"`
//	ContainerPort   int32          `json:"containerPort,omitempty"`
//	ServerPort      int32          `json:"serverPort,omitempty"`
//	LoadBalancePort int32          `json:"loadBalancePort,omitempty"`
//	NodePort        int32          `json:"nodePort,omitempty"`
//	Ingress         bool           `json:"ingress,omitempty"`
//	Host            string         `json:"host,omitempty"`
//	Path            string         `json:"path,omitempty"`
//	PathType        netv1.PathType `json:"pathType,omitempty"`
//}
//
//type Cookie struct {
//	Name string `json:"name,omitempty"`
//	Path string `json:"path,omitempty"`
//	TTL  int32  `json:"ttl,omitempty"`
//}
//
//type LoadBalancer struct {
//	LoadBalancer string  `json:"loadBalancer,omitempty"`
//	Simple       string  `json:"simple,omitempty"`
//	Hash         string  `json:"hash,omitempty"`
//	Header       string  `json:"header,omitempty"`
//	Parameter    string  `json:"parameter,omitempty"`
//	Cookie       *Cookie `json:"cookie,omitempty"`
//}
