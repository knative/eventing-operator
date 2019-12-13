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

package downgrade

import (
	"errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	eventingv1alpha1 "knative.dev/eventing-operator/pkg/apis/eventing/v1alpha1"
	"knative.dev/eventing-operator/pkg/reconciler/versioncontroller/oldschema"
)

func DowngradeCR(targetVersion string, source *eventingv1alpha1.Eventing,
	target *unstructured.Unstructured, dynamicClientSet dynamic.Interface) (schema.GroupVersionResource, error) {
	switch targetVersion {
	case "0.10.0":
		// Convert the old CR based on the information in the spec. If the spec is empty, no need to do anything.
		//plural, _ := meta.UnsafeGuessKindToResource(schema.GroupVersionKind{Group: "operator.knative.dev",
		//	Version: "v1alpha1", Kind: "Eventing"})
		//targetDyn, err := dynamicClientSet.Resource(plural).Namespace(source.Namespace).Get("knative-eventing",
		//	metav1.GetOptions{})
		//if apierrs.IsNotFound(err) {
		//	// The CR in the new version does not exist, so create one.
		//	eventing = &eventingv1alpha1.Eventing{
		//		ObjectMeta: metav1.ObjectMeta{
		//			Name:      "knative-eventing",
		//			Namespace: keVerController.Namespace,
		//		},
		//	}
		//} else if e != nil {
		//	r.Logger.Error(e, "Error getting the CR of Knative Eventing in the new version")
		//	return err
		//}
		return Convert11To10(source, target)
	case targetVersion:
		return schema.GroupVersionResource{}, nil
	default:
		err := errors.New("The old version of CR does not support downgrade.")
		return schema.GroupVersionResource{}, err
	}
	return schema.GroupVersionResource{}, nil
}

func Convert11To10(source *eventingv1alpha1.Eventing, target *unstructured.Unstructured) (schema.GroupVersionResource, error) {

	oldEventing := &oldschema.Eventing{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "knative-eventing",
			Namespace: source.Namespace,
		},
	}

	// Verify whether the spec is empty for CR in 0.11.0.
	if (source.Spec != eventingv1alpha1.EventingSpec{})  {
		err := errors.New("The old CR is in bad format, since the spec is not empty.")
		return schema.GroupVersionResource{}, err
	}

	// Convert oldEventing into the unstructured.Unstructured.
	oldEventing.Kind = "Eventing"
	oldEventing.APIVersion = "operator.knative.dev/v1alpha1"
	oldEventing.Status.Version = "0.10.0"
	content, err := runtime.DefaultUnstructuredConverter.ToUnstructured(oldEventing)
	if err != nil {
		err := errors.New("The old CR can not be converted into unstructured.Unstructured.")
		return schema.GroupVersionResource{}, err
	}
	target.SetUnstructuredContent(content)
	plural, _ := meta.UnsafeGuessKindToResource(schema.GroupVersionKind{Group: "operator.knative.dev",
		Version: "v1alpha1", Kind: "Eventing"})

	return plural, nil
}
