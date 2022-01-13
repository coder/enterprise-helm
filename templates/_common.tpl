# coder.storageClassName adds a storageClassName field to a volume claim
# if the 'storageClassName' value is non-empty.
{{- define "coder.storageClassName" }}
{{- if .Values.postgres.default.storageClassName }}
storageClassName: {{ .Values.postgres.default.storageClassName | quote }}
{{- end }}
{{- end }}

# coder.postgres.env adds environment variables that
# specify how to connect to a Postgres instance.
{{- define "coder.postgres.env" }}
{{- if .Values.postgres.default.enable }}
- name: DB_HOST
  value: timescale.{{ .Release.Namespace }}{{ .Values.services.clusterDomainSuffix }}
- name: DB_PORT
  value: "5432"
- name: DB_USER
  value: coder
- name: DB_NAME
  value: coder
- name: DB_SSL_MODE
  value: disable
- name: DB_CONNECTOR
  value: {{ .Values.postgres.connector | quote }}
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
{{- if ne .Values.postgres.ssl.certSecret.name "" }}
- name: DB_CERT
  value: "/etc/ssl/certs/pg/cert/{{ .Values.postgres.ssl.certSecret.key }}"
{{- end }}
{{- if ne .Values.postgres.ssl.keySecret.name "" }}
- name: DB_KEY
  value: "/etc/ssl/certs/pg/key/{{ .Values.postgres.ssl.keySecret.key }}"
{{- end }}
{{- if ne .Values.postgres.ssl.rootCertSecret.name "" }}
- name: DB_ROOT_CERT
  value: "/etc/ssl/certs/pg/rootcert/{{ .Values.postgres.ssl.rootCertSecret.key }}"
{{- end }}
{{- end }}
{{- end }}

# coder.volumes adds a volumes stanza if a cert.secret is provided.
{{- define "coder.volumes" }}
volumes:
  - name: tmp-pgcerts
    emptyDir: {}
{{- if .Values.certs.secret.name }}
  - name: {{ .Values.certs.secret.name | quote }}
    secret:
      secretName: {{ .Values.certs.secret.name | quote }}
{{- end }}
{{- if .Values.coderd.tls.hostSecretName }}
  - name: tls
    secret:
      secretName: {{ .Values.coderd.tls.hostSecretName | quote }}
{{- end }}
{{- if .Values.coderd.tls.devurlsHostSecretName }}
  - name: devurltls
    secret:
      secretName: {{ .Values.coderd.tls.devurlsHostSecretName | quote }}
{{- end }}
{{- if ne .Values.postgres.ssl.certSecret.name "" }}
  - name: pgcert
    secret:
      secretName: {{ .Values.postgres.ssl.certSecret.name | quote }}
{{- end }}
{{- if ne .Values.postgres.ssl.keySecret.name "" }}
  - name: pgkey
    secret:
      secretName: {{ .Values.postgres.ssl.keySecret.name | quote }}
{{- end }}
{{- if ne .Values.postgres.ssl.rootCertSecret.name "" }}
  - name: pgrootcert
    secret:
      secretName: {{ .Values.postgres.ssl.rootCertSecret.name | quote }}
{{- end }}
{{- end }}

# coder.volumeMounts adds a volume mounts stanza if a cert.secret is
# provided.
{{- define "coder.volumeMounts" }}
volumeMounts:
  - name: tmp-pgcerts
    mountPath: /tmp/pgcerts
{{- if .Values.certs.secret.name }}
  - name: {{ .Values.certs.secret.name | quote }}
    mountPath: /etc/ssl/certs/{{ .Values.certs.secret.key }}
    subPath: {{ .Values.certs.secret.key | quote }}
{{- end }}
{{- if .Values.coderd.tls.hostSecretName }}
  - name: tls
    mountPath: /etc/ssl/certs/host
    readOnly: true
{{- end }}
{{- if .Values.coderd.tls.devurlsHostSecretName }}
  - name: devurltls
    mountPath: /etc/ssl/certs/devurls
    readOnly: true
{{- end }}
{{- if ne .Values.postgres.ssl.certSecret.name "" }}
  - name: pgcert
    mountPath: /etc/ssl/certs/pg/cert
    readOnly: true
{{- end }}
{{- if ne .Values.postgres.ssl.keySecret.name "" }}
  - name: pgkey
    mountPath: /etc/ssl/certs/pg/key
    readOnly: true
{{- end }}
{{- if ne .Values.postgres.ssl.rootCertSecret.name "" }}
  - name: pgrootcert
    mountPath: /etc/ssl/certs/pg/rootcert
    readOnly: true
{{- end }}
{{- end }}

# coder.serviceTolerations adds tolerations if any are specified to
# coder-managed services.
{{- define "coder.serviceTolerations" }}
{{- if .Values.services.tolerations }}
tolerations: {{ toYaml .Values.services.tolerations | nindent 2 }}
{{- end }}
{{- end }}

# coder.accessURL is a URL for accessing the coderd.
{{- define "coder.accessURL" }}
{{- if .Values.coderd.accessURL }}
{{- .Values.coderd.accessURL -}}
{{- else -}}
http://coderd.{{ .Release.Namespace }}{{ .Values.services.clusterDomainSuffix }}:8080
{{- end }}
{{- end }}

# coder.cluster.accessURL is a URL for accessing the Kubernetes cluster.
{{- define "coder.cluster.accessURL" -}}
https://kubernetes.default{{ .Values.services.clusterDomainSuffix }}:443
{{- end }}

# coder.services.nodeSelector adds nodeSelectors if any are specified to
# coder-managed services.
{{- define "coder.services.nodeSelector" }}
{{- if .Values.services.nodeSelector }}
nodeSelector: {{ toYaml .Values.services.nodeSelector | nindent 2 }}
{{- end }}
{{- end }}

{{- define "coder.serviceName" -}}
coderd
{{- end }}

{{- define "timescale.serviceName" -}}
timescale
{{- end }}
