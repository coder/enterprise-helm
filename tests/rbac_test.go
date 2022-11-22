package tests

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"

	"github.com/stretchr/testify/assert"
)

func TestRBAC(t *testing.T) {
	t.Parallel()

	chart := LoadChart(t)

	for _, tc := range []struct {
		Name                                 string
		ValuesFunc                           func(v *CoderValues)
		AssertEnvironmentsServiceAccountFunc func(t testing.TB, spec *corev1.ServiceAccount)
		AssertCoderServiceAccountFunc        func(t testing.TB, spec *corev1.ServiceAccount)
		AssertRoleFunc                       func(t testing.TB, spec *rbacv1.Role)
		AssertRoleBindingFunc                func(t testing.TB, spec *rbacv1.RoleBinding)
	}{
		{
			Name: "Defaults",
			AssertEnvironmentsServiceAccountFunc: func(t testing.TB, spec *corev1.ServiceAccount) {
				assert.Empty(t, spec.Annotations)
			},
			AssertCoderServiceAccountFunc: func(t testing.TB, spec *corev1.ServiceAccount) {
				assert.Empty(t, spec.Annotations)
			},
			AssertRoleFunc: func(t testing.TB, spec *rbacv1.Role) {
				assert.NotEmpty(t, spec.Rules)
			},
			AssertRoleBindingFunc: func(t testing.TB, spec *rbacv1.RoleBinding) {
				if assert.Len(t, spec.Subjects, 1) {
					assert.Equal(t, spec.Subjects[0].Name, "coder")
					assert.Equal(t, spec.Subjects[0].Kind, "ServiceAccount")
				}
				assert.Equal(t, spec.RoleRef.Name, "coder")
				assert.Equal(t, spec.RoleRef.Kind, "Role")
			},
		},
	} {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			objs := chart.MustRender(t, tc.ValuesFunc)
			envSA := MustFindServiceAccount(t, objs, "environments")
			coderSA := MustFindServiceAccount(t, objs, "coder")
			role := MustFindRole(t, objs, "coder")
			roleBinding := MustFindRoleBinding(t, objs, "coder")
			tc.AssertEnvironmentsServiceAccountFunc(t, envSA)
			tc.AssertCoderServiceAccountFunc(t, coderSA)
			tc.AssertRoleFunc(t, role)
			tc.AssertRoleBindingFunc(t, roleBinding)
		})
	}
}
