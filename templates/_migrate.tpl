{{- /* Defines a KV mapping of values that have been moved. */}}
{{- /* If a value specified in this mapping is set, we notify */}}
{{- /* of deprecation in NOTES.txt */}}
{{- define "moved" }}
{{- $moved := dict }}
{{- /* To deprecate a value, map the new location to the old below */}}
{{- $_ := set $moved "coderd" "cemanager" }}
{{- $_ := set $moved "coderd.serviceSpec.loadBalancerIP" "ingress.loadBalancerIP" }}
{{- $_ := set $moved "coderd.serviceSpec.loadBalancerSourceRanges" "ingress.loadBalancerSourceRanges" }}
{{- $_ := set $moved "coderd.serviceSpec.externalTrafficPolicy" "ingress.service.externalTrafficPolicy" }}
{{- $_ := set $moved "postgres.default.storageClassName" "storageClassName" }}
{{- $_ := set $moved "postgres.default.image" "timescale.image" }}
{{- $_ := set $moved "postgres.default.resources" "timescale.resources" }}
{{- $_ := set $moved "postgres.default.resources.requests.storage" "timescale.resources.requests.storage" }}
{{- $_ := set $moved "postgres.default.enable" "postgres.useDefault" }}
{{- $_ := set $moved "services.annotations" "deploymentAnnotations" }}
{{- $_ := set $moved "services.clusterDomainSuffix" "clusterDomainSuffix" }}
{{- toJson $moved }}
{{- end }}

{{- /*
  Use when a key has been moved for deprecation.
  Prioritizes the value of "New" above "Old".
  Provide a "Default" key to set a default value.

  Example:
  {{ include "movedValue" (dict "Values" .Values "Key" "postgres.default.storageClassName") }}
*/}}
{{- define "movedValue" }}
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
    {{- toYaml $values }}
  {{- else }}
    {{- $moved := fromJson (include "moved" .) }}
    {{- $key = index $moved $key }}
    {{- if $key }}
      {{- /* We can use this function to check for the key again! */}}
      {{- include "movedValue" (dict "Values" .Values "Key" $key "Default" .Default "Nested" true) }}
    {{- else }}
        {{- if not .Nested }}
        {{ fail "Developer Error: 'movedValue' is used for deprecated values only. Reference the value directly instead!" }}
      {{- end }}
      {{- toYaml .Default | default "" }}
    {{- end }}
  {{- end }}
{{- end }}
