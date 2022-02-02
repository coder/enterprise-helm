package tests

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
)

// ensures client TLS certs are correctly applied to coderd
func Test_ClientTLS(t *testing.T) {
	t.Parallel()

	var (
		chart = LoadChart(t)

		certVolName    = "clientcert"
		certSecretName = "coder-tls-cert"
		certMountPath  = "/etc/ssl/certs/client"
	)

	objs := chart.MustRender(t, func(cv *CoderValues) {
		cv.Coderd.ClientTLS = &CoderdClientTLSValues{
			SecretName: &certSecretName,
		}
	})

	coderd := MustFindDeployment(t, objs, "coderd")
	require.Len(t, coderd.Spec.Template.Spec.Containers, 1)
	coderdCtr := coderd.Spec.Template.Spec.Containers[0]
	envVars := EnvVarsAsMap(coderdCtr.Env)

	// Assert volumes and volume mounts for both the cert and key.
	AssertVolume(t, coderd.Spec.Template.Spec.Volumes, certVolName, func(t testing.TB, v v1.Volume) {
		require.NotNil(t, v.Secret)
		assert.Equal(t, certSecretName, v.Secret.SecretName)
	})

	AssertVolumeMount(t, coderdCtr.VolumeMounts, certVolName, func(t testing.TB, v v1.VolumeMount) {
		assert.Equal(t, certMountPath, v.MountPath)
		assert.True(t, v.ReadOnly)
	})

	assert.Equal(t, envVars["SSL_CLIENT_CERT_FILE"], filepath.Join(certMountPath, "tls.crt"))
	assert.Equal(t, envVars["SSL_CLIENT_KEY_FILE"], filepath.Join(certMountPath, "tls.key"))
}
