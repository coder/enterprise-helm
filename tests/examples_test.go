package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"
)

// TestExamples loads the example values files and performs
// some basic checks.
func TestExamples(t *testing.T) {
	t.Parallel()

	chart, err := loader.LoadDir("..")
	require.NoError(t, err, "loaded chart successfully")
	require.NotNil(t, chart, "chart must be non-nil")

	exampleOpenShift, err := ReadValuesAsMap("../examples/openshift/openshift.values.yaml")
	require.NoError(t, err, "failed to load OpenShift example values")

	exampleKind, err := ReadValuesAsMap("../examples/kind/kind.values.yaml")
	require.NoError(t, err, "failed to load Kind example values")

	tests := []struct {
		Name                     string
		Values                   map[string]interface{}
		PodSecurityContext       *corev1.PodSecurityContext
		ContainerSecurityContext *corev1.SecurityContext
	}{
		{
			Name:   "default",
			Values: nil,
			PodSecurityContext: &corev1.PodSecurityContext{
				RunAsUser:    pointer.Int64(1000),
				RunAsGroup:   nil,
				RunAsNonRoot: pointer.Bool(true),
				SeccompProfile: &corev1.SeccompProfile{
					Type:             corev1.SeccompProfileTypeRuntimeDefault,
					LocalhostProfile: nil,
				},
			},
			ContainerSecurityContext: &corev1.SecurityContext{
				RunAsUser:                nil,
				RunAsGroup:               nil,
				RunAsNonRoot:             nil,
				Capabilities:             nil,
				Privileged:               nil,
				SELinuxOptions:           nil,
				WindowsOptions:           nil,
				ReadOnlyRootFilesystem:   pointer.Bool(true),
				AllowPrivilegeEscalation: pointer.Bool(false),
				ProcMount:                nil,
				SeccompProfile: &corev1.SeccompProfile{
					Type:             corev1.SeccompProfileTypeRuntimeDefault,
					LocalhostProfile: nil,
				},
			},
		}, {
			Name:   "openshift",
			Values: exampleOpenShift,
			PodSecurityContext: &corev1.PodSecurityContext{
				RunAsUser:      nil,
				RunAsGroup:     nil,
				RunAsNonRoot:   pointer.Bool(true),
				SeccompProfile: nil,
			},
			ContainerSecurityContext: &corev1.SecurityContext{
				RunAsUser:                nil,
				RunAsGroup:               nil,
				RunAsNonRoot:             nil,
				Capabilities:             nil,
				Privileged:               nil,
				SELinuxOptions:           nil,
				WindowsOptions:           nil,
				ReadOnlyRootFilesystem:   pointer.Bool(true),
				AllowPrivilegeEscalation: pointer.Bool(false),
				ProcMount:                nil,
				SeccompProfile:           nil,
			},
		},
		{
			Name:   "kind",
			Values: exampleKind,
			PodSecurityContext: &corev1.PodSecurityContext{
				RunAsUser:    pointer.Int64(1000),
				RunAsNonRoot: pointer.Bool(true),
				SeccompProfile: &corev1.SeccompProfile{
					Type: corev1.SeccompProfileTypeRuntimeDefault,
				},
			},
			ContainerSecurityContext: &corev1.SecurityContext{
				RunAsUser:                pointer.Int64(1000),
				RunAsGroup:               pointer.Int64(1000),
				RunAsNonRoot:             pointer.Bool(true),
				Capabilities:             nil,
				Privileged:               nil,
				SELinuxOptions:           nil,
				WindowsOptions:           nil,
				ReadOnlyRootFilesystem:   pointer.Bool(true),
				AllowPrivilegeEscalation: pointer.Bool(false),
				ProcMount:                nil,
				SeccompProfile: &corev1.SeccompProfile{
					Type:             corev1.SeccompProfileTypeRuntimeDefault,
					LocalhostProfile: nil,
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()

			values, err := chartutil.ToRenderValues(chart, test.Values, DefaultReleaseOptions(), chartutil.DefaultCapabilities.Copy())
			require.NoError(t, err, "failed to generate render values")

			manifests, err := engine.Render(chart, values)
			require.NoError(t, err, "failed to render chart")

			objs, err := LoadObjectsFromManifests(manifests)
			require.NoError(t, err, "failed to convert manifests to objects")

			// Find the coderd Deployment
			var found bool
			for _, obj := range objs {
				deployment, ok := obj.(*appsv1.Deployment)
				if ok && deployment.Name == "coderd" {
					found = true

					require.Equal(t, test.PodSecurityContext,
						deployment.Spec.Template.Spec.SecurityContext,
						"expected matching pod securityContext")
					require.Len(t, deployment.Spec.Template.Spec.Containers, 1,
						"expected one container")
					require.Equal(t, test.ContainerSecurityContext,
						deployment.Spec.Template.Spec.Containers[0].SecurityContext,
						"expected matching container securityContext")

					break
				}
			}
			require.True(t, found, "expected coderd deployment in manifests")
		})
	}
}
