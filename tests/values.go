package tests

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/xerrors"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/yaml"
)

var _ = fmt.Stringer(CoderValues{})

// CoderValues is a typed Go representation of Coder's values file,
// suitable for writing tests. This provides code completion for Go
// tests.
type CoderValues struct {
	Coderd *CoderdValues `json:"coderd" yaml:"coderd"`
}

// CoderdValues are values that apply to coderd.
type CoderdValues struct {
	Image              *string                   `json:"image" yaml:"image"`
	Replicas           *int                      `json:"replicas" yaml:"replicas"`
	ServiceSpec        *CoderdServiceSpecValues  `json:"serviceSpec" yaml:"serviceSpec"`
	PodSecurityContext *CoderdPodSecurityContext `json:"podSecurityContext" yaml:"podSecurityContext"`
	SecurityContext    *CoderdSecurityContext    `json:"securityContext" yaml:"securityContext"`
}

// PostgresValues are values that apply to postgres.
type PostgresValues struct {
	// TODO@jsjoeio
	// There is something called NamespaceDefault in the corev1 type/s
	// but I can't figure out how to import or use it.
	Default *PostgresDefaultValues `json:"default" yaml:"default"`
}

type PostgresDefaultValues struct {
	Resources corev1.ResourceList
}

type PostgresRequestsValues struct {
	CPU    string
	Memory string
}

type CoderdServiceSpecValues struct {
	Type                  *string `json:"type" yaml:"type"`
	ExternalTrafficPolicy *string `json:"externalTrafficPolicy" yaml:"externalTrafficPolicy"`
	LoadBalancerIP        *string `json:"loadBalancerIP" yaml:"loadBalancerIP"`
}

type CoderdPodSecurityContext struct {
	RunAsNonRoot *bool `json:"runAsNonRoot" yaml:"runAsNonRoot"`
	RunAsUser    *int  `json:"runAsUser" yaml:"runAsUser"`
}

type CoderdSecurityContext struct {
	ReadOnlyRootFilesystem   *bool `json:"readOnlyRootFilesystem" yaml:"readOnlyRootFilesystem"`
	AllowPrivilegeEscalation *bool `json:"allowPrivilegeEscalation" yaml:"allowPrivilegeEscalation"`
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

// RenderChart applies the CoderValues to the chart, and returns a list
// of Kubernetes runtime objects, or an error.
//
// values, options, and capabilities may be nil, in which case the
// function will simulate a fresh install to the "coder" namespace
// using the "coder" release, default values, and capabilities.
func RenderChart(chrt *chart.Chart, values *CoderValues, options *chartutil.ReleaseOptions, capabilities *chartutil.Capabilities) ([]runtime.Object, error) {
	vals, err := ConvertCoderValuesToMap(values)
	if err != nil {
		return nil, fmt.Errorf("failed to convert values to map: %w", err)
	}

	var opts chartutil.ReleaseOptions
	if options == nil {
		opts = chartutil.ReleaseOptions{
			Name:      "coder",
			Namespace: "coder",
			Revision:  1,
			IsInstall: true,
			IsUpgrade: false,
		}
	} else {
		opts = *options
	}

	if capabilities == nil {
		capabilities = chartutil.DefaultCapabilities.Copy()
	}

	vals, err = chartutil.ToRenderValues(chrt, vals, opts, capabilities)
	if err != nil {
		return nil, fmt.Errorf("failed to create render values: %w", err)
	}

	manifests, err := engine.Render(chrt, vals)
	if err != nil {
		return nil, fmt.Errorf("failed to render Chart: %w", err)
	}

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

// ReadValues reads the values.yaml from a file
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
