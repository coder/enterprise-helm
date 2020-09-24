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
  coder.environments.configMapEnv adds an environment variable with the name 
  of the config map holding additional user environment configuration.
*/}}
{{- define "coder.environments.configMapEnv" }}
{{- if .Values.serviceTolerations }}
- name: ENVIRONMENT_CONFIG_MAP
  value: "ce-environment-config"
{{- end }}
{{- end }}
