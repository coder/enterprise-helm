{{/*
  coder.ingress.tls configures the tls settings
  for the default nginx ingress based on the
  values.yaml settings.
*/}}
{{- define "coder.ingress.tls" }}
{{- if .Values.ingress.tls.enable }}
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
