package tests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"
)

func TestImages(t *testing.T) {
	t.Parallel()

	chart := LoadChart(t)

	tests := []struct {
		Name       string
		ValuesFunc func(v *CoderValues)
		AssertFunc func(t testing.TB, coderdSpec *corev1.PodSpec, timescaleSpec *corev1.PodSpec)
	}{
		{
			Name:       "default",
			ValuesFunc: nil,
			AssertFunc: func(t testing.TB, coderdSpec *corev1.PodSpec, timescaleSpec *corev1.PodSpec) {
				require.Len(t, coderdSpec.Containers, 1, "pod spec should have 1 container")

				imageRef := fmt.Sprintf("docker.io/coderenvs/coder-service:%s", chart.Metadata.AppVersion)
				require.Equal(t, imageRef, coderdSpec.Containers[0].Image, "coderd image ref")

				// envbox image ref
				vars := EnvVarsAsMap(coderdSpec.Containers[0].Env)
				imageRef = fmt.Sprintf("docker.io/coderenvs/envbox:%s", chart.Metadata.AppVersion)
				require.Equal(t, imageRef, vars["ENVBOX_IMAGE"], "envbox image should match expected setting")

				imageRef = fmt.Sprintf("docker.io/coderenvs/timescale:%s", chart.Metadata.AppVersion)
				require.Equal(t, imageRef, timescaleSpec.Containers[0].Image, "timescale image ref")
			},
		},
		{
			Name: "custom-images",
			ValuesFunc: func(cv *CoderValues) {
				cv.Coderd.Image = pointer.String("coder-service:latest")
				cv.Postgres.Default.Image = pointer.String("postgresql:12")
				cv.Envbox.Image = pointer.String("docker.io/abcorg/envbox:1.25")
			},
			AssertFunc: func(t testing.TB, coderdSpec *corev1.PodSpec, timescaleSpec *corev1.PodSpec) {
				require.Len(t, coderdSpec.Containers, 1, "pod spec should have 1 container")

				require.Equal(t, "coder-service:latest", coderdSpec.Containers[0].Image, "coderd image ref")

				// envbox image ref
				vars := EnvVarsAsMap(coderdSpec.Containers[0].Env)
				require.Equal(t, "docker.io/abcorg/envbox:1.25", vars["ENVBOX_IMAGE"], "envbox image should match expected setting")

				require.Equal(t, "postgresql:12", timescaleSpec.Containers[0].Image, "timescale image ref")
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()

			objs := chart.MustRender(t, test.ValuesFunc)
			deployment := MustFindDeployment(t, objs, "coderd")
			statefulset := MustFindStatefulSet(t, objs, "timescale")
			test.AssertFunc(t, &deployment.Spec.Template.Spec, &statefulset.Spec.Template.Spec)
		})
	}
}
