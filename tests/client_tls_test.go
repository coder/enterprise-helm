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
		certSecretKey  = "cert"
		certMountPath  = "/etc/ssl/certs/client/cert"
		keyVolName     = "clientkey"
		keySecretName  = "coder-tls-key"
		keySecretKey   = "cert-key"
		keyMountPath   = "/etc/ssl/certs/client/key"
	)

	objs := chart.MustRender(t, func(cv *CoderValues) {
		cv.Coderd.ClientTls = &CoderdClientTlsValues{
			CertSecret: &CertsSecretValues{
				Name: &certSecretName,
				Key:  &certSecretKey,
			},
			KeySecret: &CertsSecretValues{
				Name: &keySecretName,
				Key:  &keySecretKey,
			},
		}
	})

	coderd := MustFindDeployment(t, objs, "coderd")
	require.Len(t, coderd.Spec.Template.Spec.Containers, 1)
	coderdCtr := coderd.Spec.Template.Spec.Containers[0]
	envVars := EnvVarsAsMap(coderdCtr.Env)

	// Assert volumes and volume mounts for both the cert and key.
	cases := []struct {
		volName    string
		secretName string
		secretKey  string
		mountPath  string
		envName    string
	}{
		{
			volName:    certVolName,
			secretName: certSecretName,
			secretKey:  certSecretKey,
			mountPath:  certMountPath,
			envName:    "SSL_CLIENT_CERT_FILE",
		},
		{
			volName:    keyVolName,
			secretName: keySecretName,
			secretKey:  keySecretKey,
			mountPath:  keyMountPath,
			envName:    "SSL_CLIENT_KEY_FILE",
		},
	}
	for _, c := range cases {
		AssertVolume(t, coderd.Spec.Template.Spec.Volumes, c.volName, func(t testing.TB, v v1.Volume) {
			require.NotNil(t, v.Secret)
			assert.Equal(t, c.secretName, v.Secret.SecretName)
		})

		AssertVolumeMount(t, coderdCtr.VolumeMounts, c.volName, func(t testing.TB, v v1.VolumeMount) {
			assert.Equal(t, c.mountPath, v.MountPath)
			assert.True(t, v.ReadOnly)
		})

		assert.Equal(t, envVars[c.envName], filepath.Join(c.mountPath, c.secretKey))
	}
}
