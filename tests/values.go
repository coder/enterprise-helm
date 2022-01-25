package tests

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/yaml"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"

	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/require"
	"golang.org/x/xerrors"
)

var _ = fmt.Stringer(CoderValues{})

// Chart wraps the default Helm chart, preserving default values.
//
// This technique has the side effect of requiring that all values
// be defined in the Values struct to behave correctly.
type Chart struct {
	// chart is the original Helm chart. Callers should not need
	// to access this directly.
	chart *chart.Chart

	// Metadata is the Helm chart Metadata field.
	Metadata *chart.Metadata

	// Templates for this chart.
	Templates []*chart.File

	// Files are other miscellaneous files included in the chart.
	Files []*chart.File

	// OriginalValues contains the original chart values. This
	// is intended to be read-only and should not be modified
	// by callers. Instead, modify the Values field.
	OriginalValues *CoderValues
}

// CoderValues is a typed Go representation of Coder's values
// file, suitable for writing tests.
//
// This technique provides code completion for Go tests, and has
// the side effect of requiring that all values be defined in the
// struct to behave correctly.
//
// TODO: generate these structs from a values.schema.json
type CoderValues struct {
	Certs    *CertsValues    `json:"certs" yaml:"certs"`
	Coderd   *CoderdValues   `json:"coderd" yaml:"coderd"`
	Envbox   *EnvboxValues   `json:"envbox" yaml:"envbox"`
	Ingress  *IngressValues  `json:"ingress" yaml:"ingress"`
	Logging  *LoggingValues  `json:"logging" yaml:"logging"`
	Metrics  *MetricsValues  `json:"metrics" yaml:"metrics"`
	Postgres *PostgresValues `json:"postgres" yaml:"postgres"`
	Services *ServicesValues `json:"services" yaml:"services"`
}

// CoderdValues reflect values from coderd.
type CoderdValues struct {
	Image                         *string                                    `json:"image" yaml:"image"`
	Replicas                      *int                                       `json:"replicas" yaml:"replicas"`
	ServiceSpec                   *CoderdServiceSpecValues                   `json:"serviceSpec" yaml:"serviceSpec"`
	ServiceNodePorts              *CoderdServiceNodePortsValues              `json:"serviceNodePorts" yaml:"serviceNodePorts"`
	TrustProxyIP                  *bool                                      `json:"trustProxyIP" yaml:"trustProxyIP"`
	DevURLsHost                   *string                                    `json:"devurlsHost" yaml:"devurlsHost"`
	TLS                           *CoderdTLSValues                           `json:"tls" yaml:"tls"`
	Satellite                     *CoderdSatelliteValues                     `json:"satellite" yaml:"satellite"`
	PodSecurityContext            *corev1.PodSecurityContext                 `json:"podSecurityContext" yaml:"podSecurityContext"`
	SecurityContext               *corev1.SecurityContext                    `json:"securityContext" yaml:"securityContext"`
	Resources                     *corev1.ResourceRequirements               `json:"resources" yaml:"resources"`
	BuiltinProviderServiceAccount *CoderdBuiltinProviderServiceAccountValues `json:"builtinProviderServiceAccount" yaml:"builtinProviderServiceAccount"`
	OIDC                          *CoderdOIDCValues                          `json:"oidc" yaml:"oidc"`
	SuperAdmin                    *CoderdSuperAdminValues                    `json:"superAdmin" yaml:"superAdmin"`
	Affinity                      *corev1.Affinity                           `json:"affinity" yaml:"affinity"`
	ExtraLabels                   map[string]string                          `json:"extraLabels" yaml:"extraLabels"`
	Proxy                         *CoderdProxyValues                         `json:"proxy" yaml:"proxy"`
	ReverseProxy                  *CoderdReverseProxyValues                  `json:"reverseProxy" yaml:"reverseProxy"`
	NetworkPolicy                 *CoderdNetworkPolicyValues                 `json:"networkPolicy" yaml:"networkPolicy"`
}

// CoderdServiceNodePortsValues reflect values from
// coderd.serviceNodePorts.
type CoderdServiceNodePortsValues struct {
	HTTP  *int32 `json:"http" yaml:"http"`
	HTTPS *int32 `json:"https" yaml:"https"`
}

// CoderdSuperAdminValues reflect values from
// coderd.superAdmin.
type CoderdSuperAdminValues struct {
	PasswordSecret *CoderdSuperAdminPasswordSecretValues `json:"passwordSecret" yaml:"passwordSecret"`
}

// CoderdSuperAdminPasswordSecretValues reflect values from
// coderd.superAdmin.passwordSecret.
type CoderdSuperAdminPasswordSecretValues struct {
	Name *string `json:"name" yaml:"name"`
	Key  *string `json:"key" yaml:"key"`
}

// CoderdTLSValues reflect values from coderd.tls.
type CoderdTLSValues struct {
	HostSecretName        *string `json:"hostSecretName" yaml:"hostSecretName"`
	DevURLsHostSecretName *string `json:"devurlsHostSecretName" yaml:"devurlsHostSecretName"`
}

// CoderdProxyValues reflect values from coderd.proxy.
type CoderdProxyValues struct {
	HTTP   *string `json:"http" yaml:"http"`
	HTTPS  *string `json:"https" yaml:"https"`
	Exempt *string `json:"exempt" yaml:"exempt"`
}

// CoderdReverseProxyValues reflect values from coderd.reverseProxy
type CoderdReverseProxyValues struct {
	TrustedOrigins []string `json:"trustedOrigins" yaml:"trustedOrigins"`
	Headers        []string `json:"headers" yaml:"headers"`
}

// CoderdBuiltinProviderServiceAccountValues reflect values from
// coderd.builtinProviderServiceAccount.
type CoderdBuiltinProviderServiceAccountValues struct {
	// Labels is the same type as metav1.ObjectMeta.Labels
	Labels map[string]string `json:"labels" yaml:"labels"`
	// Annotations is the same type as metav1.ObjectMeta.Annotations
	Annotations map[string]string `json:"annotations" yaml:"annotations"`
}

// CoderdNetworkPolicyValues reflect values from coderd.networkPolicy.
type CoderdNetworkPolicyValues struct {
	Enable *bool `json:"enable" yaml:"enable"`
}

// CoderdOIDCValues reflect values from coderd.oidc.
type CoderdOIDCValues struct {
	EnableRefresh   *bool             `json:"enableRefresh" yaml:"enableRefresh"`
	RedirectOptions map[string]string `json:"redirectOptions" yaml:"redirectOptions"`
}

// CoderdSatelliteValues reflect values from coderd.satellite.
type CoderdSatelliteValues struct {
	Enable     *bool   `json:"enable" yaml:"enable"`
	AccessURL  *string `json:"accessURL" yaml:"accessURL"`
	PrimaryURL *string `json:"primaryURL" yaml:"primaryURL"`
}

// CoderdServiceSpecValues reflect values from coderd.serviceSpec.
type CoderdServiceSpecValues struct {
	Type                     *string                                  `json:"type" yaml:"type"`
	ExternalTrafficPolicy    *corev1.ServiceExternalTrafficPolicyType `json:"externalTrafficPolicy" yaml:"externalTrafficPolicy"`
	LoadBalancerIP           *string                                  `json:"loadBalancerIP" yaml:"loadBalancerIP"`
	LoadBalancerSourceRanges *[]string                                `json:"loadBalancerSourceRanges" yaml:"loadBalancerSourceRanges"`
}

// EnvboxValues reflect values from envbox.
type EnvboxValues struct {
	Image *string `json:"image" yaml:"image"`
}

// IngressValues reflect values from ingress.
type IngressValues struct {
	Enable      *bool             `json:"enable" yaml:"enable"`
	Host        *string           `json:"host" yaml:"host"`
	Annotations map[string]string `json:"annotations" yaml:"annotations"`
	TLS         *IngressTLSValues `json:"tls" yaml:"tls"`
}

// IngressTLSValues reflect values from ingress.tls.
type IngressTLSValues struct {
	Enable *bool `json:"enable" yaml:"enable"`
}

// LoggingValues reflect values from logging.
type LoggingValues struct {
	Human       *string              `json:"human" yaml:"human"`
	Stackdriver *string              `json:"stackdriver" yaml:"stackdriver"`
	JSON        *string              `json:"json" yaml:"json"`
	Splunk      *LoggingSplunkValues `json:"splunk" yaml:"splunk"`
}

// LoggingSplunkValues reflect values from logging.splunk.
type LoggingSplunkValues struct {
	URL     *string `json:"url" yaml:"url"`
	Token   *string `json:"token" yaml:"token"`
	Channel *string `json:"channel" yaml:"channel"`
}

// MetricsValues reflect values from metrics.
type MetricsValues struct {
	AmplitudeKey *string `json:"amplitudeKey" yaml:"amplitudeKey"`
}

// CertsValues reflect the values from certs.
type CertsValues struct {
	Secret *CertsSecretValues `json:"secret" yaml:"secret"`
}

// CertsSecretValues reflect the values from certs.secret.
type CertsSecretValues struct {
	Name *string `json:"name" yaml:"name"`
	Key  *string `json:"key" yaml:"key"`
}

// PostgresValues reflect the values from postgres.
type PostgresValues struct {
	Host           *string                `json:"host" yaml:"host"`
	Port           *string                `json:"port" yaml:"port"`
	User           *string                `json:"user" yaml:"user"`
	SSLMode        *string                `json:"sslMode" yaml:"sslMode"`
	Database       *string                `json:"database" yaml:"database"`
	PasswordSecret *string                `json:"passwordSecret" yaml:"passwordSecret"`
	Default        *PostgresDefaultValues `json:"default" yaml:"default"`
	SSL            *PostgresSSLValues     `json:"ssl" yaml:"ssl"`
	Connector      *string                `json:"connector" yaml:"connector"`
}

type PostgresSSLValues struct {
	CertSecret     *CertsSecretValues `json:"certSecret" yaml:"certSecret"`
	KeySecret      *CertsSecretValues `json:"keySecret" yaml:"keySecret"`
	RootCertSecret *CertsSecretValues `json:"rootCertSecret" yaml:"rootCertSecret"`
}

// PostgresDefaultValues reflect the values from
// postgres.default.
type PostgresDefaultValues struct {
	Enable           *bool                               `json:"enable" yaml:"enable"`
	Image            *string                             `json:"image" yaml:"image"`
	StorageClassName *string                             `json:"storageClassName" yaml:"storageClassName"`
	Resources        *corev1.ResourceRequirements        `json:"resources" yaml:"resources"`
	NetworkPolicy    *PostgresDefaultNetworkPolicyValues `json:"networkPolicy" yaml:"networkPolicy"`
}

// PostgresDefaultNetworkPolicyValues reflect values from
// postgres.default.networkPolicy.
type PostgresDefaultNetworkPolicyValues struct {
	Enable *bool `json:"enable" yaml:"enable"`
}

// ServicesValues reflect the values from services.
type ServicesValues struct {
	Annotations         map[string]string    `json:"annotations" yaml:"annotations"`
	ClusterDomainSuffix *string              `json:"clusterDomainSuffix" yaml:"clusterDomainSuffix"`
	Tolerations         *[]corev1.Toleration `json:"tolerations" yaml:"tolerations"`
	NodeSelector        map[string]string    `json:"nodeSelector" yaml:"nodeSelector"`
	Type                *corev1.ServiceType  `json:"type" yaml:"type"`
}

// String returns the string representation of the values.
func (v CoderValues) String() string {
	var sb strings.Builder

	encoder := json.NewEncoder(&sb)
	encoder.SetIndent("", "    ")
	err := encoder.Encode(&v)
	if err != nil {
		panic(fmt.Sprintf("failed encode: %v", err))
	}

	return sb.String()
}

// ConvertCoderValuesToMap returns the CoderValues struct encoded as
// a map[string]interface{} compatible with Helm.
func ConvertCoderValuesToMap(v *CoderValues) (map[string]interface{}, error) {
	var buf bytes.Buffer

	// Marshal the values to a buffer
	err := json.NewEncoder(&buf).Encode(v)
	if err != nil {
		return nil, xerrors.Errorf("marshal json: %w", err)
	}

	// Unmarshal the values into a map
	valuesMap := map[string]interface{}{}
	err = json.NewDecoder(&buf).Decode(&valuesMap)
	if err != nil {
		return nil, xerrors.Errorf("unmarshal json: %w", err)
	}

	return valuesMap, nil
}

// ConvertMapToCoderValues takes a map[string]interface{} as specified
// in Helm values, and converts it to a map.
//
// If strict is true, then it will error in the event that the values
// file contains unknown fields.
func ConvertMapToCoderValues(v map[string]interface{}, strict bool) (*CoderValues, error) {
	var buf bytes.Buffer

	// Marshal the values to a buffer
	err := json.NewEncoder(&buf).Encode(v)
	if err != nil {
		return nil, xerrors.Errorf("marshal json: %w", err)
	}

	// Unmarshal the values into a map
	values := &CoderValues{}
	decoder := json.NewDecoder(&buf)
	if strict {
		decoder.DisallowUnknownFields()
	}
	err = decoder.Decode(values)
	if err != nil {
		return nil, xerrors.Errorf("unmarshal json: %w", err)
	}

	return values, nil
}

// LoadChart is a utility function that loads the chart from the
// unpacked source directory.
func LoadChart(t testing.TB) *Chart {
	chart, err := loader.LoadDir("..")
	require.NoError(t, err, "loaded chart successfully")
	require.NotNil(t, chart, "chart must be non-nil")
	require.True(t, chart.IsRoot(), "chart must be a root chart")

	// Load original values so that users can override them.
	originalValues, err := ConvertMapToCoderValues(chart.Values, true)
	require.NoError(t, err, "error parsing original values")

	return &Chart{
		chart:          chart,
		Metadata:       chart.Metadata,
		Files:          chart.Files,
		Templates:      chart.Templates,
		OriginalValues: originalValues,
	}
}

// Name returns the name of the chart.
func (c *Chart) Name() string {
	return c.chart.Name()
}

// IsRoot is true if this is not a subchart and has no parents.
func (c *Chart) IsRoot() bool {
	return c.chart.IsRoot()
}

// AppVersion returns the chart appversion.
func (c *Chart) AppVersion() string {
	return c.chart.AppVersion()
}

// Validate checks that the chart metadata is valid.
func (c *Chart) Validate() error {
	return c.chart.Validate()
}

// Render creates a copy of the default chart values, runs fn to
// modify those values, then applies those values to the chart,
// returning a list of Kubernetes runtime objects, or an error.
//
// values, options, and capabilities may be nil, in which case the
// function will simulate a fresh install to the "coder" namespace
// using the "coder" release, default values, and capabilities.
func (c *Chart) Render(fn func(*CoderValues), options *chartutil.ReleaseOptions, capabilities *chartutil.Capabilities) ([]runtime.Object, error) {
	var opts chartutil.ReleaseOptions
	if options == nil {
		opts = DefaultReleaseOptions()
	} else {
		opts = *options
	}

	if capabilities == nil {
		capabilities = chartutil.DefaultCapabilities.Copy()
	}

	values := c.OriginalValues
	if fn != nil {
		values = &CoderValues{}
		copier.CopyWithOption(values, c.OriginalValues, copier.Option{
			DeepCopy: true,
		})

		fn(values)
	}

	vals, err := ConvertCoderValuesToMap(values)
	if err != nil {
		return nil, fmt.Errorf("failed to convert CoderValues to map: %w", err)
	}

	vals, err = chartutil.ToRenderValues(c.chart, vals, opts, capabilities)
	if err != nil {
		return nil, fmt.Errorf("failed to create render values: %w", err)
	}

	manifests, err := engine.Render(c.chart, vals)
	if err != nil {
		return nil, fmt.Errorf("failed to render Chart: %w", err)
	}

	objs, err := LoadObjectsFromManifests(manifests)
	if err != nil {
		return nil, fmt.Errorf("failed to load objects: %w", err)
	}

	return objs, nil
}

// MustRender renders a chart or fails the test. Use `fn` to modify the default
// chart values.
func (c *Chart) MustRender(t testing.TB, fn func(*CoderValues)) []runtime.Object {
	objs, err := c.Render(fn, nil, nil)
	require.NoError(t, err, "render chart")

	return objs
}

func DefaultReleaseOptions() chartutil.ReleaseOptions {
	return chartutil.ReleaseOptions{
		Name:      "coder",
		Namespace: "coder",
		Revision:  1,
		IsInstall: true,
		IsUpgrade: false,
	}
}

func LoadObjectsFromManifests(manifests map[string]string) ([]runtime.Object, error) {
	deserializer := NewDeserializer()

	var objs []runtime.Object

	// Helm returns a map of rendered files and contents
	for file, manifest := range manifests {
		reader := yaml.NewYAMLReader(bufio.NewReader(strings.NewReader(manifest)))
		// Split files into individual document chunks, then pass through
		// the deserializer
		for {
			document, err := reader.Read()
			if err != nil {
				// If we get an EOF, we've finished processing this file
				if err == io.EOF {
					break
				}
				return nil, fmt.Errorf("failed to read %q: %w", file, err)
			}

			// Exit the inner loop if we encounter an EOF
			if document == nil {
				break
			}

			// Skip empty documents
			if document[0] == '\n' {
				continue
			}

			obj, _, err := deserializer.Decode(document, nil, nil)
			if err != nil {
				return nil, fmt.Errorf("error deserializing %q: %w", file, err)
			}

			objs = append(objs, obj)
		}
	}

	return objs, nil
}

// NewDeserializer creates a UniversalDeserializer using the scheme
// registered by NewScheme.
func NewDeserializer() runtime.Decoder {
	scheme := NewScheme()
	deserializer := serializer.NewCodecFactory(scheme).UniversalDeserializer()
	return deserializer
}

// NewScheme creates a runtime.Scheme and registers all known types.
// Any Kubernetes types that we use must be registered here, or else
// there will be a deserialization error.
func NewScheme() *runtime.Scheme {
	scheme := runtime.NewScheme()
	if err := appsv1.AddToScheme(scheme); err != nil {
		panic(fmt.Sprintf("failed to add appsv1 scheme: %v", err))
	}
	if err := corev1.AddToScheme(scheme); err != nil {
		panic(fmt.Sprintf("failed to add corev1 scheme: %v", err))
	}
	if err := networkingv1.AddToScheme(scheme); err != nil {
		panic(fmt.Sprintf("failed to add networkingv1 scheme: %v", err))
	}
	if err := rbacv1.AddToScheme(scheme); err != nil {
		panic(fmt.Sprintf("failed to add rbacv1 scheme: %v", err))
	}
	return scheme
}

// ReadValues reads the values.yaml from a file.
func ReadValues(path string) (*CoderValues, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open %q: %w", path, err)
	}

	var values CoderValues
	decoder := yaml.NewYAMLToJSONDecoder(file)
	err = decoder.Decode(&values)
	if err != nil {
		return nil, fmt.Errorf("error decoding yaml %q: %w", path, err)
	}

	return &values, nil
}

// ReadValuesFileAsMap reads the values.yaml from a file.
func ReadValuesFileAsMap(path string) (map[string]interface{}, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open %q: %w", path, err)
	}

	var values map[string]interface{}
	decoder := yaml.NewYAMLToJSONDecoder(file)
	err = decoder.Decode(&values)
	if err != nil {
		return nil, fmt.Errorf("error decoding yaml %q: %w", path, err)
	}

	return values, nil
}
