{{- /* Defines a KV mapping of values that have been moved. */}}
{{- /* If a value specified in this mapping is set, we notify */}}
{{- /* of deprecation in NOTES.txt */}}
{{- define "moved" }}
{{- $moved := dict }}
{{- $_ := set $moved "postgres.default.storageClassName" "storageClassName" }}
{{- toJson $moved }}
{{- end }}

{{- /*
  Use when a key has been moved for deprecation.
  Prioritizes the value of "New" above "Old".
  Provide a "Default" key to set a default value.
*/}}
{{- define "lookup" }}
  {{- $key := required "`Key` must be set!" .Key }}

  {{- $values := required "`Values` must be set!" .Values }}
  {{- $found := true }}

  {{- /* Iterate through the provided key split by "." */}}
  {{- /* eg. "some.kinda.key" is ["some", "kinda", "key"] */}}
  {{- range $index, $keypart := splitList "." $key }}
    {{- /* There's no way to break this loop */}}
    {{- /* If not found once, we know the chain is broken */}}
    {{- if $found }}
      {{- if index $values $keypart  }}
        {{- $values = index $values $keypart }}
      {{- else }}
        {{- $found = false }}
      {{- end }}
    {{- end }}
  {{- end }}

  {{- if $found }}
    {{- $values }}
  {{- else }}
    {{- $moved := fromJson (include "moved" .) }}
    {{- $key = index $moved $key }}
    {{- if $key }}
      {{- /* We can use this function to check for the key again! */}}
      {{- include "lookup" (dict "Values" .Values "Key" $key "Default" .Default) }}
    {{- else }}
      {{- .Default | default "" }}
    {{- end }}
  {{- end }}
{{- end }}
