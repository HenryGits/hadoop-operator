//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	appsv1 "k8s.io/api/apps/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Hadoop) DeepCopyInto(out *Hadoop) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Hadoop.
func (in *Hadoop) DeepCopy() *Hadoop {
	if in == nil {
		return nil
	}
	out := new(Hadoop)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Hadoop) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HadoopList) DeepCopyInto(out *HadoopList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Hadoop, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HadoopList.
func (in *HadoopList) DeepCopy() *HadoopList {
	if in == nil {
		return nil
	}
	out := new(HadoopList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HadoopList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HadoopSpec) DeepCopyInto(out *HadoopSpec) {
	*out = *in
	if in.Hdfs != nil {
		in, out := &in.Hdfs, &out.Hdfs
		*out = new(Hdfs)
		(*in).DeepCopyInto(*out)
	}
	if in.WebHdfs != nil {
		in, out := &in.WebHdfs, &out.WebHdfs
		*out = new(WebHdfs)
		**out = **in
	}
	if in.ResourceManager != nil {
		in, out := &in.ResourceManager, &out.ResourceManager
		*out = new(ResourceManager)
		(*in).DeepCopyInto(*out)
	}
	if in.NodeManager != nil {
		in, out := &in.NodeManager, &out.NodeManager
		*out = new(NodeManager)
		(*in).DeepCopyInto(*out)
	}
	if in.HistoryServer != nil {
		in, out := &in.HistoryServer, &out.HistoryServer
		*out = new(HistoryServer)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HadoopSpec.
func (in *HadoopSpec) DeepCopy() *HadoopSpec {
	if in == nil {
		return nil
	}
	out := new(HadoopSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HadoopStatus) DeepCopyInto(out *HadoopStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HadoopStatus.
func (in *HadoopStatus) DeepCopy() *HadoopStatus {
	if in == nil {
		return nil
	}
	out := new(HadoopStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Hdfs) DeepCopyInto(out *Hdfs) {
	*out = *in
	if in.DaemonSet != nil {
		in, out := &in.DaemonSet, &out.DaemonSet
		*out = new(appsv1.DaemonSet)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Hdfs.
func (in *Hdfs) DeepCopy() *Hdfs {
	if in == nil {
		return nil
	}
	out := new(Hdfs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HistoryServer) DeepCopyInto(out *HistoryServer) {
	*out = *in
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = new(ResourceRequirements)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HistoryServer.
func (in *HistoryServer) DeepCopy() *HistoryServer {
	if in == nil {
		return nil
	}
	out := new(HistoryServer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Image) DeepCopyInto(out *Image) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Image.
func (in *Image) DeepCopy() *Image {
	if in == nil {
		return nil
	}
	out := new(Image)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeManager) DeepCopyInto(out *NodeManager) {
	*out = *in
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
	if in.PdbMinAvailable != nil {
		in, out := &in.PdbMinAvailable, &out.PdbMinAvailable
		*out = new(int32)
		**out = **in
	}
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = new(ResourceRequirements)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeManager.
func (in *NodeManager) DeepCopy() *NodeManager {
	if in == nil {
		return nil
	}
	out := new(NodeManager)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceList) DeepCopyInto(out *ResourceList) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceList.
func (in *ResourceList) DeepCopy() *ResourceList {
	if in == nil {
		return nil
	}
	out := new(ResourceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceManager) DeepCopyInto(out *ResourceManager) {
	*out = *in
	if in.PdbMinAvailable != nil {
		in, out := &in.PdbMinAvailable, &out.PdbMinAvailable
		*out = new(int32)
		**out = **in
	}
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = new(ResourceRequirements)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceManager.
func (in *ResourceManager) DeepCopy() *ResourceManager {
	if in == nil {
		return nil
	}
	out := new(ResourceManager)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceRequirements) DeepCopyInto(out *ResourceRequirements) {
	*out = *in
	out.Limits = in.Limits
	out.Requests = in.Requests
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceRequirements.
func (in *ResourceRequirements) DeepCopy() *ResourceRequirements {
	if in == nil {
		return nil
	}
	out := new(ResourceRequirements)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WebHdfs) DeepCopyInto(out *WebHdfs) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WebHdfs.
func (in *WebHdfs) DeepCopy() *WebHdfs {
	if in == nil {
		return nil
	}
	out := new(WebHdfs)
	in.DeepCopyInto(out)
	return out
}
