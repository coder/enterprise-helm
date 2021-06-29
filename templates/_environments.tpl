{{/* 
  coder.environments.configMap defines configuration that is applied
  to user environments.
*/}}
{{- define "coder.environments.configMap" }}
{{- if (merge .Values dict | dig "environments" "tolerations" false) }}
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
{{- if (merge .Values dict | dig "environments" "tolerations" false) }}
- name: POD_TOLERATIONS
  value: {{ toJson .Values.environments.tolerations | b64enc | quote }}
{{- end }}
{{- if (merge .Values dict | dig "environments" "nodeSelector" false) }}
- name: POD_NODESELECTOR
  value: {{ toJson .Values.environments.nodeSelector | b64enc | quote }}
{{- end }}
{{- end }}