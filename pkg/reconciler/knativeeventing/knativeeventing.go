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

package knativeeventing

import (
	"context"

	mf "github.com/jcrossley3/manifestival"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/tools/cache"

	"knative.dev/pkg/controller"
	listers "knative.dev/eventing-operator/pkg/client/listers/eventing/v1alpha1"
	"knative.dev/eventing-operator/pkg/reconciler"
)

// Reconciler implements controller.Reconciler for Knativeeventing resources.
type Reconciler struct {
	*reconciler.Base
	// Listers index properties about resources
	knativeEventingLister listers.EventingLister
	config                mf.Manifest
	eventings             sets.String
}

// Check that our Reconciler implements controller.Reconciler
var _ controller.Reconciler = (*Reconciler)(nil)

// Reconcile compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the Knativeeventing resource
// with the current status of the resource.
func (r *Reconciler) Reconcile(ctx context.Context, key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		r.Logger.Errorf("invalid resource key: %s", key)
		return nil
	}
	// Get the KnativeEventing resource with this namespace/name.
	_, err = r.knativeEventingLister.Eventings(namespace).Get(name)
	if apierrs.IsNotFound(err) {
		// The resource was deleted
		r.eventings.Delete(key)
		if r.eventings.Len() == 0 {
			r.config.DeleteAll(&metav1.DeleteOptions{})
		}
		return nil

	} else if err != nil {
		r.Logger.Error(err, "Error getting KnativeEventings")
		return err
	}
	// Keep track of the number of KnativeEventings in the cluster
	r.eventings.Insert(key)
	return nil
}
