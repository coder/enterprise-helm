package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/chartutil"
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
	tag := chart.AppVersion()
	require.Equal(t, "docker.io/coderenvs/coder-service:"+tag, podSpec.Containers[0].Image,
		"expected default image name")
	require.Empty(t, chart.OriginalValues.Coderd.Image, "expected default coderd image to be empty")
	require.Empty(t, chart.OriginalValues.Envbox.Image, "expected default envbox image to be empty")
	require.Empty(t, chart.OriginalValues.Postgres.Default.Image, "expected default timescale image to be empty")
	vars := EnvVarsAsMap(podSpec.Containers[0].Env)
	require.Equal(t, "docker.io/coderenvs/envbox:"+tag, vars["ENVBOX_IMAGE"],
		"expected default envbox image name")

	require.Len(t, podSpec.InitContainers, 1, "pod spec should have 1 init container")
	require.Equal(t, "docker.io/coderenvs/coder-service:"+tag, podSpec.InitContainers[0].Image,
		"expected default image name")
}

// TestMetadata checks that all objects are created with expected metadata.
// This checks that the release name is as expected and expected labels are
// present.
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

		// Create a release that installs into the given namespace and
		// with the given name
		opts := opts
		opts.Name = namespace + "-release"
		opts.Namespace = namespace

		t.Run(namespace, func(t *testing.T) {
			t.Parallel()

			// Render the chart with default values
			objs, err := chart.Render(nil, &opts, nil)
			require.NoError(t, err, "chart render failed")

			for _, obj := range objs {
				metaObject, err := meta.Accessor(obj)
				require.NoErrorf(t, err, "failed to get object metadata for object %q with name %q",
					obj.GetObjectKind().GroupVersionKind(), metaObject.GetName())

				// Verify that all objects are using the supplied namespace
				actualNamespace := metaObject.GetNamespace()
				require.Equalf(t, namespace, actualNamespace,
					"deployed namespace does not match for object %q with name %q",
					obj.GetObjectKind().GroupVersionKind(), metaObject.GetName())

				// Check that labels are present and values match what we expect
				labels := metaObject.GetLabels()
				require.Equalf(t, chart.Name(), labels["app.kubernetes.io/name"],
					"chart name matches for object %q with name %q",
					obj.GetObjectKind().GroupVersionKind(), metaObject.GetName())
				require.Containsf(t, labels["helm.sh/chart"], chart.Name(),
					"objects include chart name label for object %q with name %q",
					obj.GetObjectKind().GroupVersionKind(), metaObject.GetName())
				require.Equalf(t, "Helm", labels["app.kubernetes.io/managed-by"],
					"objects are managed by Helm for object %q with name %q",
					obj.GetObjectKind().GroupVersionKind(), metaObject.GetName())
				require.Equalf(t, namespace+"-release", labels["app.kubernetes.io/instance"],
					"object instance label matches Helm release name for object %q with name %q",
					obj.GetObjectKind().GroupVersionKind(), metaObject.GetName())
				require.Equalf(t, chart.AppVersion(), labels["app.kubernetes.io/version"],
					"objects version label matches Helm appVersion for object %q with name %q",
					obj.GetObjectKind().GroupVersionKind(), metaObject.GetName())
			}
		})
	}
}

func TestVersion(t *testing.T) {
	t.Parallel()

	chart := LoadChart(t)

	tests := []struct {
		Name       string
		Version    string
		Compatible bool
	}{
		{
			Name:       "gke-incompatible-1.15",
			Version:    "1.15.11-gke.15",
			Compatible: false,
		},
		{
			Name:    "gke-outdated-1.19",
			Version: "1.19.13-gke.1900",
			// Soft compatibility
			Compatible: true,
		},
		{
			Name:    "gke-outdated-1.20",
			Version: "1.20.12-gke.1500",
			// Soft compatibility
			Compatible: true,
		},
		{
			Name:       "gke-current-1.21",
			Version:    "1.21.5-gke.1802",
			Compatible: true,
		},
		{
			Name:       "gke-current-1.22",
			Version:    "1.22.3-gke.1500",
			Compatible: true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()

			compatible := chartutil.IsCompatibleRange(chart.Metadata.KubeVersion, test.Version)
			if test.Compatible {
				require.True(t, compatible, "expected version %q to be compatible with constraint %q",
					test.Version, chart.Metadata.KubeVersion)
			} else {
				require.False(t, compatible, "expected version %q not to be compatible with constraint %q",
					test.Version, chart.Metadata.KubeVersion)
			}
		})
	}
}
