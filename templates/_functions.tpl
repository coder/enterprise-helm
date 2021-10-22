{{/*
  coder.resources accepts a resource stanza as its scope and returns
  resource configuration if any of the values are set.
*/}}
{{- define "coder.resources" }}
{{- if . }}
resources:
  {{- if .requests }}
  requests:
    {{- if .requests.cpu }}
    cpu: {{ .requests.cpu | quote }}
    {{- end }} 
    {{- if .requests.memory }}
    memory: {{ .requests.memory | quote }}
    {{- end }}
  {{- end }}
  {{- if .limits }}
  limits:
    {{- if .limits.cpu }}
    cpu: {{ .limits.cpu | quote }}
    {{- end }}
    {{- if .limits.memory }}
    memory: {{ .limits.memory | quote }}
    {{- end }} 
  {{- end }}
{{- end }}
{{- end }}

{{/*
coder.template.render -- renders a value that contains template.
Usage:
{{ include "coder.template.render" ( dict "value" .Values.path.to.the.Value "context" $) }}
*/}}
{{- define "coder.template.render" -}}
    {{- if typeIs "string" .value }}
        {{- tpl .value .context }}
    {{- else }}
        {{- tpl (.value | toYaml) .context }}
    {{- end }}
{{- end -}}
