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
// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"context"
	"time"

	v1 "github.com/HenryGits/hadoop-operator/apis/hadoop/v1"
	scheme "github.com/HenryGits/hadoop-operator/pkg/client/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// HadoopsGetter has a method to return a HadoopInterface.
// A group's client should implement this interface.
type HadoopsGetter interface {
	Hadoops(namespace string) HadoopInterface
}

// HadoopInterface has methods to work with Hadoop resources.
type HadoopInterface interface {
	Create(ctx context.Context, hadoop *v1.Hadoop, opts metav1.CreateOptions) (*v1.Hadoop, error)
	Update(ctx context.Context, hadoop *v1.Hadoop, opts metav1.UpdateOptions) (*v1.Hadoop, error)
	UpdateStatus(ctx context.Context, hadoop *v1.Hadoop, opts metav1.UpdateOptions) (*v1.Hadoop, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.Hadoop, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.HadoopList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Hadoop, err error)
	HadoopExpansion
}

// hadoops implements HadoopInterface
type hadoops struct {
	client rest.Interface
	ns     string
}

// newHadoops returns a Hadoops
func newHadoops(c *HadoopV1Client, namespace string) *hadoops {
	return &hadoops{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the hadoop, and returns the corresponding hadoop object, and an error if there is any.
func (c *hadoops) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.Hadoop, err error) {
	result = &v1.Hadoop{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("hadoops").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Hadoops that match those selectors.
func (c *hadoops) List(ctx context.Context, opts metav1.ListOptions) (result *v1.HadoopList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.HadoopList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("hadoops").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested hadoops.
func (c *hadoops) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("hadoops").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a hadoop and creates it.  Returns the server's representation of the hadoop, and an error, if there is any.
func (c *hadoops) Create(ctx context.Context, hadoop *v1.Hadoop, opts metav1.CreateOptions) (result *v1.Hadoop, err error) {
	result = &v1.Hadoop{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("hadoops").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(hadoop).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a hadoop and updates it. Returns the server's representation of the hadoop, and an error, if there is any.
func (c *hadoops) Update(ctx context.Context, hadoop *v1.Hadoop, opts metav1.UpdateOptions) (result *v1.Hadoop, err error) {
	result = &v1.Hadoop{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("hadoops").
		Name(hadoop.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(hadoop).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *hadoops) UpdateStatus(ctx context.Context, hadoop *v1.Hadoop, opts metav1.UpdateOptions) (result *v1.Hadoop, err error) {
	result = &v1.Hadoop{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("hadoops").
		Name(hadoop.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(hadoop).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the hadoop and deletes it. Returns an error if one occurs.
func (c *hadoops) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("hadoops").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *hadoops) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("hadoops").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched hadoop.
func (c *hadoops) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Hadoop, err error) {
	result = &v1.Hadoop{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("hadoops").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
