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
  coder.cemanager.accessURL is a URL for accessing the cemanager.
*/}}
{{- define "coder.cemanager.accessURL" }}
{{- if ne .Values.cemanager.accessURL "" }}
{{- .Values.cemanager.accessURL -}}
{{- else -}}
{{- fail "cemanager.accessURL must be set" }}
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
    {{- .Values.ingress.host }}
{{- else }}
{{- fail "envproxy.accessURL or ingress.host must be set" }}
{{- end }}
{{- end }}
{{/*
  coder.cluster.accessURL is a URL for accessing the Kubernetes cluster.
*/}}
{{- define "coder.cluster.accessURL" }}
{{- if ne .Values.envproxy.clusterAddress "" }}
{{- .Values.envproxy.clusterAddress -}}
{{- else -}}
{{- fail "envproxy.clusterAddress must be set" }}
{{- end }}
{{- end }}