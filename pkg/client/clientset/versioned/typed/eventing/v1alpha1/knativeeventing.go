/*
Copyright 2019 The Knative Authors

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

package v1alpha1

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
	v1alpha1 "knative.dev/eventing-operator/pkg/apis/eventing/v1alpha1"
	scheme "knative.dev/eventing-operator/pkg/client/clientset/versioned/scheme"
)

// KnativeEventingsGetter has a method to return a KnativeEventingInterface.
// A group's client should implement this interface.
type KnativeEventingsGetter interface {
	KnativeEventings(namespace string) KnativeEventingInterface
}

// KnativeEventingInterface has methods to work with KnativeEventing resources.
type KnativeEventingInterface interface {
	Create(*v1alpha1.KnativeEventing) (*v1alpha1.KnativeEventing, error)
	Update(*v1alpha1.KnativeEventing) (*v1alpha1.KnativeEventing, error)
	UpdateStatus(*v1alpha1.KnativeEventing) (*v1alpha1.KnativeEventing, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.KnativeEventing, error)
	List(opts v1.ListOptions) (*v1alpha1.KnativeEventingList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.KnativeEventing, err error)
	KnativeEventingExpansion
}

// knativeEventings implements KnativeEventingInterface
type knativeEventings struct {
	client rest.Interface
	ns     string
}

// newKnativeEventings returns a KnativeEventings
func newKnativeEventings(c *EventingV1alpha1Client, namespace string) *knativeEventings {
	return &knativeEventings{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the knativeEventing, and returns the corresponding knativeEventing object, and an error if there is any.
func (c *knativeEventings) Get(name string, options v1.GetOptions) (result *v1alpha1.KnativeEventing, err error) {
	result = &v1alpha1.KnativeEventing{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("knativeeventings").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of KnativeEventings that match those selectors.
func (c *knativeEventings) List(opts v1.ListOptions) (result *v1alpha1.KnativeEventingList, err error) {
	result = &v1alpha1.KnativeEventingList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("knativeeventings").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested knativeEventings.
func (c *knativeEventings) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("knativeeventings").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a knativeEventing and creates it.  Returns the server's representation of the knativeEventing, and an error, if there is any.
func (c *knativeEventings) Create(knativeEventing *v1alpha1.KnativeEventing) (result *v1alpha1.KnativeEventing, err error) {
	result = &v1alpha1.KnativeEventing{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("knativeeventings").
		Body(knativeEventing).
		Do().
		Into(result)
	return
}

// Update takes the representation of a knativeEventing and updates it. Returns the server's representation of the knativeEventing, and an error, if there is any.
func (c *knativeEventings) Update(knativeEventing *v1alpha1.KnativeEventing) (result *v1alpha1.KnativeEventing, err error) {
	result = &v1alpha1.KnativeEventing{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("knativeeventings").
		Name(knativeEventing.Name).
		Body(knativeEventing).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *knativeEventings) UpdateStatus(knativeEventing *v1alpha1.KnativeEventing) (result *v1alpha1.KnativeEventing, err error) {
	result = &v1alpha1.KnativeEventing{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("knativeeventings").
		Name(knativeEventing.Name).
		SubResource("status").
		Body(knativeEventing).
		Do().
		Into(result)
	return
}

// Delete takes name of the knativeEventing and deletes it. Returns an error if one occurs.
func (c *knativeEventings) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("knativeeventings").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *knativeEventings) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("knativeeventings").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched knativeEventing.
func (c *knativeEventings) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.KnativeEventing, err error) {
	result = &v1alpha1.KnativeEventing{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("knativeeventings").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
