// +build postupgrade

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

package e2e

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"knative.dev/pkg/test/logstream"
	"knative.dev/eventing-operator/test"
	"knative.dev/eventing-operator/test/resources"
	"knative.dev/eventing-operator/test/client"
)

// TestKnativeEventingUpgrade verifies the KnativeEventing creation, deployment recreation, and KnativeEventing deletion
// after upgraded to the latest HEAD at master, with the latest generated manifest of KnativeEventing.
func TestKnativeEventingUpgrade(t *testing.T) {
	cancel := logstream.Start(t)
	defer cancel()
	clients := client.Setup(t)

	names := test.ResourceNames{
		KnativeEventing: test.EventingOperatorName,
		Namespace:       test.EventingOperatorNamespace,
	}

	test.CleanupOnInterrupt(func() { test.TearDown(clients, names) })
	defer test.TearDown(clients, names)

	// Create a KnativeEventing
	if _, err := resources.EnsureKnativeEventingExists(clients.KnativeEventing(), names); err != nil {
		t.Fatalf("KnativeService %q failed to create: %v", names.KnativeEventing, err)
	}

	// Test if KnativeEventing can reach the READY status
	t.Run("create", func(t *testing.T) {
		resources.AssertKEOperatorCRReadyStatus(t, clients, names)
	})

	// Verify if resources match the latest requirement after upgrade
	t.Run("verify resources", func(t *testing.T) {
		resources.AssertKEOperatorCRReadyStatus(t, clients, names)
		// TODO: We only verify the deployment, but we need to add other resources as well, like ServiceAccount, ClusterRoleBinding, etc.
		expectedDeployments := []string{"eventing-controller", "eventing-webhook", "imc-controller",
			"imc-dispatcher", "broker-controller"}
		knativeEventingVerifyDeployment(t, clients, names, expectedDeployments)
	})

	// TODO: We will add one or sections here to run the tests tagged with postupgrade in knative evening.

	// Delete the KnativeEventing to see if all resources will be removed
	t.Run("delete", func(t *testing.T) {
		resources.AssertKEOperatorCRReadyStatus(t, clients, names)
		resources.KEOperatorCRDelete(t, clients, names)
	})
}

// knativeEventingVerifyDeployment verify whether the deployments have the correct number and names.
func knativeEventingVerifyDeployment(t *testing.T, clients *test.Clients, names test.ResourceNames,
	expectedDeployments []string) {
	dpList, err := clients.KubeClient.Kube.AppsV1().Deployments(names.Namespace).List(metav1.ListOptions{})
	assertEqual(t, err, nil)
	assertEqual(t, len(dpList.Items), len(expectedDeployments))
	for _, deployment := range dpList.Items {
		assertEqual(t, stringInList(deployment.Name, expectedDeployments), true)
	}
}

func assertEqual(t *testing.T, actual, expected interface{}) {
	if actual == expected {
		return
	}
	t.Fatalf("expected does not equal actual. \nExpected: %v\nActual: %v", expected, actual)
}

func stringInList(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

