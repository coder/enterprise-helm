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

	tests := []struct {
		Name       string
		ValuesFunc func(v *CoderValues)
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
				pathTypePrefix := netv1.PathTypePrefix
				expectedRules := []netv1.IngressRule{
					{
						Host: "install.coder.com",
						IngressRuleValue: netv1.IngressRuleValue{
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
						},
					},
					{
						Host: "*.install.coder.app",
						IngressRuleValue: netv1.IngressRuleValue{
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
						},
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
				pathTypePrefix := netv1.PathTypePrefix
				expectedRules := []netv1.IngressRule{
					{
						Host: "install.coder.com",
						IngressRuleValue: netv1.IngressRuleValue{
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
						},
					},
					{
						Host: "*.install.coder.app",
						IngressRuleValue: netv1.IngressRuleValue{
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
						},
					},
				}
				require.Equal(t, expectedRules, ingress.Spec.Rules, "expected ingress spec to match")
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
			test.ValuesFunc(values)

			// Verify the results using AssertFunc
			objs, err := chart.Render(values, nil, nil)
			require.NoError(t, err, "chart render failed")

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
