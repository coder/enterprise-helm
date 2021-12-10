# coder.workspaces.configMap defines configuration that is applied
# to user workspaces.
#
# TODO: rename Values.environments to Values.coderd.workspaces,
# once we verify that it won't hurt backward compatibility
{{- define "coder.workspaces.configMap" }}
{{- if .Values.environments.tolerations }}
---
apiVersion: v1
kind: ConfigMap
metadata:
  namespace: {{ .Release.Namespace | quote }}
  # TODO: change this to coderd, and store other settings in
  # the ConfigMap
  name: ce-environment-config
data:
  tolerations: {{ toJson .Values.environments.tolerations | b64enc | quote }}
{{- end }}
{{- end }}

# coder.workspaces.configMapEnv contains a POD_TOLERATIONS environment
# variable.
#
# coderd uses this environment variable to unmarshal pod toleration objects.
{{- define "coder.workspaces.configMapEnv" }}
{{- if (merge .Values dict | dig "environments" "tolerations" false) }}
- name: POD_TOLERATIONS
  value: {{ toJson .Values.environments.tolerations | b64enc | quote }}
{{- end }}
{{- if (merge .Values dict | dig "environments" "nodeSelector" false) }}
- name: POD_NODESELECTOR
  value: {{ toJson .Values.environments.nodeSelector | b64enc | quote }}
{{- end }}
{{- end }}
