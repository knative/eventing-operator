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

package upgrade

import (
	"errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes/scheme"
	eventingv1alpha1 "knative.dev/eventing-operator/pkg/apis/eventing/v1alpha1"
	"knative.dev/eventing-operator/pkg/reconciler/versioncontroller/oldschema"
)

func UpgradeCR(sourceVersion, targetVersion string, uobject *unstructured.Unstructured, target *eventingv1alpha1.KnativeEventing) error {
	switch sourceVersion {
	case "0.10.0":
		// Convert the old CR based on the information in the spec. If the spec is empty, no need to do anything.
		err := Convert10To11(uobject, target)
		if err != nil {
			return err
		}
	case targetVersion:
		return nil
	default:
		err := errors.New("The old version of CR does not support upgrade.")
		return err
	}
	return nil
}

func Convert10To11(uobject *unstructured.Unstructured, target *eventingv1alpha1.KnativeEventing) error {
	// Keep the old schema, and convert the unstructured into the old version of CR.
	oldEventing := &oldschema.Eventing{}
	err := scheme.Scheme.Convert(uobject, oldEventing, nil)
	if err != nil {
		return err
	}

	// Verify whether the spec is empty for CR in 0.10.0.
	if (oldEventing.Spec != oldschema.EventingSpec{})  {
		err := errors.New("The old CR is in bad format, since the spec is not empty.")
		return err
	}
	return nil
}
