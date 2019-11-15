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
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis"
	duckv1beta1 "knative.dev/pkg/apis/duck/v1beta1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Eventing is the Schema for the eventings API
// +k8s:openapi-gen=true
type Eventing struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EventingSpec   `json:"spec,omitempty"`
	Status EventingStatus `json:"status,omitempty"`
}

// Registry defines image overrides of knative images.
// This affects both apps/v1.Deployment and caching.internal.knative.dev/v1alpha1.Image.
// The default value is used as a default format to override for all knative deployments.
// The override values are specific to each knative deployment.
// +k8s:openapi-gen=true
type Registry struct {
	// The default image reference template to use for all knative images.
	// It takes the form of example-registry.io/custom/path/${NAME}:custom-tag
	// ${NAME} will be replaced by the deployment container name, or caching.internal.knative.dev/v1alpha1/Image name.
	// +optional
	Default string `json:"default,omitempty"`

	// A map of a container name or image name to the full image location of the individual knative image.
	// +optional
	Override map[string]string `json:"override,omitempty"`
}

// EventingSpec defines the desired state of Eventing
// +k8s:openapi-gen=true
type EventingSpec struct {
}

// EventingStatus defines the observed state of Eventing
// +k8s:openapi-gen=true
type EventingStatus struct {
	duckv1beta1.Status `json:",inline"`

	// The version of the installed release
	// +optional
	Version string `json:"version,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EventingList contains a list of Eventing
type EventingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Eventing `json:"items"`
}

const (
	// EventingConditionReady is set when the Eventing Operator is installed, configured and ready.
	EventingConditionReady = apis.ConditionReady

	// InstallSucceeded is set when the Knative Eventing is installed.
	InstallSucceeded apis.ConditionType = "InstallSucceeded"
)
