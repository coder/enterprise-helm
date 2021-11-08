package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/chart/loader"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/utils/pointer"
)

// TestDefault loads the chart and checks metadata.
func TestDefault(t *testing.T) {
	t.Parallel()

	chart, err := loader.LoadDir("..")
	require.NoError(t, err, "loaded chart successfully")
	require.NotNil(t, chart, "chart must be non-nil")
	require.True(t, chart.IsRoot(), "chart must be a root chart")
	require.NoError(t, chart.Validate(), "chart has valid metadata")

	metadata := chart.Metadata
	require.Equal(t, "coder", metadata.Name, "unexpected chart name")
	require.False(t, metadata.Deprecated, "chart should not be deprecated")

	values, err := ConvertMapToCoderValues(chart.Values, false)
	require.NoError(t, err, "converted map to coder values")
	require.NotNil(t, values, "values must be non-nil")
	coderd := values.Coderd
	require.Equal(t, int32(1), *coderd.Replicas, "expected 1 replica by default")
}

func TestOverwriteReplica(t *testing.T) {
	t.Parallel()

	// Given
	// The default helm chart
	chart, err := loader.LoadDir("..")
	require.NoError(t, err, "loaded chart successfully")
	require.NotNil(t, chart, "chart must be non-nil")

	// When
	// We overwrite the replicas value and render the chart
	var ValuesToOverwrite = &CoderValues{
		Coderd: &CoderdValues{
			Replicas: pointer.Int32(3),
		},
	}

	objs, err := RenderChart(chart, ValuesToOverwrite, nil, nil)
	// TODO@jsjoeio - getting an error here
	// error deserializing "coder/templates/coderd.yaml": yaml: line 9: mapping values are not allowed in this context
	require.NoError(t, err, "failed to render chart")

	// Find the coderd Deployment
	var found bool
	for _, obj := range objs {
		deployment, ok := obj.(*appsv1.Deployment)
		if ok && deployment.Name == "coderd" {
			found = true

			// Then
			// We expect the rendered chart to have the values we overwrote
			expected := ValuesToOverwrite.Coderd.Replicas
			actual := deployment.Spec.Replicas
			require.Equal(t, expected, actual, "expected matching PodSecurityContext")
			break
		}
		require.True(t, found, "expected coderd deployment in manifests")
	}
}
