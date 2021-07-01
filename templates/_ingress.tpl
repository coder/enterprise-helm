{{/*
  coder.ingress.tls configures the tls settings
  for the default nginx ingress based on the
  values.yaml settings.
*/}}
{{- define "coder.ingress.tls" }}
{{- if (merge .Values dict | dig "ingress" "tls" "enable" false) }}
  tls:
    {{- if and .Values.ingress.host .Values.ingress.tls.hostSecretName }}
    - hosts:
      - {{ .Values.ingress.host | quote }}
      secretName: {{ .Values.ingress.tls.hostSecretName }}
    {{- end }}
    {{- if and .Values.devurls.host .Values.ingress.tls.devurlsHostSecretName }}
    - hosts:
      - {{ .Values.devurls.host | quote }}
      secretName: {{ .Values.ingress.tls.devurlsHostSecretName }}
    {{- end }}
{{- end }}
{{- end }}

{{/* */}}
{{- define "coder.hasNginxIngress" }}
{{- if (lookup "v1" "Service" .Release.Namespace "ingress-nginx") -}}
true
{{- else if .Values.envproxy -}}
true
{{- else -}}
false
{{- end }}
{{- end }}

{{- define "coder.useServiceNext" }}
{{- if eq (merge .Values dict | dig "coderd" "serviceNext" false) true -}}
true
{{- else if eq (merge .Values dict | dig "ingress" "useDefault" true) false -}}
false
{{- else if eq (include "coder.hasNginxIngress" .) "false" -}}
true
{{- else -}}
false
{{- end }}
{{- end }}
