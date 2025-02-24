// Copyright 2020, OpenTelemetry Authors
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

package awsemfexporter

import (
	"testing"

	commonpb "github.com/census-instrumentation/opencensus-proto/gen-go/agent/common/v1"
	agentmetricspb "github.com/census-instrumentation/opencensus-proto/gen-go/agent/metrics/v1"
	resourcepb "github.com/census-instrumentation/opencensus-proto/gen-go/resource/v1"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/model/pdata"
	conventions "go.opentelemetry.io/collector/translator/conventions/v1.5.0"
	"go.opentelemetry.io/collector/translator/internaldata"
	"go.uber.org/zap"
)

func TestReplacePatternValidTaskId(t *testing.T) {
	logger := zap.NewNop()

	input := "{TaskId}"

	attrMap := pdata.NewAttributeMap()
	attrMap.UpsertString("aws.ecs.cluster.name", "test-cluster-name")
	attrMap.UpsertString("aws.ecs.task.id", "test-task-id")

	s := replacePatterns(input, attrMap, logger)

	assert.Equal(t, "test-task-id", s)
}

func TestReplacePatternValidClusterName(t *testing.T) {
	logger := zap.NewNop()

	input := "/aws/ecs/containerinsights/{ClusterName}/performance"

	attrMap := pdata.NewAttributeMap()
	attrMap.UpsertString("aws.ecs.cluster.name", "test-cluster-name")
	attrMap.UpsertString("aws.ecs.task.id", "test-task-id")

	s := replacePatterns(input, attrMap, logger)

	assert.Equal(t, "/aws/ecs/containerinsights/test-cluster-name/performance", s)
}

func TestReplacePatternMissingAttribute(t *testing.T) {
	logger := zap.NewNop()

	input := "/aws/ecs/containerinsights/{ClusterName}/performance"

	attrMap := pdata.NewAttributeMap()
	attrMap.UpsertString("aws.ecs.task.id", "test-task-id")

	s := replacePatterns(input, attrMap, logger)

	assert.Equal(t, "/aws/ecs/containerinsights/undefined/performance", s)
}

func TestReplacePatternAttrPlaceholderClusterName(t *testing.T) {
	logger := zap.NewNop()

	input := "/aws/ecs/containerinsights/{ClusterName}/performance"

	attrMap := pdata.NewAttributeMap()
	attrMap.UpsertString("ClusterName", "test-cluster-name")

	s := replacePatterns(input, attrMap, logger)

	assert.Equal(t, "/aws/ecs/containerinsights/test-cluster-name/performance", s)
}

func TestReplacePatternWrongKey(t *testing.T) {
	logger := zap.NewNop()

	input := "/aws/ecs/containerinsights/{WrongKey}/performance"

	attrMap := pdata.NewAttributeMap()
	attrMap.UpsertString("ClusterName", "test-task-id")

	s := replacePatterns(input, attrMap, logger)

	assert.Equal(t, "/aws/ecs/containerinsights/{WrongKey}/performance", s)
}

func TestReplacePatternNilAttrValue(t *testing.T) {
	logger := zap.NewNop()

	input := "/aws/ecs/containerinsights/{ClusterName}/performance"

	attrMap := pdata.NewAttributeMap()
	attrMap.InsertNull("ClusterName")

	s := replacePatterns(input, attrMap, logger)

	assert.Equal(t, "/aws/ecs/containerinsights/undefined/performance", s)
}

func TestReplacePatternValidTaskDefinitionFamily(t *testing.T) {
	logger := zap.NewNop()

	input := "{TaskDefinitionFamily}"

	attrMap := pdata.NewAttributeMap()
	attrMap.UpsertString("aws.ecs.cluster.name", "test-cluster-name")
	attrMap.UpsertString("aws.ecs.task.family", "test-task-definition-family")

	s := replacePatterns(input, attrMap, logger)

	assert.Equal(t, "test-task-definition-family", s)
}

func TestGetNamespace(t *testing.T) {
	defaultMetric := createMetricTestData()
	testCases := []struct {
		testName        string
		metric          *agentmetricspb.ExportMetricsServiceRequest
		configNamespace string
		namespace       string
	}{
		{
			"non-empty namespace",
			defaultMetric,
			"namespace",
			"namespace",
		},
		{
			"empty namespace",
			defaultMetric,
			"",
			"myServiceNS/myServiceName",
		},
		{
			"empty namespace, no service namespace",
			&agentmetricspb.ExportMetricsServiceRequest{
				Resource: &resourcepb.Resource{
					Labels: map[string]string{
						conventions.AttributeServiceName: "myServiceName",
					},
				},
			},
			"",
			"myServiceName",
		},
		{
			"empty namespace, no service name",
			&agentmetricspb.ExportMetricsServiceRequest{
				Resource: &resourcepb.Resource{
					Labels: map[string]string{
						conventions.AttributeServiceNamespace: "myServiceNS",
					},
				},
			},
			"",
			"myServiceNS",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			rms := internaldata.OCToMetrics(tc.metric.Node, tc.metric.Resource, tc.metric.Metrics)
			rm := rms.ResourceMetrics().At(0)
			namespace := getNamespace(&rm, tc.configNamespace)
			assert.Equal(t, tc.namespace, namespace)
		})
	}
}

func TestGetLogInfo(t *testing.T) {
	metrics := []*agentmetricspb.ExportMetricsServiceRequest{
		{
			Node: &commonpb.Node{
				ServiceInfo: &commonpb.ServiceInfo{Name: "test-emf"},
				LibraryInfo: &commonpb.LibraryInfo{ExporterVersion: "SomeVersion"},
			},
			Resource: &resourcepb.Resource{
				Labels: map[string]string{
					"aws.ecs.cluster.name":          "test-cluster-name",
					"aws.ecs.task.id":               "test-task-id",
					"k8s.node.name":                 "ip-192-168-58-245.ec2.internal",
					"aws.ecs.container.instance.id": "203e0410260d466bab7873bb4f317b4e",
					"aws.ecs.task.family":           "test-task-definition-family",
				},
			},
		},
		{
			Node: &commonpb.Node{
				ServiceInfo: &commonpb.ServiceInfo{Name: "test-emf"},
				LibraryInfo: &commonpb.LibraryInfo{ExporterVersion: "SomeVersion"},
			},
			Resource: &resourcepb.Resource{
				Labels: map[string]string{
					"ClusterName":          "test-cluster-name",
					"TaskId":               "test-task-id",
					"NodeName":             "ip-192-168-58-245.ec2.internal",
					"ContainerInstanceId":  "203e0410260d466bab7873bb4f317b4e",
					"TaskDefinitionFamily": "test-task-definition-family",
				},
			},
		},
	}

	var rms []pdata.ResourceMetrics
	for _, md := range metrics {
		rms = append(rms, internaldata.OCToMetrics(md.Node, md.Resource, md.Metrics).ResourceMetrics().At(0))
	}

	testCases := []struct {
		testName        string
		namespace       string
		configLogGroup  string
		configLogStream string
		logGroup        string
		logStream       string
	}{
		{
			"non-empty namespace, no config",
			"namespace",
			"",
			"",
			"/metrics/namespace",
			"",
		},
		{
			"empty namespace, no config",
			"",
			"",
			"",
			"",
			"",
		},
		{
			"non-empty namespace, config w/o pattern",
			"namespace",
			"test-logGroupName",
			"test-logStreamName",
			"test-logGroupName",
			"test-logStreamName",
		},
		{
			"empty namespace, config w/o pattern",
			"",
			"test-logGroupName",
			"test-logStreamName",
			"test-logGroupName",
			"test-logStreamName",
		},
		{
			"non-empty namespace, config w/ pattern",
			"namespace",
			"/aws/ecs/containerinsights/{ClusterName}/performance",
			"{TaskId}",
			"/aws/ecs/containerinsights/test-cluster-name/performance",
			"test-task-id",
		},
		{
			"empty namespace, config w/ pattern",
			"",
			"/aws/ecs/containerinsights/{ClusterName}/performance",
			"{TaskId}",
			"/aws/ecs/containerinsights/test-cluster-name/performance",
			"test-task-id",
		},
		//test case for aws container insight usage
		{
			"empty namespace, config w/ pattern",
			"",
			"/aws/containerinsights/{ClusterName}/performance",
			"{NodeName}",
			"/aws/containerinsights/test-cluster-name/performance",
			"ip-192-168-58-245.ec2.internal",
		},
		// test case for AWS ECS EC2 container insights usage
		{
			"empty namespace, config w/ pattern",
			"",
			"/aws/containerinsights/{ClusterName}/performance",
			"instanceTelemetry/{ContainerInstanceId}",
			"/aws/containerinsights/test-cluster-name/performance",
			"instanceTelemetry/203e0410260d466bab7873bb4f317b4e",
		},
		{
			"empty namespace, config w/ pattern",
			"",
			"/aws/containerinsights/{ClusterName}/performance",
			"{TaskDefinitionFamily}-{TaskId}",
			"/aws/containerinsights/test-cluster-name/performance",
			"test-task-definition-family-test-task-id",
		},
	}

	for i := range rms {
		for _, tc := range testCases {
			t.Run(tc.testName, func(t *testing.T) {
				config := &Config{
					LogGroupName:  tc.configLogGroup,
					LogStreamName: tc.configLogStream,
				}
				logGroup, logStream := getLogInfo(&rms[i], tc.namespace, config)
				assert.Equal(t, tc.logGroup, logGroup)
				assert.Equal(t, tc.logStream, logStream)
			})
		}
	}

}
