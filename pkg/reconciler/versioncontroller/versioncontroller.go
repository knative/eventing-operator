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

package versioncontroller

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"knative.dev/eventing-operator/pkg/reconciler/versioncontroller/oldschema"
	"knative.dev/eventing-operator/pkg/reconciler/versioncontroller/upgrade"
	"knative.dev/eventing-operator/pkg/reconciler/versioncontroller/downgrade"
	"knative.dev/eventing-operator/version"
	"reflect"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/equality"

	corev1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"

	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	eventingv1alpha1 "knative.dev/eventing-operator/pkg/apis/eventing/v1alpha1"
	listers "knative.dev/eventing-operator/pkg/client/listers/eventing/v1alpha1"
	"knative.dev/eventing-operator/pkg/reconciler"
	"knative.dev/eventing-operator/pkg/reconciler/knativeeventing/common"
	"knative.dev/pkg/controller"
)

var (
	platform common.Platforms
)

// Reconciler implements controller.Reconciler for Knativeeventing resources.
type Reconciler struct {
	*reconciler.Base
	versionControllerLister listers.KEVersionControllerLister
	kubeClient              apiextension.Interface
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
	original, err := r.versionControllerLister.KEVersionControllers(namespace).Get(name)
	if apierrs.IsNotFound(err) {
		// The resource was deleted
		return nil

	} else if err != nil {
		r.Logger.Error(err, "Error getting the CR of version controller for Knative Eventing")
		return err
	}

	// Don't modify the informers copy.
	keVerController := original.DeepCopy()

	// Reconcile this copy of the Eventing resource and then write back any status
	// updates regardless of whether the reconciliation errored out.
	reconcileErr := r.reconcile(ctx, keVerController)
	if equality.Semantic.DeepEqual(original.Status, keVerController.Status) {
		// If we didn't change anything then don't call updateStatus.
		// This is important because the copy we loaded from the informer's
		// cache may be stale and we don't want to overwrite a prior update
		// to status with this stale state.
	} else if _, err = r.updateStatus(keVerController); err != nil {
		r.Logger.Warnw("Failed to update Eventing status", zap.Error(err))
		r.Recorder.Eventf(keVerController, corev1.EventTypeWarning, "UpdateFailed",
			"Failed to update status for Eventing %q: %v", keVerController.Name, err)
		return err
	}
	if reconcileErr != nil {
		r.Recorder.Event(keVerController, corev1.EventTypeWarning, "InternalError", reconcileErr.Error())
		return reconcileErr
	}
	return nil
}

func (r *Reconciler) reconcile(ctx context.Context, keVerController *eventingv1alpha1.KEVersionController) error {
	r.Logger.Info("Start to reconcile the version controller")
	// Check the target version to see if it is upgrade or downgrade.
	if keVerController.Spec.TargetVersion == version.Version || keVerController.Spec.TargetVersion == "" {
		// 1. Get the old CR by the DynamicClientSet.
		plural := oldschema.GetOldGroupVersionResource(keVerController.Spec.SourceVersion)
		object, err := r.DynamicClientSet.Resource(plural).Namespace(keVerController.Namespace).Get("knative-eventing",
			metav1.GetOptions{})

		if apierrs.IsNotFound(err) {
			// The CR in the old version does not exist.
			r.Logger.Infof("The CR in the version of %s does not exist.", keVerController.Spec.SourceVersion)
			return nil

		} else if err != nil {
			return fmt.Errorf("Error getting the CR of Knative Eventing in the old version")
		}

		// 2. Get the version of the new CR, and use it as the target CR.
		var eventing *eventingv1alpha1.KnativeEventing
		var newCRExist = true
		eventing, e := r.KnativeEventingClientSet.OperatorV1alpha1().KnativeEventings(keVerController.Namespace).Get("knative-eventing", metav1.GetOptions{})
		if apierrs.IsNotFound(e) {
			// The CR in the new version does not exist, so create one.
			r.Logger.Info("Prepare to create the new CR in the current version ", version.Version)
			eventing = &eventingv1alpha1.KnativeEventing{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "knative-eventing",
					Namespace: keVerController.Namespace,
				},
			}
			newCRExist = false
		} else if e != nil {
			return fmt.Errorf("Error getting the CR of Knative Eventing in the new version")
		}

		// 3. Upgrade the CR.
		r.Logger.Info("Upgrade the CR into the current version ", version.Version)
		err = upgrade.UpgradeCR(keVerController.Spec.SourceVersion, version.Version,
			object, eventing)
		if err != nil {
			return err
		}
		// 4. Create or update the new CR after the conversion for upgrade.
		r.Logger.Info("Update or create the CR in the new version ", version.Version)
		if newCRExist {
			_, e = r.KnativeEventingClientSet.OperatorV1alpha1().KnativeEventings(keVerController.Namespace).UpdateStatus(eventing)
			if e != nil {
				return fmt.Errorf("Failed to update the version of the CR to the new version")
			} else {
				_, e = r.KnativeEventingClientSet.OperatorV1alpha1().KnativeEventings(keVerController.Namespace).Update(eventing)
				if e != nil {
					return fmt.Errorf("Failed to update the CR to the new version")
				}
			}
		} else {
			_, e = r.KnativeEventingClientSet.OperatorV1alpha1().KnativeEventings(keVerController.Namespace).Create(eventing)
			if e != nil {
				return fmt.Errorf("Failed to creating the CR of Knative Eventing in the new version")
			}
		}
		r.Logger.Infof("Successfully upgraded the CR to the new version %s.", version.Version)
	} else if (keVerController.Spec.SourceVersion == version.Version || keVerController.Spec.SourceVersion == "") &&
		keVerController.Spec.TargetVersion < version.Version {
		// 1. Get the current CR.
		eventing, e := r.KnativeEventingClientSet.OperatorV1alpha1().KnativeEventings(keVerController.Namespace).Get("knative-eventing", metav1.GetOptions{})
		if apierrs.IsNotFound(e) {
			// The CR in the new version does not exist, so there is no need to downgrade.
			r.Logger.Infof("The Knative Eventing CR does not exist. You can directly install the older version" +
				"%s of the operator.", keVerController.Spec.TargetVersion)
			return nil
		} else if e != nil {
			return fmt.Errorf("Unable to get the CR of Knative Eventing in the new version target version %s",
				keVerController.Spec.TargetVersion)
		}

		// 2. Since the schema of the old CR may be removed, we use the unstructured.Unstructured and covert the CR of
		// the current version into unstructured.Unstructured.
		var target = &unstructured.Unstructured{}
		plural, err := downgrade.DowngradeCR(keVerController.Spec.TargetVersion,
			eventing, target, r.DynamicClientSet)
		if err != nil {
			return fmt.Errorf("Failed to downgrade the CR.")
		}

		err = oldschema.InstallOldCRD(keVerController.Spec.TargetVersion, r.kubeClient)
		if err != nil {
			return err
		}
		// 3. Create or update the target CR after the conversion for downgrade.
		targetDyn, err := r.DynamicClientSet.Resource(plural).Namespace(keVerController.Namespace).Get(target.GetName(),
			metav1.GetOptions{})
		if err == nil {
			// 4. If we can get the old CR, it means the CRD exists. We only need to update the unstructured.Unstructured.
			target.SetResourceVersion(targetDyn.GetResourceVersion())
			// 5. Update the version by updating the status.
			_, e = r.DynamicClientSet.Resource(plural).Namespace(keVerController.Namespace).UpdateStatus(target, metav1.UpdateOptions{})
			if e != nil {
				return fmt.Errorf("Failed to update the version of the CR into %s", keVerController.Spec.TargetVersion)
			} else {
				// 5. Update the CR itself.
				_, e = r.DynamicClientSet.Resource(plural).Namespace(keVerController.Namespace).Update(target, metav1.UpdateOptions{})
				if e != nil {
					return fmt.Errorf("Failed to update the the CR")
				}
			}
		} else {
			// 4. If the old CR does not exist, we need to create the CRD.
			// TODO Recreate the CRD of the old version.
			err := oldschema.InstallOldCRD(keVerController.Spec.TargetVersion, r.kubeClient)
			if err != nil {
				return err
			}

			// 5. When the CRD is ready, create the CR of the old version.
			_, e = r.DynamicClientSet.Resource(plural).Namespace(keVerController.Namespace).Create(target, metav1.CreateOptions{})
			if e != nil {
				return fmt.Errorf("Failed to create the CR in the version %s.", keVerController.Spec.TargetVersion)
			}
		}
		r.Logger.Infof("Successfully downgraded to the old CR version %s. Please install the older version of the operator %s.",
			keVerController.Spec.TargetVersion, keVerController.Spec.TargetVersion)
	} else {
		return fmt.Errorf("The target version you specified %s is not supported.", keVerController.Spec.TargetVersion)
	}

	return nil
}

func (r *Reconciler) updateStatus(desired *eventingv1alpha1.KEVersionController) (*eventingv1alpha1.KEVersionController, error) {
	ke, err := r.KnativeEventingClientSet.OperatorV1alpha1().KEVersionControllers(desired.Namespace).Get(desired.Name,
		metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	// If there's nothing to update, just return.
	if reflect.DeepEqual(ke.Status, desired.Status) {
		return ke, nil
	}
	// Don't modify the informers copy
	existing := ke.DeepCopy()
	existing.Status = desired.Status
	return r.KnativeEventingClientSet.OperatorV1alpha1().KEVersionControllers(desired.Namespace).UpdateStatus(existing)
}
