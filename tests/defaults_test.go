package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/meta"
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

// TestNamespace checks that all objects are created in the specified
// release namespace.
func TestNamespace(t *testing.T) {
	t.Parallel()

	chart := LoadChart(t)
	opts := DefaultReleaseOptions()
	namespaces := []string{
		opts.Namespace,
		"coder-test",
	}
	for _, namespace := range namespaces {
		namespace := namespace
		opts := opts
		opts.Namespace = namespace
		t.Run(namespace, func(t *testing.T) {
			t.Parallel()

			// Render the chart with default values
			objs, err := chart.Render(nil, &opts, nil)
			require.NoError(t, err, "chart render failed")

			// Verify that all objects are using the supplied namespace
			for _, obj := range objs {
				metaObject, err := meta.Accessor(obj)
				require.NoError(t, err, "failed to get object metadata")

				actualNamespace := metaObject.GetNamespace()
				require.Equal(t, namespace, actualNamespace,
					"deployed namespace does not match target")
			}
		})
	}
}
