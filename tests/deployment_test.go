package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeployment(t *testing.T) {
	t.Parallel()

	t.Run("Labels", func(t *testing.T) {
		var (
			expectedLabels = map[string]string{
				"app":                       "coderd",
				"app.kubernetes.io/name":    "coderd",
				"app.kubernetes.io/part-of": "coder",
				"coder.deployment":          "coderd",
			}
			extraLabels = map[string]string{
				"foo": "bar",
			}

			objs = LoadChart(t).MustRender(t, func(cv *CoderValues) {
				cv.Coderd.ExtraLabels = extraLabels
			})
			coderd = MustFindDeployment(t, objs, "coderd")
		)

		for k, v := range extraLabels {
			if _, found := expectedLabels[k]; !found {
				expectedLabels[k] = v
			}
		}

		require.EqualValues(t, expectedLabels, coderd.Spec.Template.Labels)
	})
}
