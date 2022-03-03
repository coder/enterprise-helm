{{- if .Values.cemanager }}
{{- fail "The 'cemanager' setting was deprecated in 1.21 and removed in 1.27; use 'coderd' instead" }}
{{- end }}

{{- if .Values.devurls }}
{{- fail "The 'devurls.host' setting was deprecated in 1.21 and removed in 1.27; use 'coderd.devurlsHost' instead" }}
{{- end }}

{{- if .Values.ingress.loadBalancerIP }}
{{- fail "The 'ingress.loadbalancerIP' setting was deprecated in 1.21 and removed in 1.27; use 'coderd.serviceSpec.loadBalancerIP' instead" }}
{{- end }}

{{- if .Values.ingress.loadBalancerSourceRanges }}
{{- fail "The 'ingress.loadbalancerIP' setting was deprecated in 1.21 and removed in 1.27; use 'coderd.serviceSpec.loadBalancerSourceRanges' instead" }}
{{- end }}

{{- if hasKey .Values "ingress.service.externalTrafficPolicy" }}
{{- fail "The 'ingress.loadbalancerIP' setting was deprecated in 1.21 and removed in 1.27; use 'coderd.serviceSpec.externalTrafficPolicy' instead" }}
{{- end }}

{{- if .Values.storageClassName }}
{{- fail "The 'storageClassName' setting was deprecated in 1.21 and removed in 1.27; use 'postgres.default.storageClassName' instead" }}
{{- end }}

{{- if .Values.timescale }}
{{- fail "The 'timescale' setting was deprecated in 1.21 and removed in 1.27; use 'postgres.default' instead" }}
{{- end }}

{{- if .Values.postgres.useDefault }}
{{- fail "The 'postgres.useDefault' setting was deprecated in 1.21 and removed in 1.27; use 'postgres.default.enable' instead" }}
{{- end }}

{{- if .Values.deploymentAnnotations }}
{{- fail "The 'deploymentAnnotations' setting was deprecated in 1.21 and removed in 1.27; use 'services.annotations' instead" }}
{{- end }}

{{- if .Values.serviceTolerations }}
{{- fail "The 'serviceTolerations' setting was deprecated in 1.21 and removed in 1.27; use 'services.tolerations' instead" }}
{{- end }}

{{- if .Values.clusterDomainSuffix }}
{{- fail "The 'clusterDomainSuffix' setting was deprecated in 1.21 and removed in 1.27; use 'services.clusterDomainSuffix' instead" }}
{{- end }}

{{- if .Values.serviceType }}
{{- fail "The 'serviceType' setting was deprecated in 1.21 and removed in 1.27; use 'services.type' instead" }}
{{- end }}

{{- if .Values.serviceAccount }}
{{- fail "The 'serviceAccount' setting was deprecated in 1.21 and removed in 1.27; use 'coderd.builtinProviderServiceAccount' instead" }}
{{- end }}
