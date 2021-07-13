{{- /* Defines a KV mapping of values that have been moved. */}}
{{- /* If a value specified in this mapping is set, we notify */}}
{{- /* of deprecation in NOTES.txt */}}
{{- define "moved" }}
{{- $moved := dict }}
{{- /* To deprecate a value, map the new location to the old below */}}
{{- $_ := set $moved "coderd" "cemanager" }}
{{- $_ := set $moved "coderd.replicas" "cemanager.replicas" }}
{{- $_ := set $moved "coderd.image" "cemanager.image" }}
{{- $_ := set $moved "coderd.resources" "cemanager.resources" }}
{{- $_ := set $moved "coderd.devurlsHost" "devurls.host" }}
{{- $_ := set $moved "coderd.serviceSpec.loadBalancerIP" "ingress.loadBalancerIP" }}
{{- $_ := set $moved "coderd.serviceSpec.loadBalancerSourceRanges" "ingress.loadBalancerSourceRanges" }}
{{- $_ := set $moved "coderd.serviceSpec.externalTrafficPolicy" "ingress.service.externalTrafficPolicy" }}
{{- $_ := set $moved "coderd.tls.hostSecretName" "ingress.tls.hostSecretName" }}
{{- $_ := set $moved "coderd.tls.devurlsHostSecretName" "ingress.tls.devurlsHostSecretName" }}
{{- $_ := set $moved "postgres.default.storageClassName" "storageClassName" }}
{{- $_ := set $moved "postgres.default.image" "timescale.image" }}
{{- $_ := set $moved "postgres.default.resources" "timescale.resources" }}
{{- $_ := set $moved "postgres.default.resources.requests.storage" "timescale.resources.requests.storage" }}
{{- $_ := set $moved "postgres.default.enable" "postgres.useDefault" }}
{{- $_ := set $moved "services.annotations" "deploymentAnnotations" }}
{{- $_ := set $moved "services.tolerations" "serviceTolerations" }}
{{- $_ := set $moved "services.clusterDomainSuffix" "clusterDomainSuffix" }}
{{- $_ := set $moved "services.type" "serviceType" }}
{{- $_ := set $moved "coderd.builtinProviderServiceAccount.annotations" "serviceAccount.annotations" }}
{{- $_ := set $moved "coderd.builtinProviderServiceAccount.labels" "serviceAccount.labels" }}
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

  {{- $moved := fromJson (include "moved" .) }}
  {{- $oldkey := index $moved $key }}
  {{- $oldvalue := "" }}
  {{- if $oldkey }}
    {{- $oldvalue = include "movedValue" (dict "Values" .Values "Key" $oldkey "Default" .Default "Nested" true) }}
  {{- else if not .Nested }}
    {{ fail "Developer Error: 'movedValue' is used for deprecated values only. Reference the value directly instead!" }}
  {{- end }}

  {{- if ne $oldvalue "" }}
    {{- $oldvalue }}
  {{- else }}
    {{- /* Iterate through the provided key split by "." */}}
    {{- /* eg. "some.kinda.key" is ["some", "kinda", "key"] */}}
    {{- range $index, $keypart := splitList "." $key }}
      {{- /* There's no way to break this loop */}}
      {{- /* If not found once, we know the chain is broken */}}
      {{- if $found }}
        {{- $values = index $values $keypart }}
        {{- if not $values }}
          {{- $found = false }}
        {{- end }}
        {{- if kindIs "bool" $values }}
          {{- $found = true }}
        {{- end }}
      {{- end }}
    {{- end }}

    {{- if $found }}
      {{- toYaml $values }}
    {{- else }}
      {{- if and (not .Nested) .Default }}
        {{- toYaml .Default }}
      {{- end }}
    {{- end }}
  {{- end }}
{{- end }}
