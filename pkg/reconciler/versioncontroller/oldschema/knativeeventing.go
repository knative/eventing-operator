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
package oldschema

import (
	apiextensionv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	eventingv1alpha1 "knative.dev/eventing-operator/pkg/apis/eventing/v1alpha1"
)

// The following versions of Knative eventing operators use this schema for the custom resource:
// * 0.10.0
type Eventing struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   eventingv1alpha1.EventingSpec   `json:"spec,omitempty"`
	Status eventingv1alpha1.EventingStatus `json:"status,omitempty"`
}

const (
	GroupName = "operator.knative.dev"
	SchemaVersion = "v1alpha1"
	Kind = "Eventing"
	CRDPlural = "eventings"
	CRDListKind = "EventingList"
	Singular = "eventing"
	FullCRDName string = CRDPlural + "." + GroupName
)

func GetOldGroupVersionResource(version string) schema.GroupVersionResource {
	// GroupVersionResource should be returned based on the old version.
	plural, _ := meta.UnsafeGuessKindToResource(schema.GroupVersionKind{Group: GroupName,
		Version: SchemaVersion, Kind: Kind})
	return plural
}

func InstallOldCRD(version string, clientset apiextension.Interface) error {
	crd := &apiextensionv1beta1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{Name: FullCRDName},
		Spec: apiextensionv1beta1.CustomResourceDefinitionSpec{
			Group:   GroupName,
			Version: SchemaVersion,
			Scope:   apiextensionv1beta1.NamespaceScoped,
			Names: apiextensionv1beta1.CustomResourceDefinitionNames{
				Plural:   CRDPlural,
				Kind:     Kind,
				ListKind: CRDListKind,
				Singular: Singular,
			},
		},
	}

	_, err := clientset.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
	if err != nil && apierrors.IsAlreadyExists(err) {
		return nil
	}
	return err
}
