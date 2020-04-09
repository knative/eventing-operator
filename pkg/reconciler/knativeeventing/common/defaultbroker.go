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

package common

import (
	mf "github.com/manifestival/manifestival"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes/scheme"

	eventingv1alpha1 "knative.dev/eventing-operator/pkg/apis/eventing/v1alpha1"
)

const (
	channelBasedBrokerClass = "ChannelBasedBroker"
)

var defaultBrokerConfigMapData = map[string]map[string]string{
	"clusterDefault": {
		"brokerClass": channelBasedBrokerClass,
		"apiVersion":  "v1",
		"kind":        "ConfigMap",
		"name":        "config-br-default-channel",
		"namespace":   "knative-eventing",
	},
}

// DefaultBrokerConfigMapTransform updates the default broker configMap with the value defined in the spec
func DefaultBrokerConfigMapTransform(instance *eventingv1alpha1.KnativeEventing, log *zap.SugaredLogger) mf.Transformer {
	return func(u *unstructured.Unstructured) error {
		if u.GetKind() == "ConfigMap" && u.GetName() == "config-br-defaults" {
			var configMap = &corev1.ConfigMap{}
			err := scheme.Scheme.Convert(u, configMap, nil)
			if err != nil {
				log.Error(err, "Error converting Unstructured to ConfigMap", "unstructured", u, "configMap", configMap)
				return err
			}

			defaultBrokerClass := instance.Spec.DefaultBrokerClass
			if defaultBrokerClass == "" {
				defaultBrokerClass = channelBasedBrokerClass
			}

			defaultBrokerConfigMapData["clusterDefault"]["brokerClass"] = defaultBrokerClass
			out, err := yaml.Marshal(&defaultBrokerConfigMapData)

			if err != nil {
				return err
			}

			configMap.Data["default-br-config"] = string(out)

			err = scheme.Scheme.Convert(configMap, u, nil)
			if err != nil {
				return err
			}
			// The zero-value timestamp defaulted by the conversion causes
			// superfluous updates
			u.SetCreationTimestamp(metav1.Time{})
			log.Debugw("Finished updating config-br-defaults configMap", "name", u.GetName(), "unstructured", u.Object)
		}
		return nil
	}
}
