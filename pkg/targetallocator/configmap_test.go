// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package targetallocator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDesiredConfigMap(t *testing.T) {
	expectedLables := map[string]string{
		"app.kubernetes.io/managed-by": "opentelemetry-operator",
		"app.kubernetes.io/instance":   "default.my-instance",
		"app.kubernetes.io/part-of":    "opentelemetry",
		"app.kubernetes.io/version":    "0.47.0",
	}

	t.Run("should return expected target allocator config map", func(t *testing.T) {
		expectedLables["app.kubernetes.io/component"] = "opentelemetry-targetallocator"
		expectedLables["app.kubernetes.io/name"] = "my-instance-targetallocator"

		expectedData := map[string]string{
			"targetallocator.yaml": `allocation_strategy: least-weighted
config:
  scrape_configs:
  - job_name: otel-collector
    scrape_interval: 10s
    static_configs:
    - targets:
      - 0.0.0.0:8888
      - 0.0.0.0:9999
label_selector:
  app.kubernetes.io/component: opentelemetry-collector
  app.kubernetes.io/instance: default.my-instance
  app.kubernetes.io/managed-by: opentelemetry-operator
  app.kubernetes.io/part-of: opentelemetry
`,
		}

		actual, err := ConfigMap(collectorInstance())
		assert.NoError(t, err)

		assert.Equal(t, "my-instance-targetallocator", actual.Name)
		assert.Equal(t, expectedLables, actual.Labels)
		assert.Equal(t, expectedData, actual.Data)

	})
	t.Run("should return expected target allocator config map with label selectors", func(t *testing.T) {
		expectedLables["app.kubernetes.io/component"] = "opentelemetry-targetallocator"
		expectedLables["app.kubernetes.io/name"] = "my-instance-targetallocator"

		expectedData := map[string]string{
			"targetallocator.yaml": `allocation_strategy: least-weighted
config:
  scrape_configs:
  - job_name: otel-collector
    scrape_interval: 10s
    static_configs:
    - targets:
      - 0.0.0.0:8888
      - 0.0.0.0:9999
label_selector:
  app.kubernetes.io/component: opentelemetry-collector
  app.kubernetes.io/instance: default.my-instance
  app.kubernetes.io/managed-by: opentelemetry-operator
  app.kubernetes.io/part-of: opentelemetry
pod_monitor_selector:
  release: my-instance
service_monitor_selector:
  release: my-instance
`,
		}
		instance := collectorInstance()
		instance.Spec.TargetAllocator.PrometheusCR.PodMonitorSelector = map[string]string{
			"release": "my-instance",
		}
		instance.Spec.TargetAllocator.PrometheusCR.ServiceMonitorSelector = map[string]string{
			"release": "my-instance",
		}
		actual, err := ConfigMap(instance)
		assert.NoError(t, err)

		assert.Equal(t, "my-instance-targetallocator", actual.Name)
		assert.Equal(t, expectedLables, actual.Labels)
		assert.Equal(t, expectedData, actual.Data)

	})

}