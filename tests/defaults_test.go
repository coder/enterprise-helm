package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/utils/pointer"
)

// TestDefault loads the chart and checks metadata.
func TestDefault(t *testing.T) {
	t.Parallel()

	chart := LoadChart(t)
	require.NoError(t, chart.Validate(), "chart has valid metadata")

	metadata := chart.Metadata
	require.Equal(t, "coder", metadata.Name, "unexpected chart name")
	require.False(t, metadata.Deprecated, "chart should not be deprecated")

	objs := chart.MustRender(t, nil)
	deployment := MustFindDeployment(t, objs, "coderd")

	require.Equal(t, pointer.Int32(1), deployment.Spec.Replicas, "expected 1 replica by default")
	podSpec := deployment.Spec.Template.Spec
	require.Len(t, podSpec.Containers, 1, "pod spec should have 1 container")
	require.Equal(t, "docker.io/coderenvs/coder-service:1.25.0", podSpec.Containers[0].Image,
		"expected default image name")
	vars := EnvVarsAsMap(podSpec.Containers[0].Env)
	t.Logf("vars: %v", vars)
	require.Equal(t, "docker.io/coderenvs/envbox:1.25.0", vars["ENVBOX_IMAGE"],
		"expected default envbox image name")

	require.Len(t, podSpec.InitContainers, 1, "pod spec should have 1 init container")
	require.Equal(t, "docker.io/coderenvs/coder-service:1.25.0", podSpec.InitContainers[0].Image,
		"expected default image name")
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
