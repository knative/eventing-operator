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

package versioncontroller

import (
	"context"
	"k8s.io/client-go/tools/clientcmd"
	versionControllerinformer "knative.dev/eventing-operator/pkg/client/injection/informers/eventing/v1alpha1/keversioncontroller"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"

	rbase "knative.dev/eventing-operator/pkg/reconciler"
	"knative.dev/eventing-operator/pkg/reconciler/knativeeventing"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
)

const (
	controllerAgentName = "version-controller"
	reconcilerName      = "VersionController"
)

// NewController initializes the controller and is called by the generated code
// Registers eventhandlers to enqueue events
func NewController(
	ctx context.Context,
	cmw configmap.Watcher,
) *controller.Impl {
	versionControllerInformer := versionControllerinformer.Get(ctx)

	c := &Reconciler{
		Base:                    rbase.NewBase(ctx, controllerAgentName, cmw),
		versionControllerLister: versionControllerInformer.Lister(),
	}

	cfg, err := clientcmd.BuildConfigFromFlags(*knativeeventing.MasterURL, *knativeeventing.Kubeconfig)
	if err != nil {
		c.Logger.Error(err, "Error building kubeconfig")
	}

	kubeClient, err := apiextension.NewForConfig(cfg)
	if err != nil {
		c.Logger.Error(err, "Failed to create client to access the CRD")
	}
	c.kubeClient = kubeClient
	impl := controller.NewImpl(c, c.Logger, reconcilerName)

	c.Logger.Info("Setting up event handlers for %s", reconcilerName)

	versionControllerInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

	return impl
}
