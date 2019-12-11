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

func TestKnativeEventingStatusGetCondition(t *testing.T) {
	ke := &KnativeEventingStatus{}
	if a := ke.GetCondition(InstallSucceeded); a != nil {
		t.Errorf("empty EventingStatus returned %v when expected nil", a)
	}
	mc := &apis.Condition{
		Type:   InstallSucceeded,
		Status: corev1.ConditionTrue,
	}
	ke.MarkEventingInstalled()
	if diff := cmp.Diff(mc, ke.GetCondition(InstallSucceeded), cmpopts.IgnoreFields(apis.Condition{}, "LastTransitionTime")); diff != "" {
		t.Errorf("GetCondition refs diff (-want +got): %v", diff)
	}
}

func TestKnativeEventingStatusEventingInstalled(t *testing.T) {
	ke := &KnativeEventingStatus{}
	mc := &apis.Condition{
		Type:   InstallSucceeded,
		Status: corev1.ConditionTrue,
	}
	ke.MarkEventingInstalled()
	if diff := cmp.Diff(mc, ke.GetCondition(InstallSucceeded), cmpopts.IgnoreFields(apis.Condition{}, "LastTransitionTime")); diff != "" {
		t.Errorf("GetCondition refs diff (-want +got): %v", diff)
	}
}

func TestKnativeEventingStatusEventingFailed(t *testing.T) {
	reason := "NotReady"
	message := "Waiting on deployments"
	ke := &KnativeEventingStatus{}
	mc := &apis.Condition{
		Type:    EventingConditionReady,
		Status:  corev1.ConditionFalse,
		Reason:  reason,
		Message: message,
	}
	ke.MarkEventingFailed(reason, message)
	if diff := cmp.Diff(mc, ke.GetCondition(EventingConditionReady), cmpopts.IgnoreFields(apis.Condition{}, "LastTransitionTime")); diff != "" {
		t.Errorf("GetCondition refs diff (-want +got): %v", diff)
	}
}

func TestKnativeEventingStatusNotReady(t *testing.T) {
	reason := "NotReady"
	message := "Waiting on deployments"
	ke := &KnativeEventingStatus{}
	mc := &apis.Condition{
		Type:    EventingConditionReady,
		Status:  corev1.ConditionUnknown,
		Reason:  reason,
		Message: message,
	}
	ke.MarkEventingNotReady(reason, message)
	if diff := cmp.Diff(mc, ke.GetCondition(EventingConditionReady), cmpopts.IgnoreFields(apis.Condition{}, "LastTransitionTime")); diff != "" {
		t.Errorf("GetCondition refs diff (-want +got): %v", diff)
	}
}

func TestKnativeEventingStatusReady(t *testing.T) {
	ke := &KnativeEventingStatus{}
	mc := &apis.Condition{
		Type:   EventingConditionReady,
		Status: corev1.ConditionTrue,
	}
	ke.MarkEventingReady()
	if diff := cmp.Diff(mc, ke.GetCondition(EventingConditionReady), cmpopts.IgnoreFields(apis.Condition{}, "LastTransitionTime")); diff != "" {
		t.Errorf("GetCondition refs diff (-want +got): %v", diff)
	}
}

func TestKnativeEventingStatusIsReady(t *testing.T) {
	ke := &KnativeEventingStatus{}
	ke.MarkEventingReady()
	if diff := cmp.Diff(true, ke.IsReady()); diff != "" {
		t.Errorf("IsReady refs diff (-want +got): %v", diff)
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
					Type:   InstallSucceeded,
					Status: corev1.ConditionUnknown,
				}, {
					Type:   EventingConditionReady,
					Status: corev1.ConditionUnknown,
				}},
			},
		},
	}, {
		name: "eventingConditionNotReady",
		ke: &KnativeEventingStatus{
			Status: duckv1beta1.Status{
				Conditions: []apis.Condition{{
					Type:   EventingConditionReady,
					Status: corev1.ConditionFalse,
				}},
			},
		},
		want: &KnativeEventingStatus{
			Status: duckv1beta1.Status{
				Conditions: []apis.Condition{{
					Type:   InstallSucceeded,
					Status: corev1.ConditionUnknown,
				}, {
					Type:   EventingConditionReady,
					Status: corev1.ConditionFalse,
				}},
			},
		},
	}, {
		name: "eventingConditionReady",
		ke: &KnativeEventingStatus{
			Status: duckv1beta1.Status{
				Conditions: []apis.Condition{{
					Type:   EventingConditionReady,
					Status: corev1.ConditionTrue,
				}},
			},
		},
		want: &KnativeEventingStatus{
			Status: duckv1beta1.Status{
				Conditions: []apis.Condition{{
					Type:   InstallSucceeded,
					Status: corev1.ConditionTrue,
				}, {
					Type:   EventingConditionReady,
					Status: corev1.ConditionTrue,
				}},
			},
		},
	}, {
		name: "installSucceeded",
		ke: &KnativeEventingStatus{
			Status: duckv1beta1.Status{
				Conditions: []apis.Condition{{
					Type:   InstallSucceeded,
					Status: corev1.ConditionTrue,
				}},
			},
		},
		want: &KnativeEventingStatus{
			Status: duckv1beta1.Status{
				Conditions: []apis.Condition{{
					Type:   InstallSucceeded,
					Status: corev1.ConditionTrue,
				}, {
					Type:   EventingConditionReady,
					Status: corev1.ConditionUnknown,
				}},
			},
		},
	}, {
		name: "installNotSucceeded",
		ke: &KnativeEventingStatus{
			Status: duckv1beta1.Status{
				Conditions: []apis.Condition{{
					Type:   InstallSucceeded,
					Status: corev1.ConditionFalse,
				}},
			},
		},
		want: &KnativeEventingStatus{
			Status: duckv1beta1.Status{
				Conditions: []apis.Condition{{
					Type:   InstallSucceeded,
					Status: corev1.ConditionFalse,
				}, {
					Type:   EventingConditionReady,
					Status: corev1.ConditionUnknown,
				}},
			},
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.ke.InitializeConditions()
			if diff := cmp.Diff(test.want, test.ke, ignoreAllButTypeAndStatus); diff != "" {
				t.Errorf("unexpected conditions (-want, +got) = %v", diff)
			}
		})
	}
}
