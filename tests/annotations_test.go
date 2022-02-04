package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ensures services.annotations values are applied to both coderd deployment &
// service
func TestAnnotations(t *testing.T) {
	t.Parallel()

	chart := LoadChart(t)

	expected := map[string]string{}

	objs := chart.MustRender(t, func(cv *CoderValues) {
		cv.Coderd.Annotations = expected
		cv.Coderd.ServiceSpec.Annotations = expected
		cv.Postgres.Default.Annotations = expected
	})

	depl := MustFindDeployment(t, objs, "coderd")
	assert.Equal(t, expected, depl.Annotations)

	svc := MustFindService(t, objs, "coderd")
	assert.Equal(t, expected, svc.Annotations)
}

// check if values are empty
func TestAnnotationsEmpty(t *testing.T) {
	t.Parallel()

	chart := LoadChart(t)

	objs := chart.MustRender(t, nil)

	depl := MustFindDeployment(t, objs, "coderd")
	assert.Empty(t, depl.Annotations)

	svc := MustFindService(t, objs, "coderd")
	assert.Empty(t, svc.Annotations)
}
