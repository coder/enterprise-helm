package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/chart/loader"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"
)

func TestSecurityContext(t *testing.T) {
	t.Parallel()

	chart, err := loader.LoadDir("..")
	require.NoError(t, err, "loaded chart successfully")

	exampleOpenShift, err := ReadValues("../examples/openshift/openshift.values.yaml")
	require.NoError(t, err, "failed to load OpenShift example values")

	tests := []struct {
		Name                     string
		Values                   *CoderValues
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
		}, {
			Name:   "openshift",
			Values: exampleOpenShift,
			PodSecurityContext: &corev1.PodSecurityContext{
				RunAsUser:    nil,
				RunAsGroup:   nil,
				RunAsNonRoot: nil,
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

			objs, err := RenderChart(chart, test.Values, nil, nil)
			require.NoError(t, err, "failed to render chart")

			// Find the coderd Deployment
			var found bool
			for _, obj := range objs {
				deployment, ok := obj.(*appsv1.Deployment)
				if ok && deployment.Name == "coderd" {
					found = true

					expected := test.PodSecurityContext
					actual := deployment.Spec.Template.Spec.SecurityContext
					require.Equal(t, expected, actual, "expected matching PodSecurityContext")
					break
				}
			}
			require.True(t, found, "expected coderd deployment in manifests")
		})
	}
}
