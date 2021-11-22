package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestDefault loads the chart and checks metadata.
func TestDefault(t *testing.T) {
	t.Parallel()

	chart := LoadChart(t)
	require.NoError(t, chart.Validate(), "chart has valid metadata")

	metadata := chart.Metadata
	require.Equal(t, "coder", metadata.Name, "unexpected chart name")
	require.False(t, metadata.Deprecated, "chart should not be deprecated")

	coderd := chart.OriginalValues.Coderd
	require.Equal(t, 1, *coderd.Replicas, "expected 1 replica by default")
}
