package tests

import (
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// MustFindService finds a service in the given slice of objects with the
// given name, or fails the test.
func MustFindService(t testing.TB, objs []runtime.Object, name string) *corev1.Service {
	names := []string{}
	for _, obj := range objs {
		if service, ok := obj.(*corev1.Service); ok {
			if service.Name == name {
				return service
			}
			names = append(names, service.Name)
		}
	}

	t.Fatalf("failed to find service %q, found %v", name, names)
	return nil
}

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

// MustFindStatefulSet finds a stateful set in the given slice of objects with
// the given name, or fails the test.
func MustFindStatefulSet(t testing.TB, objs []runtime.Object, name string) *appsv1.StatefulSet {
	names := []string{}
	for _, obj := range objs {
		if statefulset, ok := obj.(*appsv1.StatefulSet); ok {
			if statefulset.Name == name {
				return statefulset
			}
			names = append(names, statefulset.Name)
		}
	}

	t.Fatalf("failed to find statefulset %q, found %v", name, names)
	return nil
}

// MustFindServiceAccount finds a service account in the given slice of objects
// with the given name, or fails the test.
func MustFindServiceAccount(t testing.TB, objs []runtime.Object, name string) *corev1.ServiceAccount {
	names := []string{}
	for _, obj := range objs {
		if serviceAccount, ok := obj.(*corev1.ServiceAccount); ok {
			if serviceAccount.Name == name {
				return serviceAccount
			}
			names = append(names, serviceAccount.Name)
		}
	}

	t.Fatalf("failed to find serviceaccount %q, found %v", name, names)
	return nil
}

// MustFindRole finds a role in the given slice of objects
// with the given name, or fails the test.
func MustFindRole(t testing.TB, objs []runtime.Object, name string) *rbacv1.Role {
	names := []string{}
	for _, obj := range objs {
		if role, ok := obj.(*rbacv1.Role); ok {
			if role.Name == name {
				return role
			}
			names = append(names, role.Name)
		}
	}

	t.Fatalf("failed to find role %q, found %v", name, names)
	return nil
}

// MustFindRoleBinding finds a role in the given slice of objects
// with the given name, or fails the test.
func MustFindRoleBinding(t testing.TB, objs []runtime.Object, name string) *rbacv1.RoleBinding {
	names := []string{}
	for _, obj := range objs {
		if roleBinding, ok := obj.(*rbacv1.RoleBinding); ok {
			if roleBinding.Name == name {
				return roleBinding
			}
			names = append(names, roleBinding.Name)
		}
	}

	t.Fatalf("failed to find rolebinding %q, found %v", name, names)
	return nil
}

// FindNetworkPolicy finds a network policy in the given slice of objects with
// the given name, or returns false if no policy with that name was found.
func FindNetworkPolicy(objs []runtime.Object, name string) (*networkingv1.NetworkPolicy, bool) {
	for _, obj := range objs {
		if policy, ok := obj.(*networkingv1.NetworkPolicy); ok {
			if policy.Name == name {
				return policy, true
			}
		}
	}

	return nil, false
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

// AssertNoVolumeMount asserts that no volume mount exists of the given name in the given
// slice of volumes.
func AssertNoVolumeMount(t testing.TB, vols []corev1.VolumeMount, name string) {
	for _, v := range vols {
		if v.Name == name {
			t.Fatalf("did not expect to find volume %q", name)
			return
		}
	}
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

// AssertEnvVar asserts that an environment variable exists with the given name.
// If it exists, it runs fn against the named environment variable.
func AssertEnvVar(t testing.TB, envs []corev1.EnvVar, name string, fn func(t testing.TB, env corev1.EnvVar)) {
	names := []string{}
	for _, env := range envs {
		if env.Name == name {
			fn(t, env)
			return
		}
		names = append(names, env.Name)
	}
	t.Fatalf("failed to find env var %q, found %v", name, names)
}

// AssertNoEnvVar asserts that an environment variable does not exist with the given name.
func AssertNoEnvVar(t testing.TB, envs []corev1.EnvVar, name string) {
	for _, env := range envs {
		if env.Name == name {
			t.Fatalf("did not expect to find env var %q", name)
			return
		}
	}
}
