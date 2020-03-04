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
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/apis"
	duckv1beta1 "knative.dev/pkg/apis/duck/v1beta1"
)

var ignoreAllButTypeAndStatus = cmpopts.IgnoreFields(
	apis.Condition{},
	"LastTransitionTime", "Message", "Reason", "Severity")

func TestKnativeEventingGroupVersionKind(t *testing.T) {
	r := &KnativeEventing{}
	want := schema.GroupVersionKind{
		Group:   GroupName,
		Version: SchemaVersion,
		Kind:    Kind,
	}
	if got := r.GroupVersionKind(); got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func TestKnativeEventingInitializeConditions(t *testing.T) {
	tests := []struct {
		name string
		ke   *KnativeEventingStatus
		want *KnativeEventingStatus
	}{{
		name: "empty",
		ke:   &KnativeEventingStatus{},
		want: &KnativeEventingStatus{
			Status: duckv1beta1.Status{
				Conditions: []apis.Condition{{
					Type:   DependenciesInstalled,
					Status: corev1.ConditionUnknown,
				}, {
					Type:   DeploymentsAvailable,
					Status: corev1.ConditionUnknown,
				}, {
					Type:   InstallSucceeded,
					Status: corev1.ConditionUnknown,
				}, {
					Type:   Ready,
					Status: corev1.ConditionUnknown,
				}},
			},
		},
	},}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.ke.InitializeConditions()
			if diff := cmp.Diff(test.want, test.ke, ignoreAllButTypeAndStatus); diff != "" {
				t.Errorf("unexpected conditions (-want, +got) = %v", diff)
			}
		})
	}
}
