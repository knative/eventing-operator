/*
Copyright 2020 The Knative Authors

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

package knativeeventing

import (
	"context"
	listersObsolete "knative.dev/eventing-operator/pkg/client/listers/eventing/v1alpha1"
	"knative.dev/pkg/controller"
	"knative.dev/eventing-operator/pkg/reconciler"
)

// ReconcilerObsolete implements controller.Reconciler for Eventing resources.
type ReconcilerObsolete struct {
	*reconciler.Base
	// Listers index properties about resources
	knativeEventingObsoleteLister listersObsolete.EventingLister
}

// Check that our Reconciler implements controller.Reconciler
var _ controller.Reconciler = (*ReconcilerObsolete)(nil)

// Reconcile compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the Knativeeventing resource
// with the current status of the resource.
func (r *ReconcilerObsolete) Reconcile(ctx context.Context, key string) error {
	r.Logger.Info("Reconcile is OK")
	return nil
}
