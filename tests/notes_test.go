package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/chartutil"
	"k8s.io/utils/pointer"
)

// TestVersions ensures that a warning appears for versions
// that are incompatible.
func TestVersions(t *testing.T) {
	t.Parallel()

	chart := LoadChart(t)

	capab := chartutil.DefaultCapabilities.Copy()
	capab.KubeVersion.Version = "1.19.13-gke.1900"
	capab.KubeVersion.Major = "1"
	capab.KubeVersion.Minor = "19"

	warning := `======================= KUBERNETES SUPPORT =======================

NOTICE: Coder follows the Kubernetes upstream version support
policy, and the latest stable release version of Coder supports
the previous two minor releases as well as the current release of
Kubernetes at time of publication.

Your Kubernetes version is: 1.19.13-gke.1900
Coder 1.28.0 requires Kubernetes >= 1.21

Coder cannot provide any guarantees of compatibility nor technical
support for this version, in accordance with our support policy:
https://coder.com/docs/coder/latest/setup/kubernetes#supported-kubernetes-versions

======================= KUBERNETES SUPPORT =======================`

	notes, err := chart.RenderNotes(nil, nil, capab)
	require.NoError(t, err, "error rendering NOTES.txt")
	require.Equal(t, warning, notes, "warning should match expected")
}

// TestDeprecatedTrustProxyIP checks that the chart emits a warning when
// the deprecated coderd.trustProxyIP setting is set.
func TestDeprecatedTrustProxyIP(t *testing.T) {
	t.Parallel()

	chart := LoadChart(t)

	warning := `======================= DEPRECATION NOTICE =======================

WARNING: The coderd "trustProxyIP" setting is deprecated. Instead,
use the coderd "reverseProxy" setting to configure trusted headers
and origins.

See https://coder.com/docs/coder/latest/guides/deployments/proxy

======================= DEPRECATION NOTICE =======================`

	capab := chartutil.DefaultCapabilities.Copy()
	capab.KubeVersion.Version = "1.23.1-gke.500"
	capab.KubeVersion.Major = "1"
	capab.KubeVersion.Minor = "23"

	notes, err := chart.RenderNotes(func(cv *CoderValues) {
		cv.Coderd.TrustProxyIP = pointer.Bool(true)
	}, nil, capab)
	require.NoError(t, err, "error rendering NOTES.txt")
	require.Equal(t, warning, notes, "warning should match expected")
}
