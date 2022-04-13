package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
)

func TestBuiltInProvider(t *testing.T) {
	t.Parallel()

	chart := LoadChart(t)

	migrate := true
	objs := chart.MustRender(t, func(cv *CoderValues) {
		cv.Coderd.BuiltinProviderServiceAccount.Migrate = &migrate
	})

	depl := MustFindDeployment(t, objs, "coderd")
	for _, container := range depl.Spec.Template.Spec.Containers {
		assert.Subsetf(t, container.Env, []v1.EnvVar{
			{
				Name:  "CODER_MIGRATE_BUILT_IN_PROVIDER",
				Value: "true",
			},
		}, "container %q", container.Name)
	}

	migrate = false
	objs = chart.MustRender(t, func(cv *CoderValues) {
		cv.Coderd.BuiltinProviderServiceAccount.Migrate = &migrate
	})

	depl = MustFindDeployment(t, objs, "coderd")
	for _, container := range depl.Spec.Template.Spec.Containers {
		assert.Subsetf(t, container.Env, []v1.EnvVar{
			{
				Name:  "CODER_MIGRATE_BUILT_IN_PROVIDER",
				Value: "false",
			},
		}, "container %q", container.Name)
	}
}
