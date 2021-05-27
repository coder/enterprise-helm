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