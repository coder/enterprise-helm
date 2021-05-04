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
{{- if ne .Values.postgres.passwordSecret "" }}
- name: DB_PASSWORD
  valueFrom:
    secretKeyRef:
      name: {{ .Values.postgres.passwordSecret | quote }}    
      key: password
{{- end }}
- name: DB_SSL_MODE
  value: {{ .Values.postgres.sslMode | quote }}
- name: DB_NAME
  value: {{ .Values.postgres.database | quote }}
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
    http://cemanager.{{ .Release.Namespace }}{{ .Values.clusterDomainSuffix }}:8080
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
{{- end }}
{{- end }}
{{/*
  coder.cluster.accessURL is a URL for accessing the Kubernetes cluster.
*/}}
{{- define "coder.cluster.accessURL" }}
{{- if ne .Values.envproxy.clusterAddress "" }}
{{- .Values.envproxy.clusterAddress -}}
{{- else -}}
    https://kubernetes.default{{ .Values.clusterDomainSuffix }}:443
{{- end }}
{{- end }}

{{/*
  coder.services.nodeSelector adds nodeSelectors if any are specified to
  coder-managed services.
*/}}
{{- define "coder.services.nodeSelector" }}
{{- if .Values.services.nodeSelector }}
nodeSelector:
{{ toYaml .Values.services.nodeSelector | indent 1 }}
{{- end }}
{{- end }}

{{/*
  coder.array includes an array if it has one or more values.
*/}}
{{- define "coder.array" }}
{{- if ne (len .) 0 }}
{{ toYaml .}}
{{- end }}
{{- end }}
