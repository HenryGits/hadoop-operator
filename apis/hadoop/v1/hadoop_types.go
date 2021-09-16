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
	appsv1 "k8s.io/api/apps/v1"
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

	// +optional
	Hdfs    *Hdfs    `json:"hdfs,omitempty"`
	WebHdfs *WebHdfs `json:"webHdfs,omitempty"`
	// Yarn is hadoop yarn components include ResourceManager, NodeManager
	// +optional
	ResourceManager *ResourceManager `json:"resourceManager,omitempty"`
	NodeManager     *NodeManager     `json:"nodeManager,omitempty"`
	HistoryServer   *HistoryServer   `json:"historyServer,omitempty"`
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
	// +kubebuilder:default="hadoop"
	Repository string `json:"repository,omitempty"`
	// Tag is the Hadoop image tag
	// +optional
	// +kubebuilder:default="v3.3.1"
	Tag string `json:"tag,omitempty"`
	// PullPolicy is the pull policy for the images
	// +optional
	// +kubebuilder:validation:Enum={Always,Never,IfNotPresent}
	// +kubebuilder:default="IfNotPresent"
	PullPolicy corev1.PullPolicy `json:"pullPolicy,omitempty"`
}

// Hdfs  Hadoop nn & dn
type Hdfs struct {
	// +optional
	// +kubebuilder:validation:Enum={NameNode,DataNode}
	Type string `json:"type,omitempty"`
	// DaemonSet specifies the hadoop should be deployed as a DaemonSet, and allows providing its spec.
	// Cannot be used along with `deployment`. If both are absent a default for the Type is used.
	// +optional
	DaemonSet DaemonSetSpec `json:"daemonSet,omitempty"`
}

type DaemonSetSpec struct {
	// +optional
	Template corev1.PodTemplateSpec `json:"template,omitempty"`
	// +optional
	UpdateStrategy appsv1.DaemonSetUpdateStrategy `json:"updateStrategy,omitempty"`
}

type WebHdfs struct {
	// Enabled is whether to enable WebHDFS REST API or not
	// +optional
	// +kubebuilder:validation:Enum=true;false
	// +kubebuilder:default=true
	Enabled bool `json:"enabled,omitempty"`
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
	// +kubebuilder:default:=1
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

type HistoryServer struct {
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
