/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Phase string

const (
	Reconciling Phase = "Reconciling"
	Running     Phase = "Running"
	Deleting    Phase = "Deleting"
)

//+genclient
//+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Hadoop is the Schema for the hadoops API
type Hadoop struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HadoopSpec   `json:"spec,omitempty"`
	Status HadoopStatus `json:"status,omitempty"`
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HadoopSpec defines the desired state of Hadoop
type HadoopSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// ID is the id of the hadoop release
	// +optional
	ID string `json:"id"`
	// Title is the title of the hadoop release
	// +optional
	Title string `json:"title,omitempty"`
	// Describe is the description of the hadoop release
	// +optional
	Describe string `json:"describe,omitempty"`
	// Image is the base hadoop image to use for all components.
	// +optional
	// +kubebuilder:default:={repository: "jasonchrion/hadoop", tag: "3.2.2-nolib", pullPolicy: "IfNotPresent"}
	Image Image `json:"image,omitempty"`
	// HadoopVersion is the version of the hadoop libraries being used in the image.
	// +optional
	// +kubebuilder:default="3.2.2"
	HadoopVersion string `json:"hadoopVersion,omitempty"`
	// AntiAffinity Select antiAffinity as either hard or soft, default is soft
	// +optional
	// +kubebuilder:validation:Enum={soft,hard}
	// +kubebuilder:default="soft"
	AntiAffinity AntiAffinity `json:"antiAffinity,omitempty"`
	// Hdfs is hadoop hdfs components include NameNode, DataNode, WebHdfs
	// +optional
	Hdfs *Hdfs `json:"hdfs,omitempty"`
	// Yarn is hadoop yarn components include ResourceManager, NodeManager
	// +optional
	Yarn *Yarn `json:"yarn,omitempty"`
	// Persistence is persistent volume for Hadoop components
	// +optional
	Persistence *Persistence `json:"persistence,omitempty"`
	// PostInstallCommands is the hdfs commands to be run after the installation
	// +optional
	// +kubebuilder:default={}
	PostInstallCommands []string `json:"postInstallCommands,omitempty"`
}

// HadoopStatus defines the observed state of Hadoop
type HadoopStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Phase Phase `json:"phase,omitempty"`
}

//+kubebuilder:object:root=true
//+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HadoopList contains a list of Hadoop
type HadoopList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Hadoop `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Hadoop{}, &HadoopList{})
}

type Image struct {
	// Repository is Hadoop image repository
	// +optional
	// +kubebuilder:default="jasonchrion/hadoop"
	Repository string `json:"repository,omitempty"`
	// Tag is the Hadoop image tag
	// +optional
	// +kubebuilder:default="3.2.2-nolib"
	Tag string `json:"tag,omitempty"`
	// PullPolicy is the pull policy for the images
	// +optional
	// +kubebuilder:validation:Enum={Always,Never,IfNotPresent}
	// +kubebuilder:default="IfNotPresent"
	PullPolicy corev1.PullPolicy `json:"pullPolicy,omitempty"`
}

type AntiAffinity string

const (
	Soft AntiAffinity = "soft"
	Hard AntiAffinity = "hard"
)

type Hdfs struct {
	NameNode *NameNode `json:"nameNode,omitempty"`
	DataNode *DataNode `json:"dataNode,omitempty"`
	WebHdfs  *WebHdfs  `json:"webHdfs,omitempty"`
}

type NameNode struct {
	// PdbMinAvailable is the minimum available number of PodDisruptionBudget for Hadoop component
	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=1
	PdbMinAvailable *int32 `json:"pdbMinAvailable,omitempty"`
	// Resources is the CPU and memory resource (requests and limits) allocated to each Hadoop component pod.
	// This should be tuned to fit your workload.
	// +optional
	// +kubebuilder:default:={requests: {cpu: "100m", memory: "256Mi"}, limits: {cpu: "500m", memory: "1Gi"}}
	Resources *ResourceRequirements `json:"resources,omitempty"`
}

type DataNode struct {
	// Replicas is the pod number of Hadoop component.
	// +optional
	// +kubebuilder:default:=3
	Replicas *int32 `json:"replicas,omitempty"`
	// PdbMinAvailable is the minimum available number of PodDisruptionBudget for Hadoop component
	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=1
	PdbMinAvailable *int32 `json:"pdbMinAvailable,omitempty"`
	// Resources is the CPU and memory resource (requests and limits) allocated to each Hadoop component pod.
	// This should be tuned to fit your workload.
	// +optional
	// +kubebuilder:default:={requests: {cpu: "100m", memory: "256Mi"}, limits: {cpu: "500m", memory: "1Gi"}}
	Resources *ResourceRequirements `json:"resources,omitempty"`
}

type WebHdfs struct {
	// Enabled is whether to enable WebHDFS REST API or not
	// +optional
	// +kubebuilder:validation:Enum=true;false
	// +kubebuilder:default=true
	Enabled bool `json:"enabled,omitempty"`
}

type Yarn struct {
	ResourceManager *ResourceManager `json:"resourceManager,omitempty"`
	NodeManager     *NodeManager     `json:"nodeManager,omitempty"`
}

type ResourceManager struct {
	// PdbMinAvailable is the minimum available number of PodDisruptionBudget for Hadoop component
	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=1
	PdbMinAvailable *int32 `json:"pdbMinAvailable,omitempty"`
	// Resources is the CPU and memory resource (requests and limits) allocated to each Hadoop component pod.
	// This should be tuned to fit your workload.
	// +optional
	// +kubebuilder:default:={requests: {cpu: "100m", memory: "256Mi"}, limits: {cpu: "500m", memory: "1Gi"}}
	Resources *ResourceRequirements `json:"resources,omitempty"`
}

type NodeManager struct {
	// Replicas is the pod number of Hadoop component.
	// +optional
	// +kubebuilder:default:=3
	Replicas *int32 `json:"replicas,omitempty"`
	// PdbMinAvailable is the minimum available number of PodDisruptionBudget for Hadoop component
	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=1
	PdbMinAvailable *int32 `json:"pdbMinAvailable,omitempty"`
	// Resources is the CPU and memory resource (requests and limits) allocated to each Hadoop component pod.
	// This should be tuned to fit your workload.
	// +optional
	// +kubebuilder:default:={requests: {cpu: "100m", memory: "256Mi"}, limits: {cpu: "500m", memory: "1Gi"}}
	Resources *ResourceRequirements `json:"resources,omitempty"`
	// ParallelCreate is whether to create all nodeManager statefulset pods in parallel or not (K8S 1.7+)
	// +optional
	// +kubebuilder:validation:Enum={true,false}
	// +kubebuilder:default=true
	ParallelCreate bool `json:"parallelCreate,omitempty"`
}

type Persistence struct {
	NameNode NameNodePersistence `json:"nameNode,omitempty"`
	DataNode DataNodePersistence `json:"dataNode,omitempty"`
}

type NameNodePersistence struct {
	// Enabled is whether to enable Hadoop component persistence or not
	// +optional
	// +kubebuilder:validation:Enum={true,false}
	// +kubebuilder:default=false
	Enabled bool `json:"enabled,omitempty"`
	// StorageClass is the name of the StorageClass required by the claim.
	// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#class-1
	// +optional
	// +kubebuilder:default="-"
	StorageClass *string `json:"storageClass,omitempty"`
	// AccessMode contains the desired access modes the volume should have.
	// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#access-modes-1
	// +optional
	// +kubebuilder:validation:Enum={ReadWriteOnce,ReadOnlyMany,ReadWriteMany}
	// +kubebuilder:default="ReadWriteOnce"
	AccessMode corev1.PersistentVolumeAccessMode `json:"accessMode,omitempty"`
	// Size is the size of the volume
	// +optional
	// +kubebuilder:default:="50Gi"
	Size string `json:"size,omitempty"`
}

type DataNodePersistence struct {
	// Enabled is whether to enable Hadoop component persistence or not
	// +optional
	// +kubebuilder:validation:Enum={true,false}
	// +kubebuilder:default=false
	Enabled bool `json:"enabled,omitempty"`
	// StorageClass is the name of the StorageClass required by the claim.
	// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#class-1
	// +optional
	// +kubebuilder:default="-"
	StorageClass *string `json:"storageClass,omitempty"`
	// AccessMode contains the desired access modes the volume should have.
	// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#access-modes-1
	// +optional
	// +kubebuilder:validation:Enum={ReadWriteOnce,ReadOnlyMany,ReadWriteMany}
	// +kubebuilder:default="ReadWriteOnce"
	AccessMode corev1.PersistentVolumeAccessMode `json:"accessMode,omitempty"`
	// Size is the size of the volume
	// +optional
	// +kubebuilder:default:="200Gi"
	Size string `json:"size,omitempty"`
}

// ResourceRequirements describes the compute resource requirements.
type ResourceRequirements struct {
	// Limits describes the maximum amount of compute resources allowed.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	// +optional
	Limits ResourceList `json:"limits,omitempty"`
	// Requests describes the minimum amount of compute resources required.
	// If Requests is omitted for a container, it defaults to Limits if that is explicitly specified,
	// otherwise to an implementation-defined value.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	// +optional
	Requests ResourceList `json:"requests,omitempty"`
}

type ResourceList struct {
	Cpu    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
}
