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
				require.Equal(t, "cluster.local", vars["no_proxy"], "no_proxy did not match")

				require.Len(t, spec.InitContainers, 1, "pod spec should have 1 init container")
				vars = EnvVarsAsMap(spec.InitContainers[0].Env)
				require.Empty(t, vars["https_proxy"], "https_proxy should be empty")
				require.Empty(t, vars["http_proxy"], "http_proxy should be empty")
				require.Equal(t, "cluster.local", vars["no_proxy"], "no_proxy did not match")
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

func TestReverseProxy(t *testing.T) {
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
				require.Empty(t, vars["PROXY_TRUSTED_ORIGINS"], "PROXY_TRUSTED_ORIGINS should be empty")
				require.Empty(t, vars["PROXY_TRUSTED_HEADERS"], "http_proxy should be empty")

				require.Len(t, spec.InitContainers, 1, "pod spec should have 1 init container")
				vars = EnvVarsAsMap(spec.InitContainers[0].Env)
				require.NotContains(t, vars, "PROXY_TRUSTED_ORIGINS", "init container should not have PROXY_TRUSTED_ORIGINS")
				require.NotContains(t, vars, "PROXY_TRUSTED_HEADERS", "init container should not have PROXY_TRUSTED_HEADERS")
			},
		},
		{
			Name: "single-proxy-header",
			ValuesFunc: func(v *CoderValues) {
				v.Coderd.ReverseProxy.Headers = []string{"X-Forwarded-For"}
				v.Coderd.ReverseProxy.TrustedOrigins = []string{"127.0.0.1/8"}
			},
			AssertFunc: func(t testing.TB, spec *corev1.PodSpec) {
				require.Len(t, spec.Containers, 1, "pod spec should have 1 container")
				vars := EnvVarsAsMap(spec.Containers[0].Env)
				require.Equal(t, "127.0.0.1/8", vars["PROXY_TRUSTED_ORIGINS"], "PROXY_TRUSTED_ORIGINS should contain 127.0.0.1/8")
				require.Equal(t, "X-Forwarded-For", vars["PROXY_TRUSTED_HEADERS"], "PROXY_TRUSTED_HEADERS should contain X-Forwarded-For")

				require.Len(t, spec.InitContainers, 1, "pod spec should have 1 init container")
				vars = EnvVarsAsMap(spec.InitContainers[0].Env)
				require.NotContains(t, vars, "PROXY_TRUSTED_ORIGINS", "init container should not have PROXY_TRUSTED_ORIGINS")
				require.NotContains(t, vars, "PROXY_TRUSTED_HEADERS", "init container should not have PROXY_TRUSTED_HEADERS")
			},
		},
		{
			Name: "multiple-proxy-headers",
			ValuesFunc: func(v *CoderValues) {
				v.Coderd.ReverseProxy.Headers = []string{"X-Real-IP", "X-Forwarded-For"}
				v.Coderd.ReverseProxy.TrustedOrigins = []string{"127.0.0.1/8", "10.0.0.0/8"}
			},
			AssertFunc: func(t testing.TB, spec *corev1.PodSpec) {
				require.Len(t, spec.Containers, 1, "pod spec should have 1 container")
				vars := EnvVarsAsMap(spec.Containers[0].Env)
				require.Equal(t, "127.0.0.1/8,10.0.0.0/8", vars["PROXY_TRUSTED_ORIGINS"], "PROXY_TRUSTED_ORIGINS should be 127.0.0.1/8,10.0.0.0/8")
				require.Equal(t, "X-Real-IP,X-Forwarded-For", vars["PROXY_TRUSTED_HEADERS"], "PROXY_TRUSTED_HEADERS should contain X-Real-IP,X-Forwarded-For")

				require.Len(t, spec.InitContainers, 1, "pod spec should have 1 init container")
				vars = EnvVarsAsMap(spec.InitContainers[0].Env)
				require.NotContains(t, vars, "PROXY_TRUSTED_ORIGINS", "init container should not have PROXY_TRUSTED_ORIGINS")
				require.NotContains(t, vars, "PROXY_TRUSTED_HEADERS", "init container should not have PROXY_TRUSTED_HEADERS")
			},
		},
		{
			Name: "cloudflare-headers",
			ValuesFunc: func(v *CoderValues) {
				v.Coderd.ReverseProxy.Headers = []string{"CF-Connecting-IP", "X-Forwarded-For"}
				// List published by Cloudflare: https://www.cloudflare.com/ips/
				v.Coderd.ReverseProxy.TrustedOrigins = []string{
					"103.21.244.0/22",
					"103.22.200.0/22",
					"103.31.4.0/22",
					"104.16.0.0/13",
					"104.24.0.0/14",
					"108.162.192.0/18",
					"131.0.72.0/22",
					"141.101.64.0/18",
					"162.158.0.0/15",
					"172.64.0.0/13",
					"173.245.48.0/20",
					"188.114.96.0/20",
					"190.93.240.0/20",
					"197.234.240.0/22",
					"198.41.128.0/17",
					"2400:cb00::/32",
					"2606:4700::/32",
					"2803:f800::/32",
					"2405:b500::/32",
					"2405:8100::/32",
					"2a06:98c0::/29",
					"2c0f:f248::/32",
				}
			},
			AssertFunc: func(t testing.TB, spec *corev1.PodSpec) {
				require.Len(t, spec.Containers, 1, "pod spec should have 1 container")
				vars := EnvVarsAsMap(spec.Containers[0].Env)
				require.Equal(t, "103.21.244.0/22,103.22.200.0/22,103.31.4.0/22,104.16.0.0/13,104.24.0.0/14,108.162.192.0/18,131.0.72.0/22,141.101.64.0/18,162.158.0.0/15,172.64.0.0/13,173.245.48.0/20,188.114.96.0/20,190.93.240.0/20,197.234.240.0/22,198.41.128.0/17,2400:cb00::/32,2606:4700::/32,2803:f800::/32,2405:b500::/32,2405:8100::/32,2a06:98c0::/29,2c0f:f248::/32", vars["PROXY_TRUSTED_ORIGINS"], "PROXY_TRUSTED_ORIGINS should be Cloudflare range")
				require.Equal(t, "CF-Connecting-IP,X-Forwarded-For", vars["PROXY_TRUSTED_HEADERS"], "PROXY_TRUSTED_HEADERS should contain CF-Connecting-IP,X-Forwarded-For")

				require.Len(t, spec.InitContainers, 1, "pod spec should have 1 init container")
				vars = EnvVarsAsMap(spec.InitContainers[0].Env)
				require.NotContains(t, vars, "PROXY_TRUSTED_ORIGINS", "init container should not have PROXY_TRUSTED_ORIGINS")
				require.NotContains(t, vars, "PROXY_TRUSTED_HEADERS", "init container should not have PROXY_TRUSTED_HEADERS")
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
