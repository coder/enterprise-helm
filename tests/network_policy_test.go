package tests

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"

	"github.com/stretchr/testify/require"
)

// TestNetworkPolicyCoder tests that chart creates a network policy for the control plane.
func TestNetworkPolicyCoder(t *testing.T) {
	t.Parallel()

	chart := LoadChart(t)

	tests := []struct {
		Name       string
		ValuesFunc func(cv *CoderValues)
		// ExpectCoderPolicy is true if we expect a network policy for coderd
		// or the satellite
		ExpectCoderPolicy bool
		// ExpectDatabasePolicy is true if we expect a network policy for
		// the built-in database
		ExpectDatabasePolicy bool
	}{
		{
			Name:       "default-primary",
			ValuesFunc: nil,
			// The default install has coderd and a built-in database, so both
			// policies should be set
			ExpectCoderPolicy:    true,
			ExpectDatabasePolicy: true,
		},
		{
			Name: "default-satellite",
			ValuesFunc: func(cv *CoderValues) {
				cv.Coderd.Satellite.Enable = pointer.Bool(true)
			},
			ExpectCoderPolicy: true,
			// For a satellite, the built-in database is not deployed, so there
			// should not be network policies.
			ExpectDatabasePolicy: false,
		},
		{
			Name: "no-policies",
			ValuesFunc: func(cv *CoderValues) {
				cv.Coderd.NetworkPolicy.Enable = pointer.Bool(false)
				cv.Postgres.Default.NetworkPolicy.Enable = pointer.Bool(false)
			},
			ExpectCoderPolicy:    false,
			ExpectDatabasePolicy: false,
		},
		{
			Name: "no-builtin-database-policy",
			ValuesFunc: func(cv *CoderValues) {
				cv.Postgres.Default.Enable = pointer.Bool(false)
				cv.Postgres.Default.NetworkPolicy.Enable = pointer.Bool(true)
			},
			// If we're not using the built-in database, then the corresponding
			// policy should not exist, even if the user asks for it
			ExpectCoderPolicy:    true,
			ExpectDatabasePolicy: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()

			objs := chart.MustRender(t, test.ValuesFunc)

			policy, exist := FindNetworkPolicy(objs, "coderd")
			require.Equal(t, test.ExpectCoderPolicy, exist, "coderd network policy")
			if test.ExpectCoderPolicy {
				require.Contains(t, policy.Spec.PolicyTypes, networkingv1.PolicyTypeIngress, "expected to restrict ingress")
				require.NotContains(t, policy.Spec.PolicyTypes, networkingv1.PolicyTypeEgress, "expected all egress to be allowed")
				require.Empty(t, policy.Spec.Egress, "expected empty egress rules")
				protocolTCP := corev1.ProtocolTCP
				expectedRules := []networkingv1.NetworkPolicyIngressRule{
					{
						From: []networkingv1.NetworkPolicyPeer{},
						Ports: []networkingv1.NetworkPolicyPort{
							{
								Protocol: &protocolTCP,
								Port: &intstr.IntOrString{
									Type:   intstr.Int,
									IntVal: 8080,
								},
							},
							{
								Protocol: &protocolTCP,
								Port: &intstr.IntOrString{
									Type:   intstr.Int,
									IntVal: 8443,
								},
							},
						},
					},
				}
				require.Equal(t, expectedRules, policy.Spec.Ingress, "expected ingress rules to match")
			}

			policy, exist = FindNetworkPolicy(objs, "timescale")
			require.Equal(t, test.ExpectDatabasePolicy, exist, "timescale network policy")
			if test.ExpectDatabasePolicy {
				require.Contains(t, policy.Spec.PolicyTypes, networkingv1.PolicyTypeIngress, "expected to restrict ingress")
				require.Contains(t, policy.Spec.PolicyTypes, networkingv1.PolicyTypeEgress, "expected to restrict egress")
				require.Empty(t, policy.Spec.Egress, "expected empty egress rules")
				protocolTCP := corev1.ProtocolTCP

				podSelector := &metav1.LabelSelector{}
				metav1.AddLabelToSelector(podSelector, "app", "coderd")

				expectedRules := []networkingv1.NetworkPolicyIngressRule{
					{
						From: []networkingv1.NetworkPolicyPeer{
							{
								PodSelector: podSelector,
							},
						},
						Ports: []networkingv1.NetworkPolicyPort{
							{
								Protocol: &protocolTCP,
								Port: &intstr.IntOrString{
									Type:   intstr.Int,
									IntVal: 5432,
								},
							},
						},
					},
				}
				require.Equal(t, expectedRules, policy.Spec.Ingress, "expected ingress rules to match")
			}
		})
	}
}
