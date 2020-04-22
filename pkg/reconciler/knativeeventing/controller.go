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

package knativeeventing

import (
	"context"
	"os"
	"path/filepath"

	"github.com/go-logr/zapr"
	mfc "github.com/manifestival/client-go-client"
	mf "github.com/manifestival/manifestival"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/tools/cache"

	"k8s.io/client-go/rest"
	"knative.dev/eventing-operator/pkg/apis/eventing/v1alpha1"
	knativeEventinginformer "knative.dev/eventing-operator/pkg/client/injection/informers/eventing/v1alpha1/knativeeventing"
	rbase "knative.dev/eventing-operator/pkg/reconciler"
	deploymentinformer "knative.dev/pkg/client/injection/kube/informers/apps/v1/deployment"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
)

const (
	controllerAgentName = "knativeeventing-controller"
	reconcilerName      = "KnativeEventing"
)

// NewControllerWithConfig initializes the controller and is called by the generated code
// Registers eventhandlers to enqueue events
func NewControllerWithConfig(cfg *rest.Config) func(context.Context, configmap.Watcher) *controller.Impl {
	return func(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
		knativeEventingInformer := knativeEventinginformer.Get(ctx)
		deploymentInformer := deploymentinformer.Get(ctx)

		c := &Reconciler{
			Base:                  rbase.NewBase(ctx, controllerAgentName, cmw),
			knativeEventingLister: knativeEventingInformer.Lister(),
			eventings:             sets.String{},
		}

		koDataDir := os.Getenv("KO_DATA_PATH")

		config, err := mfc.NewManifest(filepath.Join(koDataDir, "knative-eventing/"), cfg,
			mf.UseLogger(zapr.NewLogger(c.Logger.Desugar()).WithName("manifestival")))
		if err != nil {
			c.Logger.Error(err, "Error creating the Manifest for knative-eventing")
			os.Exit(1)
		}

		c.config = config
		impl := controller.NewImpl(c, c.Logger, reconcilerName)

		c.Logger.Info("Setting up event handlers for %s", reconcilerName)

		knativeEventingInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

		deploymentInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
			FilterFunc: controller.Filter(v1alpha1.SchemeGroupVersion.WithKind("KnativeEventing")),
			Handler:    controller.HandleAll(impl.EnqueueControllerOf),
		})

		return impl
	}
}
