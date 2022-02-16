package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Ensures services.annotations values and the individual annotations per object
// are applied correctly.
func TestAnnotations(t *testing.T) {
	t.Parallel()

	var (
		chart = LoadChart(t)

		expectedGlobal = map[string]string{
			"global-key": "global-value",
			// Should be overwritten by some children.
			"key": "global-value",
		}
		expectedCoderd = map[string]string{
			"key":  "value",
			"key2": "value2",
		}
		expectedCoderdService = map[string]string{
			"key": "value",
			"service.beta.kubernetes.io/aws-load-balancer-backend-protocol": "http",
		}
		expectedTimescale = map[string]string{
			"key2": "value",
		}
	)

	objs := chart.MustRender(t, func(cv *CoderValues) {
		// Ensure backwards compatibility and merging order.
		cv.Services.Annotations = expectedGlobal

		cv.Coderd.Annotations = expectedCoderd
		cv.Coderd.ServiceAnnotations = expectedCoderdService
		cv.Postgres.Default.Annotations = expectedTimescale
	})

	depl := MustFindDeployment(t, objs, "coderd")
	assert.Equal(t, mergeAnnotations(expectedGlobal, expectedCoderd), depl.Annotations)

	svc := MustFindService(t, objs, "coderd")
	assert.Equal(t, mergeAnnotations(expectedGlobal, expectedCoderdService), svc.Annotations)

	db := MustFindStatefulSet(t, objs, "timescale")
	assert.Equal(t, mergeAnnotations(expectedGlobal, expectedTimescale), db.Annotations)
}

func TestAnnotationsEmpty(t *testing.T) {
	t.Parallel()

	var (
		chart = LoadChart(t)
		objs  = chart.MustRender(t, nil)
	)

	depl := MustFindDeployment(t, objs, "coderd")
	assert.Empty(t, depl.Annotations)

	svc := MustFindService(t, objs, "coderd")
	assert.Empty(t, svc.Annotations)

	db := MustFindStatefulSet(t, objs, "timescale")
	assert.Empty(t, db.Annotations)
}

func TestAnnotationsNull(t *testing.T) {
	t.Parallel()

	var (
		chart = LoadChart(t)
		objs  = chart.MustRender(t, func(cv *CoderValues) {
			cv.Coderd.Annotations = nil
			cv.Coderd.ServiceAnnotations = nil
			cv.Postgres.Default.Annotations = nil
			cv.Services.Annotations = nil
		})
	)

	depl := MustFindDeployment(t, objs, "coderd")
	assert.Empty(t, depl.Annotations)

	svc := MustFindService(t, objs, "coderd")
	assert.Empty(t, svc.Annotations)

	db := MustFindStatefulSet(t, objs, "timescale")
	assert.Empty(t, db.Annotations)
}

// mergeAnnotations copies `a` into a new map, then it copies all key/value
// pairs from `b` on top of that copy.
func mergeAnnotations(a, b map[string]string) map[string]string {
	out := map[string]string{}
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		out[k] = v
	}

	return out
}
