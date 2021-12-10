package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"
)

func TestProxy(t *testing.T) {
	t.Parallel()

	chart := LoadChart(t)

	tests := []struct {
		Name       string
		ValuesFunc func(v *CoderValues)
		AssertFunc func(t testing.TB, spec *corev1.PodSpec)
	}{
		{
			Name:       "default",
			ValuesFunc: nil,
			AssertFunc: func(t testing.TB, spec *corev1.PodSpec) {
				require.Len(t, spec.Containers, 1, "pod spec should have 1 container")
				vars := EnvVarsAsMap(spec.Containers[0].Env)
				require.Empty(t, vars["https_proxy"], "https_proxy should be empty")
				require.Empty(t, vars["http_proxy"], "http_proxy should be empty")
				require.Empty(t, vars["no_proxy"], "no_proxy should be empty")

				require.Len(t, spec.InitContainers, 1, "pod spec should have 1 init container")
				vars = EnvVarsAsMap(spec.InitContainers[0].Env)
				require.Empty(t, vars["https_proxy"], "https_proxy should be empty")
				require.Empty(t, vars["http_proxy"], "http_proxy should be empty")
				require.Empty(t, vars["no_proxy"], "no_proxy should be empty")
			},
		},
		{
			Name: "all_proxy",
			ValuesFunc: func(v *CoderValues) {
				v.Coderd.Proxy.HTTPS = pointer.String("http://proxy.coder.com:3128")
				v.Coderd.Proxy.HTTP = pointer.String("https://proxy.coder.com:8888")
				v.Coderd.Proxy.Exempt = pointer.String("coder.com,coder.app")
			},
			AssertFunc: func(t testing.TB, spec *corev1.PodSpec) {
				require.Len(t, spec.Containers, 1, "pod spec should have 1 container")
				vars := EnvVarsAsMap(spec.Containers[0].Env)
				require.Equal(t, "http://proxy.coder.com:3128", vars["https_proxy"], "http_proxy did not match")
				require.Equal(t, "https://proxy.coder.com:8888", vars["http_proxy"], "https_proxy did not match")
				require.Equal(t, "coder.com,coder.app", vars["no_proxy"], "no_proxy did not match")

				require.Len(t, spec.InitContainers, 1, "pod spec should have 1 init container")
				vars = EnvVarsAsMap(spec.InitContainers[0].Env)
				require.Equal(t, "http://proxy.coder.com:3128", vars["https_proxy"], "http_proxy did not match")
				require.Equal(t, "https://proxy.coder.com:8888", vars["http_proxy"], "https_proxy did not match")
				require.Equal(t, "coder.com,coder.app", vars["no_proxy"], "no_proxy did not match")
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()

			objs := chart.MustRender(t, test.ValuesFunc)
			deployment := MustFindDeployment(t, objs, "coderd")
			test.AssertFunc(t, &deployment.Spec.Template.Spec)
		})
	}
}
