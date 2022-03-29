package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/chartutil"
	v1 "k8s.io/api/core/v1"
)

// TestHealthz loads the chart and checks overriding liveness and readiness
func TestHealthz(t *testing.T) {
	t.Parallel()

	chart := LoadChart(t)
	opts := DefaultReleaseOptions()
	override := func(cv *CoderValues) {
		cv.Coderd.Liveness.InitialDelaySeconds = 1
		cv.Coderd.Liveness.FailureThreshold = 2
		cv.Coderd.Liveness.PeriodSeconds = 3
		cv.Coderd.Readiness.InitialDelaySeconds = 4
		cv.Coderd.Readiness.FailureThreshold = 5
		cv.Coderd.Readiness.PeriodSeconds = 6
	}
	objs, err := chart.Render(override, &opts, chartutil.DefaultCapabilities)
	require.NoError(t, err)
	deployment := MustFindDeployment(t, objs, "coderd")
	AssertContainer(t, deployment.Spec.Template.Spec.Containers, "coderd", func(t testing.TB, v v1.Container) {
		require.EqualValues(t, v.LivenessProbe.InitialDelaySeconds, 1)
		require.EqualValues(t, v.LivenessProbe.FailureThreshold, 2)
		require.EqualValues(t, v.LivenessProbe.PeriodSeconds, 3)
		require.EqualValues(t, v.ReadinessProbe.InitialDelaySeconds, 4)
		require.EqualValues(t, v.ReadinessProbe.FailureThreshold, 5)
		require.EqualValues(t, v.ReadinessProbe.PeriodSeconds, 6)
	})
}
