package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"
)

// TestExamples loads the example values files and performs
// some basic checks.
func TestExamples(t *testing.T) {
	t.Parallel()

	chart := LoadChart(t)

	exampleOpenShift, err := ReadValuesFileAsMap("../examples/openshift/openshift.values.yaml")
	require.NoError(t, err, "failed to load OpenShift example values")

	exampleKind, err := ReadValuesFileAsMap("../examples/kind/kind.values.yaml")
	require.NoError(t, err, "failed to load Kind example values")

	tests := []struct {
		Name                     string
		Values                   map[string]interface{}
		PodSecurityContext       *corev1.PodSecurityContext
		ContainerSecurityContext *corev1.SecurityContext
		Postgres                 *PostgresValues
	}{
		{
			Name:   "default",
			Values: nil,
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

	var (
		defaultPsp = &corev1.PodSecurityContext{
			RunAsUser:    pointer.Int64(1000),
			RunAsNonRoot: pointer.Bool(true),
			SeccompProfile: &corev1.SeccompProfile{
				Type: corev1.SeccompProfileTypeRuntimeDefault,
			},
		}

		defaultCsc = &corev1.SecurityContext{
			ReadOnlyRootFilesystem:   pointer.Bool(true),
			AllowPrivilegeEscalation: pointer.Bool(false),
			SeccompProfile: &corev1.SeccompProfile{
				Type: corev1.SeccompProfileTypeRuntimeDefault,
			},
		}
	)

	for _, test := range tests {
		test := test

		if test.PodSecurityContext == nil {
			test.PodSecurityContext = defaultPsp
		}
		if test.ContainerSecurityContext == nil {
			test.ContainerSecurityContext = defaultCsc
		}

		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()

			values, err := chartutil.ToRenderValues(chart.chart, test.Values, DefaultReleaseOptions(), chartutil.DefaultCapabilities.Copy())
			require.NoError(t, err, "failed to generate render values")

			manifests, err := engine.Render(chart.chart, values)
			require.NoError(t, err, "failed to render chart")

			objs, err := LoadObjectsFromManifests(manifests)
			require.NoError(t, err, "failed to convert manifests to objects")

			// Find the coderd Deployment
			coderd := FindDeployment(t, objs, "coderd")

			assert.Equal(t, test.PodSecurityContext, coderd.Spec.Template.Spec.SecurityContext,
				"expected matching pod securityContext",
			)
			require.Len(t, coderd.Spec.Template.Spec.Containers, 1,
				"expected one container",
			)
			assert.Equal(t, test.ContainerSecurityContext, coderd.Spec.Template.Spec.Containers[0].SecurityContext,
				"expected matching container securityContext",
			)
		})
	}
}
