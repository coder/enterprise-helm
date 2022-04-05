package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestExtraEnvs(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		envs []corev1.EnvVar
	}{
		{
			name: "Empty",
			envs: nil,
		},
		{
			name: "One",
			envs: []corev1.EnvVar{
				{
					Name:  "DEAN_WAS_HERE",
					Value: "true",
				},
			},
		},
		{
			name: "Many",
			envs: []corev1.EnvVar{
				{
					Name:  "DEAN_WAS_HERE",
					Value: "true",
				},
				{
					Name:  "COLIN_WAS_HERE",
					Value: "false",
				},
			},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			chart := LoadChart(t)
			objs := chart.MustRender(t, func(cv *CoderValues) {
				cv.Coderd.ExtraEnvs = c.envs
			})

			depl := MustFindDeployment(t, objs, "coderd")
			for _, container := range depl.Spec.Template.Spec.InitContainers {
				assert.Subsetf(t, container.Env, c.envs, "init container %q", container.Name)
			}
			for _, container := range depl.Spec.Template.Spec.Containers {
				assert.Subsetf(t, container.Env, c.envs, "container %q", container.Name)
			}
		})
	}
}
