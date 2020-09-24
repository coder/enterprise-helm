{{/*
  coder.envproxy.roles creates a role for every valid namespace to provide
  envproxy permissions to administer environments.
*/}}
{{- define "coder.envproxy.roles" }}
{{ $namespaces := append .Values.namespaceWhitelist .Release.Namespace }}
{{- range $namespaces }}
---
# The envproxy role lists the permissions the envproxy requires.
# It conditionally sets the pod security policy based on the
# helm value.
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: envproxy
  namespace: {{ . | quote }}
rules:
  - apiGroups:
      - ""
      - "apps"
    resources:
      - services
      - pods
      - pods/exec
      - secrets
      - deployments
      - configmaps
    verbs:
      - create
      - list
      - get
      - watch
  - apiGroups:
      - ""
    resources:
      - secrets
    verbs:
      - create
  - apiGroups:
      - ""
    resources:
      - pods/exec
    verbs:
      - create
  - apiGroups:
      - metrics.k8s.io
    resources:
      - pods
    # Necessary for gathering stats about an environment.
    verbs:
      - list
      - get
{{- end }}
{{- end }}
{{/*
  coder.envproxy.rolebindings creates a role 
  binding for every valid namespace to provide envproxy
  permissions to administer environments.
*/}}
{{- define "coder.envproxy.rolebindings" }}
{{ $namespaces := append .Values.namespaceWhitelist .Release.Namespace }}
{{- range $namespaces }}
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  namespace: {{ . | quote }}
  name: envproxy
subjects:
  - kind: ServiceAccount
    name: envproxy
    namespace: {{ $.Release.Namespace | quote }}
roleRef:
  kind: Role
  name: envproxy
  apiGroup: rbac.authorization.k8s.io
{{- end }}
{{- end }}
