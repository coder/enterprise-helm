package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/chart/loader"
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
	require.Equal(t, 1, *coderd.Replicas, "expected 1 replica by default")
}
