{{/*
  coder.cemanager.roles creates a role for every valid namespace to provide
  cemanager the necessary permissionse to administer environments.
*/}}
{{- define "coder.cemanager.roles" }}
{{- $namespaces := append .Values.namespaceWhitelist .Release.Namespace }}
{{- range $namespaces }}
---
# The roles the manager needs in the environments
# namespace in order to administer pods.
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: cemanager
  namespace: {{ . | quote }}
rules:
  - apiGroups:
      - "" # indicates the core API group
      - "apps"
    resources:
        # Every environment deployment has a service.
      - services
      - pods
        # Required for creating registry pull secrets.
      - secrets
        # Required for the code-server volume and the environments' home volumes.
      - persistentvolumeclaims
        # Required for legacy reasons.
      - statefulsets
        # Every environment is a deployment.
      - deployments
    verbs:
      - create
      - get
      - list
      - watch
      - update
      - patch
      - delete
      - deletecollection
  - apiGroups:
      - "" # indicates the core API group
    resources:
      - pods/exec
      - pods/log
      - events
      - configmaps
    verbs:
      - create
      - get
      - list
      - watch
  - apiGroups:
      - "apps"
    resources:
      # Required for legacy reasons.
      - statefulsets/scale
      # Auto-off requires being able to update the
      # replicas of a deployment from 1 to 0.
      - deployments/scale
    verbs:
      - update
  - apiGroups:
      - metrics.k8s.io
    resources:
      - pods
    # Necessary for gathering stats about an environment.
    verbs:
      - list
      - get
  - apiGroups:
      - networking.k8s.io
    resources:
        # Necessary for preventing inter-environment communication.
      - networkpolicies
    verbs:
      - create
      - delete
      - get
      - list
      - patch
{{- end }}
{{- end }}
{{/*
  coder.cemanager.rolebindings creates a
  role binding for every namespace in the list
  to allow cemanager permissions to create environments
  in the additional namespaces
*/}}
{{- define "coder.cemanager.rolebindings" }}
{{ $namespaces := append .Values.namespaceWhitelist .Release.Namespace }}
{{- range $namespaces }}
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  namespace: {{ . | quote }}
  name: cemanager
subjects:
  - kind: ServiceAccount
    name: cemanager
    namespace: {{ $.Release.Namespace | quote }}
roleRef:
  kind: Role
  name: cemanager
  apiGroup: rbac.authorization.k8s.io
{{- end }}
{{- end }}
