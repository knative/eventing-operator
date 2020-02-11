// +build e2e

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

package e2e

import (
	"errors"
	"path/filepath"
	"runtime"
	"testing"

	"k8s.io/apimachinery/pkg/api/meta"

	mf "github.com/manifestival/manifestival"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"knative.dev/eventing-operator/test"
	"knative.dev/eventing-operator/test/resources"
	"knative.dev/pkg/test/logstream"
)

// TestKnativeEventingDeployment verifies the KnativeEventing creation, deployment recreation, and KnativeEventing deletion.
func TestKnativeEventingDeployment(t *testing.T) {
	cancel := logstream.Start(t)
	defer cancel()
	clients := Setup(t)

	names := test.ResourceNames{
		KnativeEventing: test.EventingOperatorName,
		Namespace:       test.EventingOperatorNamespace,
	}

	test.CleanupOnInterrupt(func() { test.TearDown(clients, names) })
	defer test.TearDown(clients, names)

	// Create a KnativeEventing
	if _, err := resources.CreateKnativeEventing(clients.KnativeEventing(), names); err != nil {
		t.Fatalf("KnativeService %q failed to create: %v", names.KnativeEventing, err)
	}

	// Test if KnativeEventing can reach the READY status
	t.Run("create", func(t *testing.T) {
		knativeEventingVerify(t, clients, names)
	})

	// Delete the deployments one by one to see if they will be recreated.
	t.Run("restore", func(t *testing.T) {
		knativeEventingVerify(t, clients, names)
		deploymentRecreation(t, clients, names)
	})

	// Delete the KnativeEventing to see if all resources will be removed
	t.Run("delete", func(t *testing.T) {
		knativeEventingVerify(t, clients, names)
		knativeEventingDelete(t, clients, names)
	})
}

// knativeEventingVerify verifies if the KnativeEventing can reach the READY status.
func knativeEventingVerify(t *testing.T, clients *test.Clients, names test.ResourceNames) {
	if _, err := resources.WaitForKnativeEventingState(clients.KnativeEventing(), names.KnativeEventing,
		resources.IsKnativeEventingReady); err != nil {
		t.Fatalf("KnativeService %q failed to get to the READY status: %v", names.KnativeEventing, err)
	}

}

// deploymentRecreation verify whether all the deployments for knative eventing are able to recreate, when they are deleted.
func deploymentRecreation(t *testing.T, clients *test.Clients, names test.ResourceNames) {
	dpList, err := clients.KubeClient.Kube.AppsV1().Deployments(names.Namespace).List(metav1.ListOptions{})
	if err != nil {
		t.Fatalf("Failed to get any deployment under the namespace %q: %v",
			test.EventingOperatorNamespace, err)
	}
	if len(dpList.Items) == 0 {
		t.Fatalf("No deployment under the namespace %q was found",
			test.EventingOperatorNamespace)
	}
	// Delete the first deployment and verify the operator recreates it
	deployment := dpList.Items[0]
	if err := clients.KubeClient.Kube.AppsV1().Deployments(deployment.Namespace).Delete(deployment.Name,
		&metav1.DeleteOptions{}); err != nil {
		t.Fatalf("Failed to delete deployment %s/%s: %v", deployment.Namespace, deployment.Name, err)
	}

	waitErr := wait.PollImmediate(resources.Interval, resources.Timeout, func() (bool, error) {
		dep, err := clients.KubeClient.Kube.AppsV1().Deployments(deployment.Namespace).Get(deployment.Name, metav1.GetOptions{})
		if err != nil {
			// If the deployment is not found, we continue to wait for the availability.
			if apierrs.IsNotFound(err) {
				return false, nil
			}
			return false, err
		}
		return resources.IsDeploymentAvailable(dep)
	})

	if waitErr != nil {
		t.Fatalf("The deployment %s/%s failed to reach the desired state: %v", deployment.Namespace, deployment.Name, err)
	}

	if _, err := resources.WaitForKnativeEventingState(clients.KnativeEventing(), test.EventingOperatorName,
		resources.IsKnativeEventingReady); err != nil {
		t.Fatalf("KnativeService %q failed to reach the desired state: %v", test.EventingOperatorName, err)
	}
	t.Logf("The deployment %s/%s reached the desired state.", deployment.Namespace, deployment.Name)
}

// knativeEventingDelete deletes tha KnativeEventing to see if all resources will be deleted
func knativeEventingDelete(t *testing.T, clients *test.Clients, names test.ResourceNames) {
	if err := clients.KnativeEventing().Delete(names.KnativeEventing, &metav1.DeleteOptions{}); err != nil {
		t.Fatalf("KnativeEventing %q failed to delete: %v", names.KnativeEventing, err)
	}
	_, b, _, _ := runtime.Caller(0)
	m, err := mf.NewManifest(filepath.Join((filepath.Dir(b)+"/.."), "config/"), false, clients.Config)
	if err != nil {
		t.Fatal("Failed to load manifest", err)
	}
	if err := verifyNoKnativeEventings(clients); err != nil {
		t.Fatal(err)
	}
	for _, u := range m.Resources {
		if u.GetKind() == "Namespace" {
			// The namespace should be skipped, because when the CR is removed, the Manifest to be removed has
			// been modified, since the namespace can be injected.
			continue
		}
		waitErr := wait.PollImmediate(resources.Interval, resources.Timeout, func() (bool, error) {
			gvrs, _ := meta.UnsafeGuessKindToResource(u.GroupVersionKind())
			if _, err := clients.Dynamic.Resource(gvrs).Get(u.GetName(), metav1.GetOptions{}); apierrs.IsNotFound(err) {
				return true, nil
			}
			return false, err
		})

		if waitErr != nil {
			t.Fatalf("The %s %s failed to be deleted: %v", u.GetKind(), u.GetName(), waitErr)
		}
		t.Logf("The %s %s has been deleted.", u.GetKind(), u.GetName())
	}
}

func verifyNoKnativeEventings(clients *test.Clients) error {
	eventings, err := clients.KnativeEventingAll().List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	if len(eventings.Items) > 0 {
		return errors.New("Unable to verify cluster-scoped resources are deleted if any KnativeEventing exists")
	}
	return nil
}
