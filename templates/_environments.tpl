{{/*
  coder.environments.saName defines the environments service account name.
*/}}
{{- define "coder.environments.saName" -}}
{{ .Release.Name }}-environments
{{- end}}

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
{{- end }}