package tests

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"
)

// TestExamples loads the example values files and performs
// some basic checks.
func TestExamples(t *testing.T) {
	t.Parallel()

	chart := LoadChart(t)

	exampleIngress, err := ReadValuesFileAsMap("../examples/ingress/ingress.values.yaml")
	require.NoError(t, err, "failed to load ingress example values")

	exampleOffline, err := ReadValuesFileAsMap("../examples/offline/offline.values.yaml")
	require.NoError(t, err, "failed to load offline example values")

	exampleOpenShift, err := ReadValuesFileAsMap("../examples/openshift/openshift.values.yaml")
	require.NoError(t, err, "failed to load OpenShift example values")

	exampleKind, err := ReadValuesFileAsMap("../examples/kind/kind.values.yaml")
	require.NoError(t, err, "failed to load Kind example values")

	tests := []struct {
		Name                     string
		Values                   map[string]interface{}
		PodSecurityContext       *corev1.PodSecurityContext
		ContainerSecurityContext *corev1.SecurityContext
		ServiceType              corev1.ServiceType
		CoderdImageRef           string
	}{
		{
			Name:        "default",
			Values:      nil,
			ServiceType: corev1.ServiceTypeLoadBalancer,
		},
		{
			Name:        "ingress",
			Values:      exampleIngress,
			ServiceType: corev1.ServiceTypeClusterIP,
		},
		{
			Name:           "offline",
			Values:         exampleOffline,
			ServiceType:    corev1.ServiceTypeLoadBalancer,
			CoderdImageRef: "us-docker.pkg.dev/airgap-project/test/coder-service:1.25.0",
		},
		{
			Name:   "openshift",
			Values: exampleOpenShift,
			PodSecurityContext: &corev1.PodSecurityContext{
				RunAsUser:      nil,
				RunAsGroup:     nil,
				RunAsNonRoot:   pointer.Bool(true),
				SeccompProfile: nil,
			},
			ContainerSecurityContext: &corev1.SecurityContext{
				RunAsUser:                nil,
				RunAsGroup:               nil,
				RunAsNonRoot:             pointer.Bool(true),
				Capabilities:             nil,
				Privileged:               nil,
				SELinuxOptions:           nil,
				WindowsOptions:           nil,
				ReadOnlyRootFilesystem:   pointer.Bool(true),
				AllowPrivilegeEscalation: pointer.Bool(false),
				ProcMount:                nil,
				SeccompProfile:           nil,
			},
			ServiceType: corev1.ServiceTypeClusterIP,
		},
		{
			Name:   "kind",
			Values: exampleKind,
			PodSecurityContext: &corev1.PodSecurityContext{
				RunAsUser:    pointer.Int64(1000),
				RunAsGroup:   pointer.Int64(1000),
				RunAsNonRoot: pointer.Bool(true),
				SeccompProfile: &corev1.SeccompProfile{
					Type: corev1.SeccompProfileTypeRuntimeDefault,
				},
			},
			ContainerSecurityContext: &corev1.SecurityContext{
				RunAsUser:                pointer.Int64(1000),
				RunAsGroup:               pointer.Int64(1000),
				RunAsNonRoot:             pointer.Bool(true),
				Capabilities:             nil,
				Privileged:               nil,
				SELinuxOptions:           nil,
				WindowsOptions:           nil,
				ReadOnlyRootFilesystem:   pointer.Bool(true),
				AllowPrivilegeEscalation: pointer.Bool(false),
				ProcMount:                nil,
				SeccompProfile: &corev1.SeccompProfile{
					Type:             corev1.SeccompProfileTypeRuntimeDefault,
					LocalhostProfile: nil,
				},
			},
			ServiceType: corev1.ServiceTypeClusterIP,
		},
	}

	var (
		defaultPsp = &corev1.PodSecurityContext{
			RunAsUser:    pointer.Int64(1000),
			RunAsGroup:   pointer.Int64(1000),
			RunAsNonRoot: pointer.Bool(true),
			SeccompProfile: &corev1.SeccompProfile{
				Type: corev1.SeccompProfileTypeRuntimeDefault,
			},
		}

		defaultCsc = &corev1.SecurityContext{
			RunAsUser:                pointer.Int64(1000),
			RunAsGroup:               pointer.Int64(1000),
			RunAsNonRoot:             pointer.Bool(true),
			ReadOnlyRootFilesystem:   pointer.Bool(true),
			AllowPrivilegeEscalation: pointer.Bool(false),
			SeccompProfile: &corev1.SeccompProfile{
				Type: corev1.SeccompProfileTypeRuntimeDefault,
			},
		}
	)

	for _, test := range tests {
		test := test

		if test.PodSecurityContext == nil {
			test.PodSecurityContext = defaultPsp
		}
		if test.ContainerSecurityContext == nil {
			test.ContainerSecurityContext = defaultCsc
		}

		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()

			values, err := chartutil.ToRenderValues(chart.chart, test.Values, DefaultReleaseOptions(), chartutil.DefaultCapabilities.Copy())
			require.NoError(t, err, "failed to generate render values")

			manifests, err := engine.Render(chart.chart, values)
			require.NoError(t, err, "failed to render chart")

			// As a special case, ignore any .txt files (e.g. NOTES.txt)
			for key := range manifests {
				if filepath.Ext(key) == ".txt" {
					delete(manifests, key)
				}
			}

			objs, err := LoadObjectsFromManifests(manifests)
			require.NoError(t, err, "failed to convert manifests to objects")

			// Find the coderd Deployment
			coderd := MustFindDeployment(t, objs, "coderd")
			assert.Equal(t, test.PodSecurityContext, coderd.Spec.Template.Spec.SecurityContext,
				"expected matching pod securityContext",
			)
			require.Len(t, coderd.Spec.Template.Spec.Containers, 1,
				"expected one container",
			)
			if test.CoderdImageRef != "" {
				require.Equal(t, test.CoderdImageRef, coderd.Spec.Template.Spec.Containers[0].Image, "expected image ref to match")
			} else {
				imageRef := fmt.Sprintf("docker.io/coderenvs/coder-service:%s", chart.Metadata.AppVersion)
				require.Equal(t, imageRef, coderd.Spec.Template.Spec.Containers[0].Image, "expected image ref to be default")
			}
			assert.Equal(t, test.ContainerSecurityContext, coderd.Spec.Template.Spec.Containers[0].SecurityContext,
				"expected matching container securityContext",
			)

			service := MustFindService(t, objs, "coderd")
			assert.Equal(t, test.ServiceType, service.Spec.Type, "service type should match")
			switch test.ServiceType {
			case corev1.ServiceTypeLoadBalancer:
				assert.Empty(t, service.Spec.ExternalName, "external name should not be set")
			case corev1.ServiceTypeClusterIP:
				assert.Empty(t, service.Spec.ExternalName, "external name should not be set")
				assert.Nil(t, service.Spec.LoadBalancerClass, "loadBalancerClass should not be set")
				assert.Empty(t, service.Spec.LoadBalancerSourceRanges, "loadBalancerSourceRanges should not be set")
				assert.Empty(t, service.Spec.ExternalTrafficPolicy, "externalTrafficPolicy should not be set")
			}
		})
	}
}
