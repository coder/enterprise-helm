{{/*
  coder.namespaceWhitelist.envSA contains a service account used by all
  environments in each namespace.
  LEGACY: .Values.namespaceWhitelist is deprecated. This exists to support
          existing environments.
*/}}
{{- define "coder.namespaceWhitelist.envSA" }}
{{- if .Values.namespaceWhitelist }}
{{- range .Values.namespaceWhitelist }}
---
apiVersion: v1
kind: ServiceAccount
metadata:
  namespace: {{ . | quote }}
  name: environments
{{- end }}
{{- end }}
{{- end }}
{{/* 
  coder.environments.configMap defines configuration that is applied
  to user environments.
*/}}
{{- define "coder.environments.configMap" }}
{{- if .Values.environments.tolerations }}
---
apiVersion: v1
kind: ConfigMap
metadata:
  namespace: {{ .Release.Namespace | quote }}
  name: ce-environment-config
data:
  tolerations: {{ toJson .Values.environments.tolerations | b64enc | quote }}
{{- end}}
{{- end}}
{{/*
  coder.environments.configMapEnv contains a POD_TOLERATIONS environment variable.
  ce-manager uses this environment variable to unmarshal pod toleration objects.
*/}}
{{- define "coder.environments.configMapEnv" }}
{{- if .Values.environments.tolerations }}
- name: POD_TOLERATIONS
  value: {{ toJson .Values.environments.tolerations | b64enc | quote }}
{{- end }}
{{- if .Values.environments.nodeSelector }}
- name: POD_NODESELECTOR
  value: {{ toJson .Values.environments.nodeSelector | b64enc | quote }}
{{- end }}
{{- end }}