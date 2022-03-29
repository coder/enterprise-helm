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
		cv.Coderd.Liveness.TimeoutSeconds = 4
		cv.Coderd.Readiness.InitialDelaySeconds = 5
		cv.Coderd.Readiness.FailureThreshold = 6
		cv.Coderd.Readiness.PeriodSeconds = 7
		cv.Coderd.Readiness.TimeoutSeconds = 8
	}
	objs, err := chart.Render(override, &opts, chartutil.DefaultCapabilities)
	require.NoError(t, err)
	deployment := MustFindDeployment(t, objs, "coderd")
	AssertContainer(t, deployment.Spec.Template.Spec.Containers, "coderd", func(t testing.TB, v v1.Container) {
		require.EqualValues(t, 1, v.LivenessProbe.InitialDelaySeconds)
		require.EqualValues(t, 2, v.LivenessProbe.FailureThreshold)
		require.EqualValues(t, 3, v.LivenessProbe.PeriodSeconds)
		require.EqualValues(t, 4, v.LivenessProbe.TimeoutSeconds)
		require.EqualValues(t, 5, v.ReadinessProbe.InitialDelaySeconds)
		require.EqualValues(t, 6, v.ReadinessProbe.FailureThreshold)
		require.EqualValues(t, 7, v.ReadinessProbe.PeriodSeconds)
		require.EqualValues(t, 8, v.ReadinessProbe.TimeoutSeconds)
	})
}
