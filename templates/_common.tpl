{{/* 
	coder.storageClassName adds a storageClassName field to a volume claim
	if the 'storageClassName' value is non-empty.
*/}}
{{- define "coder.storageClassName" }}
{{- if ne .Values.storageClassName "" }}
storageClassName: {{ .Values.storageClassName | quote }}
{{- end }}
{{- end }}
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
  coder.postgres.env adds environment variables that
  specify how to connect to a Postgres instance.
*/}}
{{- define "coder.postgres.env" }}
{{- if eq .Values.postgres.useDefault true }}
- name: DB_HOST
  value: timescale.{{ .Release.Namespace }}{{ .Values.clusterDomainSuffix}}
- name: DB_PORT
  value: "5432"
- name: DB_USER
  value: coder
- name: DB_NAME
  value: coder
- name: DB_SSL_MODE
  value: disable
{{- else }}
- name: DB_HOST
  value: {{ .Values.postgres.host | quote }}
- name: DB_PORT
  value: {{ .Values.postgres.port | quote }}
- name: DB_USER
  value: {{ .Values.postgres.user | quote }}
- name: DB_SECRET
  value: {{ .Values.postgres.passwordSecret | quote }}
- name: DB_SSL_MODE
  value: {{ .Values.postgres.sslMode | quote }}
- name: DB_NAME
  value: {{ .Values.postgres.database | quote }}
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
