package tests

import (
	"testing"

	netv1 "k8s.io/api/networking/v1"
	"k8s.io/utils/pointer"

	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/require"
)

func TestIngress(t *testing.T) {
	chart := LoadChart(t)

	pathTypePrefix := netv1.PathTypePrefix
	coderdIngressRule := netv1.IngressRuleValue{
		HTTP: &netv1.HTTPIngressRuleValue{
			Paths: []netv1.HTTPIngressPath{
				{
					Path:     "/",
					PathType: &pathTypePrefix,
					Backend: netv1.IngressBackend{
						Service: &netv1.IngressServiceBackend{
							Name: "coderd",
							Port: netv1.ServiceBackendPort{
								Name: "tcp-coderd",
							},
						},
					},
				},
			},
		},
	}

	tests := []struct {
		Name string
		// ValuesFunc is called to configure the values used in this test.
		// The function should override the CoderValues as appropriate for
		// the test in question.
		ValuesFunc func(v *CoderValues)
		// AssertFunc is called after rendering the chart, with the resulting
		// Ingress object. You can use it to assert properties about the
		// Ingress object.
		AssertFunc func(t *testing.T, ingress *netv1.Ingress)
	}{
		{
			Name: "simple-ingress",
			ValuesFunc: func(v *CoderValues) {
				v.Ingress.Enable = pointer.Bool(true)
				v.Ingress.Host = pointer.String("install.coder.com")
				v.Coderd.DevURLsHost = pointer.String("*.install.coder.app")
			},
			AssertFunc: func(t *testing.T, ingress *netv1.Ingress) {
				defaultAnnotations := map[string]string{
					"nginx.ingress.kubernetes.io/proxy-body-size": "0",
				}
				require.Equal(t, defaultAnnotations, ingress.Annotations)

				require.Empty(t, ingress.Spec.IngressClassName)

				expectedRules := []netv1.IngressRule{
					{
						Host:             "install.coder.com",
						IngressRuleValue: coderdIngressRule,
					},
					{
						Host:             "*.install.coder.app",
						IngressRuleValue: coderdIngressRule,
					},
				}
				require.Equal(t, expectedRules, ingress.Spec.Rules, "expected ingress spec to match")
			},
		},
		{
			Name: "devurl-suffix",
			ValuesFunc: func(v *CoderValues) {
				v.Ingress.Enable = pointer.Bool(true)
				v.Ingress.Host = pointer.String("install.coder.com")
				v.Coderd.DevURLsHost = pointer.String("*-dev.install.coder.app")
			},
			AssertFunc: func(t *testing.T, ingress *netv1.Ingress) {
				defaultAnnotations := map[string]string{
					"nginx.ingress.kubernetes.io/proxy-body-size": "0",
				}
				require.Equal(t, defaultAnnotations, ingress.Annotations)
				expectedRules := []netv1.IngressRule{
					{
						Host:             "install.coder.com",
						IngressRuleValue: coderdIngressRule,
					},
					{
						Host:             "*.install.coder.app",
						IngressRuleValue: coderdIngressRule,
					},
				}
				require.Equal(t, expectedRules, ingress.Spec.Rules, "expected ingress spec to match")
			},
		},
		{
			Name: "ingress-className",
			ValuesFunc: func(v *CoderValues) {
				v.Ingress.Enable = pointer.Bool(true)
				v.Ingress.ClassName = pointer.String("test")
			},
			AssertFunc: func(t *testing.T, ingress *netv1.Ingress) {
				require.Equal(t, pointer.String("test"), ingress.Spec.IngressClassName, "expected classname to match")
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.Name, func(t *testing.T) {
			// Clone the original values
			values := &CoderValues{}
			copier.Copy(values, chart.OriginalValues)

			// Run function to perform test-specific modifications of defaults
			objs := chart.MustRender(t, test.ValuesFunc)

			var found bool
			for _, obj := range objs {
				ingress, ok := obj.(*netv1.Ingress)
				if ok && ingress.Name == "coderd-ingress" {
					found = true
					test.AssertFunc(t, ingress)
					break
				}
			}
			require.True(t, found, "expected ingress in manifests")
		})
	}
}
