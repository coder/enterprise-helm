package tests

import (
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// MustFindDeployment finds a deployment in the given slice of objects with the
// given name, or fails the test.
func MustFindDeployment(t testing.TB, objs []runtime.Object, name string) *appsv1.Deployment {
	names := []string{}
	for _, obj := range objs {
		if deployment, ok := obj.(*appsv1.Deployment); ok {
			if deployment.Name == name {
				return deployment
			}
			names = append(names, deployment.Name)
		}
	}

	t.Fatalf("failed to find deployment %q, found %v", name, names)
	return nil
}

// EnvVarsAsMap converts simple key/value environment variable pairs into a
// map, ignoring variables using a ConfigMap or Secret source. If a variable
// is defined multiple times, the last value will be returned.
func EnvVarsAsMap(variables []corev1.EnvVar) map[string]string {
	values := map[string]string{}

	for _, v := range variables {
		if v.ValueFrom != nil {
			continue
		}

		values[v.Name] = v.Value
	}

	return values
}

// AssertVolume asserts that a volume exists of the given name in the given
// slice of volumes. If it exists, it also runs fn against the named volume.
func AssertVolume(t testing.TB, vols []corev1.Volume, name string, fn func(t testing.TB, v corev1.Volume)) {
	names := []string{}
	for _, v := range vols {
		if v.Name == name {
			fn(t, v)
			return
		}
		names = append(names, v.Name)
	}

	t.Fatalf("failed to find volume %q, found %v", name, names)
}

// AssertVolumeMount asserts that a volume mount exists of the given name in the
// given slice of volume mounts. If it exists, it also runs fn against the named
// volume mount.
func AssertVolumeMount(t testing.TB, vols []corev1.VolumeMount, name string, fn func(t testing.TB, v corev1.VolumeMount)) {
	names := []string{}
	for _, v := range vols {
		if v.Name == name {
			fn(t, v)
			return
		}
		names = append(names, v.Name)
	}

	t.Fatalf("failed to find volume mount %q, found %v", name, names)
}

// AssertContainer asserts that a container exists of the given name in the
// given slice of containers. If it exists, it also runs fn against the named
// container.
func AssertContainer(t testing.TB, cnts []corev1.Container, name string, fn func(t testing.TB, v corev1.Container)) {
	names := []string{}
	for _, c := range cnts {
		if c.Name == name {
			fn(t, c)
			return
		}
		names = append(names, c.Name)
	}

	t.Fatalf("failed to find container %q, found %v", name, names)
}
