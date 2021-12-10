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
    {{- if .Values.devurls }}
    {{- if and .Values.devurls.host .Values.ingress.tls.devurlsHostSecretName }}
    - hosts:
      - {{ .Values.coderd.devurlsHost }}
      secretName: {{ .Values.ingress.tls.devurlsHostSecretName }}
    {{- end }}
    {{- end }}
{{- end }}
{{- end }}
