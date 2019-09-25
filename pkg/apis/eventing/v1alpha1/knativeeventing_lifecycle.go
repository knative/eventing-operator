/*
Copyright 2019 The Knative Authors.

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
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/apis"
	duckv1beta1 "knative.dev/pkg/apis/duck/v1beta1"
)

var eventingCondSet = apis.NewLivingConditionSet(
	EventingConditionReady,
	InstallSucceeded,
)

// GetGroupVersionKind returns SchemeGroupVersion of an Ingress
func (e *Eventing) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind("Eventing")
}

// GetCondition returns the current condition of a given condition type
func (es *EventingStatus) GetCondition(t apis.ConditionType) *apis.Condition {
	return eventingCondSet.Manage(es).GetCondition(t)
}

// InitializeConditions initializes conditions of an IngressStatus
func (es *EventingStatus) InitializeConditions() {
	eventingCondSet.Manage(es).InitializeConditions()
}

// MarkEventingInstalled set InstallSucceeded in EventingStatus as true
func (es *EventingStatus) MarkEventingInstalled() {
	eventingCondSet.Manage(es).MarkTrue(InstallSucceeded)
}

// IsReady looks at the conditions and if the Status has a condition
// EventingConditionReady returns true if ConditionStatus is True
func (es *EventingStatus) IsReady() bool {
	return eventingCondSet.Manage(es).IsHappy()
}

// MarkEventingReady marks the Eventing status as ready
func (es *EventingStatus) MarkEventingReady() {
	eventingCondSet.Manage(es).MarkTrue(EventingConditionReady)
}

// MarkEventingNotReady marks the Eventing status as ready == Unknown
func (es *EventingStatus) MarkEventingNotReady(reason, message string) {
	eventingCondSet.Manage(es).MarkUnknown(EventingConditionReady, reason, message)
}

// MarkEventingFailed marks the Eventing status as failed
func (es *EventingStatus) MarkEventingFailed(reason, message string) {
	eventingCondSet.Manage(es).MarkFalse(EventingConditionReady, reason, message)
}

func (es *EventingStatus) duck() *duckv1beta1.Status {
	return &es.Status
}
