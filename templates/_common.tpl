{{/* 
  coder.storageClassName adds a storageClassName field to a volume claim
  if the 'storageClassName' value is non-empty.
*/}}
{{- define "coder.storageClassName" }}
{{ $storageClass := include "movedValue" (dict "Values" .Values "Key" "postgres.default.storageClassName") }}
{{- if ne $storageClass "" }}
storageClassName: {{ $storageClass | default "" | quote }}
{{- end }}
{{- end }}
{{/*
  coder.postgres.env adds environment variables that
  specify how to connect to a Postgres instance.
*/}}
{{- define "coder.postgres.env" }}
{{- if eq (include "movedValue" (dict "Values" .Values "Key" "postgres.default.enable" "Default" true)) "true" }}
- name: DB_HOST
  value: timescale.{{ .Release.Namespace }}{{ include "movedValue" (dict "Values" .Values "Key" "services.clusterDomainSuffix") }}
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
- name: DB_PASSWORD
  valueFrom:
    secretKeyRef:
      name: {{ .Values.postgres.passwordSecret | quote }}    
      key: password
- name: DB_SSL_MODE
  value: {{ .Values.postgres.sslMode | quote }}
- name: DB_NAME
  value: {{ .Values.postgres.database | quote }}
- name: DB_IDLE_CONNECTIONS
  value: {{ .Values.postgres.idleConnections | quote }}
{{- end }}
{{- end }}
{{/*
  coder.volumes adds a volumes stanza if a cert.secret is provided.
*/}}
{{- define "coder.volumes" }}
{{- if or (merge .Values dict | dig "certs" "secret" "name" false) (ne (include "movedValue" (dict "Values" .Values "Key" "coderd.tls.hostSecretName")) "") }}
volumes:
{{- end }}
{{- if (merge .Values dict | dig "certs" "secret" "name" false) }}
  - name: {{ .Values.certs.secret.name | quote }}
    secret:
      secretName: {{ .Values.certs.secret.name | quote }}
{{- end }}
{{- if ne (include "movedValue" (dict "Values" .Values "Key" "coderd.tls.hostSecretName")) "" }}
  - name: tls
    secret:
      secretName: {{ include "movedValue" (dict "Values" .Values "Key" "coderd.tls.hostSecretName") }}
{{- end }}
{{- end }}

{{/* 
  coder.volumeMounts adds a volume mounts stanza if a cert.secret is provided.
*/}}
{{- define "coder.volumeMounts" }}
{{- if or (merge .Values dict | dig "certs" "secret" "name" false) (ne (include "movedValue" (dict "Values" .Values "Key" "coderd.tls.hostSecretName")) "") }}
volumeMounts:
{{- end }}
{{- if (merge .Values dict | dig "certs" "secret" "name" false) }}
  - name: {{ .Values.certs.secret.name | quote }}
    mountPath: /etc/ssl/certs/{{ .Values.certs.secret.key }}
    subPath: {{ .Values.certs.secret.key | quote }}
{{- end }}
{{- if ne (include "movedValue" (dict "Values" .Values "Key" "coderd.tls.hostSecretName")) "" }}
  - name: tls
    mountPath: /etc/ssl/certs/host
    readOnly: true
{{- end }}
{{- end }}
{{/*
  coder.serviceTolerations adds tolerations if any are specified to 
  coder-managed services.
*/}}
{{- define "coder.serviceTolerations" }}
{{- if ne (include "movedValue" (dict "Values" .Values "Key" "services.tolerations")) "" }}
tolerations:
{{ include "movedValue" (dict "Values" .Values "Key" "services.tolerations") }}
{{- end }}
{{- end }}
{{/*
  coder.accessURL is a URL for accessing the coderd.
*/}}
{{- define "coder.accessURL" }}
{{- if .Values.cemanager }}
{{- if ne (merge .Values dict | dig "cemanager" "accessURL" "") "" }}
{{- .Values.cemanager.accessURL -}}
{{- else -}}
    http://cemanager.{{ .Release.Namespace }}{{ include "movedValue" (dict "Values" .Values "Key" "services.clusterDomainSuffix") }}:8080
{{- end }}
{{- else -}}
{{- if ne (merge .Values dict | dig "coderd" "accessURL" "") "" }}
{{- .Values.coderd.accessURL -}}
{{- else -}}
    http://coderd.{{ .Release.Namespace }}{{ include "movedValue" (dict "Values" .Values "Key" "services.clusterDomainSuffix") }}:8080
{{- end }}
{{- end }}
{{- end }}
{{/*
  coder.envproxy.accessURL is a URL for accessing the envproxy.
*/}}
{{- define "coder.envproxy.accessURL" }}
{{- if ne (merge .Values dict | dig "envproxy" "accessURL" "") "" }}
{{- .Values.envproxy.accessURL -}}
{{- else if ne (merge .Values dict | dig "ingress" "host" "") "" }}
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
{{- if ne (merge .Values dict | dig "envproxy" "clusterAddress" "") "" }}
{{- .Values.envproxy.clusterAddress -}}
{{- else -}}
    https://kubernetes.default{{ include "movedValue" (dict "Values" .Values "Key" "services.clusterDomainSuffix") }}:443
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

{{- define "coder.serviceName" }}
{{- if hasKey .Values "cemanager" -}}
cemanager
{{- else -}}
coderd
{{- end }}
{{- end }}
