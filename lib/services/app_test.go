/*
Copyright 2021 Gravitational, Inc.

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

package services

import (
	"testing"

	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gravitational/teleport/api/types"
	"github.com/gravitational/teleport/lib/utils"
)

// TestApplicationUnmarshal verifies an app resource can be unmarshaled.
func TestApplicationUnmarshal(t *testing.T) {
	expected, err := types.NewAppV3(types.Metadata{
		Name:        "test-app",
		Description: "Test description",
		Labels:      map[string]string{"env": "dev"},
	}, types.AppSpecV3{
		URI: "http://localhost:8080",
	})
	require.NoError(t, err)
	data, err := utils.ToJSON([]byte(appYAML))
	require.NoError(t, err)
	actual, err := UnmarshalApp(data)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

// TestApplicationMarshal verifies a marshaled app resource can be unmarshaled back.
func TestApplicationMarshal(t *testing.T) {
	expected, err := types.NewAppV3(types.Metadata{
		Name:        "test-app",
		Description: "Test description",
		Labels:      map[string]string{"env": "dev"},
	}, types.AppSpecV3{
		URI: "http://localhost:8080",
	})
	require.NoError(t, err)
	data, err := MarshalApp(expected)
	require.NoError(t, err)
	actual, err := UnmarshalApp(data)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

var appYAML = `kind: app
version: v3
metadata:
  name: test-app
  description: "Test description"
  labels:
    env: dev
spec:
  uri: "http://localhost:8080"`

func TestGetAppName(t *testing.T) {
	tests := []struct {
		serviceName string
		namespace   string
		clusterName string
		portName    string
		annotation  string
		wantErr     string
		expected    string
	}{
		{
			serviceName: "service1",
			namespace:   "ns1",
			clusterName: "cluster1",
			portName:    "port1",
			expected:    "service1-port1-ns1-cluster1",
		},
		{
			serviceName: "service2",
			namespace:   "ns2",
			clusterName: "cluster2",
			expected:    "service2-ns2-cluster2",
		},
		{
			serviceName: "service3",
			namespace:   "ns3",
			clusterName: "cluster.3",
			expected:    "service3-ns3-cluster-3",
		},
		{
			serviceName: "service3",
			namespace:   "ns3",
			clusterName: "cluster.3",
			annotation:  "overridden-name",
			expected:    "overridden-name",
		},
		{
			serviceName: "service3",
			namespace:   "ns3",
			clusterName: "cluster.3",
			portName:    "http",
			annotation:  "overridden-name",
			expected:    "overridden-name-http",
		},
		{
			serviceName: "service3",
			namespace:   "ns3",
			clusterName: "cluster.3",
			portName:    "http",
			annotation:  "overridden*name",
			wantErr:     "s",
		},
	}

	for _, tt := range tests {
		result, err := getAppName(tt.serviceName, tt.namespace, tt.clusterName, tt.portName, tt.annotation)
		if tt.wantErr != "" {
			require.ErrorContains(t, err, tt.wantErr)
		} else {
			require.Equal(t, tt.expected, result)
		}
	}
}

func TestGetServiceFQDN(t *testing.T) {
	tests := []struct {
		name         string
		namespace    string
		externalName string
		expected     string
	}{
		{
			name:      "service1",
			namespace: "ns1",
			expected:  "service1.ns1.svc.cluster.local",
		},
		{
			name:         "service2",
			externalName: "external-service2",
			namespace:    "ns2",
			expected:     "external-service2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			service := v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      tt.name,
					Namespace: tt.namespace,
				},
				Spec: v1.ServiceSpec{
					ExternalName: tt.externalName,
				},
			}
			if tt.externalName != "" {
				service.Spec.Type = v1.ServiceTypeExternalName
			}
			require.Equal(t, tt.expected, getServiceFQDN(service))
		})
	}
}

func TestBuildAppURI(t *testing.T) {
	tests := []struct {
		serviceFQDN string
		port        int32
		protocol    string
		expected    string
	}{
		{
			serviceFQDN: "service.example",
			port:        8080,
			protocol:    "http",
			expected:    "http://service.example:8080",
		},
		{
			serviceFQDN: "service.example",
			port:        443,
			protocol:    "https",
			expected:    "https://service.example:443",
		},
		{
			serviceFQDN: "special.service.example",
			port:        42,
			protocol:    "tcp",
			expected:    "tcp://special.service.example:42",
		},
	}

	for _, tt := range tests {
		require.Equal(t, buildAppURI(tt.protocol, tt.serviceFQDN, tt.port), tt.expected)
	}
}

func TestGetAppRewriteConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected types.Rewrite
		wantErr  string
	}{
		{
			input: `redirect:
  - "redirect1"
  - "redirect2"`,
			expected: types.Rewrite{
				Redirect: []string{"redirect1", "redirect2"},
			},
		},
		{
			input:    `wrong:"`,
			expected: types.Rewrite{},
			wantErr:  "failed decoding rewrite config",
		},
		{
			input: `redirect:
  - "localhost"
  - "jenkins.internal.dev"
headers:
  - name: "X-Custom-Header"
    value: "example"
  - name: "Authorization"
    value: "Bearer {{internal.jwt}}"`,
			expected: types.Rewrite{
				Redirect: []string{"localhost", "jenkins.internal.dev"},
				Headers: []*types.Header{
					{
						Name:  "X-Custom-Header",
						Value: "example",
					},
					{
						Name:  "Authorization",
						Value: "Bearer {{internal.jwt}}",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		result, err := getAppRewriteConfig(map[string]string{types.DiscoveryAppRewriteLabel: tt.input})
		if tt.wantErr != "" {
			require.ErrorContains(t, err, tt.wantErr)
		} else {
			require.NoError(t, err)
			require.Equal(t, tt.expected, *result)
		}
	}
}

func TestGetAppLabels(t *testing.T) {
	tests := []struct {
		labels      map[string]string
		clusterName string
		expected    map[string]string
	}{
		{
			labels:      map[string]string{"label1": "value1"},
			clusterName: "cluster1",
			expected:    map[string]string{"label1": "value1", types.KubernetesClusterLabel: "cluster1"},
		},
		{
			labels:      map[string]string{"label1": "value1", "label2": "value2"},
			clusterName: "cluster2",
			expected:    map[string]string{"label1": "value1", "label2": "value2", types.KubernetesClusterLabel: "cluster2"},
		},
	}

	for _, tt := range tests {
		result, err := getAppLabels(tt.labels, tt.clusterName)
		require.NoError(t, err)

		require.Equal(t, tt.expected, result)
	}
}
