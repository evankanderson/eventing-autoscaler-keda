/*
Copyright 2021 The Knative Authors

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

package kafkasource

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"

	duckv1 "knative.dev/pkg/apis/duck/v1"

	"knative.dev/reconciler-test/pkg/feature"
	"knative.dev/reconciler-test/pkg/k8s"
	"knative.dev/reconciler-test/pkg/manifest"
)

type CfgFn func(map[string]interface{})

func Gvr() schema.GroupVersionResource {
	return schema.GroupVersionResource{Group: "sources.knative.dev", Version: "v1beta1", Resource: "kafkasources"}
}

// Install will create a KafkaSource resource, augmented with the config fn options.
func Install(name string, opts ...CfgFn) feature.StepFn {
	cfg := map[string]interface{}{
		"name": name,
	}
	for _, fn := range opts {
		fn(cfg)
	}
	return func(ctx context.Context, t feature.T) {
		if _, err := manifest.InstallLocalYaml(ctx, cfg); err != nil {
			t.Fatal(err, cfg)
		}
	}
}

// IsReady tests to see if a KafkaSource becomes ready within the time given.
func IsReady(name string, timings ...time.Duration) feature.StepFn {
	return k8s.IsReady(Gvr(), name, timings...)
}

// WithBootstrapServers adds the bootstrapServers config to a KafkaSource spec.
func WithBootstrapServers(bootstrapServers []string) CfgFn {
	return func(cfg map[string]interface{}) {
		if bootstrapServers != nil {
			cfg["bootstrapServers"] = bootstrapServers
		}
	}
}

// WithTopics adds the topics config to a KafkaSource spec.
func WithTopics(topics []string) CfgFn {
	return func(cfg map[string]interface{}) {
		if topics != nil {
			cfg["topics"] = topics
		}
	}
}

// WithSink adds the sink related config to a PingSource spec.
func WithSink(ref *duckv1.KReference, uri string) CfgFn {
	return func(cfg map[string]interface{}) {
		if _, set := cfg["sink"]; !set {
			cfg["sink"] = map[string]interface{}{}
		}
		sink := cfg["sink"].(map[string]interface{})

		if uri != "" {
			sink["uri"] = uri
		}
		if ref != nil {
			if _, set := sink["ref"]; !set {
				sink["ref"] = map[string]interface{}{}
			}
			sref := sink["ref"].(map[string]interface{})
			sref["apiVersion"] = ref.APIVersion
			sref["kind"] = ref.Kind
			// skip namespace
			sref["name"] = ref.Name
		}
	}
}
