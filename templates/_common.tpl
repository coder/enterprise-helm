{{/*
  coder.devurls.hostEnv adds an environment variable indicating
  the dev URL host.
*/}}
{{- define "coder.devurls.hostEnv" }}
{{- if ne .Values.devurls.host "" }}
- name: "DEVURL_DOMAIN"
  value: {{ .Values.devurls.host | quote }}
{{- end }}
{{- end }}
{{/*
  coder.namespaceWhitelist provides an environment variable
  that lists available namespaces for environment creation.
*/}}
{{- define "coder.namespaceWhitelist.env" }}
{{- if .Values.namespaceWhitelist }}
- name: NAMESPACE_WHITELIST
  value: {{ join "," .Values.namespaceWhitelist | quote }}
{{- end }}
{{- end }}
{{/* 
  coder.volumes adds a volumes stanza if a cert.secret is provided.
*/}}
{{- define "coder.volumes" }}
{{- if .Values.certs.secret.name }}
volumes:
{{- end }}
{{- if .Values.certs.secret.name }}
  - name: {{ .Values.certs.secret.name | quote }}
    secret:
      secretName: {{ .Values.certs.secret.name | quote }}
{{- end }}
{{- end }}

{{/* 
  coder.volumeMounts adds a volume mounts stanza if a cert.secret is provided.
*/}}
{{- define "coder.volumeMounts" }}
{{- if .Values.certs.secret.name }}
volumeMounts:
{{- end }}
{{- if .Values.certs.secret.name }}
  - name: {{ .Values.certs.secret.name | quote }}
    mountPath: /etc/ssl/certs/{{ .Values.certs.secret.key }}
    subPath: {{ .Values.certs.secret.key | quote }}
{{- end }}
{{- end }}
{{/*
  coder.serviceTolerations adds tolerations if any are specified to 
  coder-managed services.
*/}}
{{- define "coder.serviceTolerations" }}
{{- if .Values.serviceTolerations }}
tolerations:
{{ toYaml .Values.serviceTolerations }}
{{- end }}
{{- end }}
{{/*
  coder.envproxy.accessURL is a URL for accessing the envproxy.
*/}}
{{- define "coder.envproxy.accessURL" }}
{{- if ne .Values.envproxy.accessURL "" }}
{{- .Values.envproxy.accessURL -}}
{{- else if ne .Values.ingress.host "" }}
    {{- if .Values.ingress.tls.enable -}}
    https://
    {{- else -}}
    http://
    {{- end -}}
    {{- .Values.ingress.host }}/ingress
{{- else }}
{{- fail "envproxy.access.URL or ingress.host must be set" }}
{{- end }}
{{- end }}
