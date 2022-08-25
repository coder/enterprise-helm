<!-- DO NOT EDIT. THIS IS GENERATED FROM README.md.gotmpl -->

# Coder Helm Chart

[![build](https://github.com/cdr/enterprise-helm/actions/workflows/build.yml/badge.svg?event=push)](https://github.com/cdr/enterprise-helm/actions/workflows/build.yml)
[![Twitter Follow](https://img.shields.io/twitter/follow/CoderHQ?label=%40CoderHQ&style=social)](https://twitter.com/coderhq)

Coder moves developer workspaces to your cloud and centralizes their creation and management. Keep developers in flow with the power of the cloud and a superior developer experience.

The Coder Helm Chart is the best way to install and operate Coder on Kubernetes. It contains all the required components, and can scale to large deployments.

![Coder Dashboard](./assets/coder.svg)

## Getting Started

> ⚠️ **Warning**: This repository will not represent the latest Coder release. Reference
our installation docs for instructions on a tagged release.

View [our docs](https://coder.com/docs/setup/installation) for detailed installation instructions.

## Values

| Key                 | Type | Description | Default                         |
| ------------------- | ---- | ----------- | ------------------------------- |
| certs | object | Certificate that will be mounted inside Coder services. | `{"secret":{"key":"","name":""}}` |
| certs.secret.key | string | Key pointing to a certificate in the secret. | `""` |
| certs.secret.name | string | Name of the secret. | `""` |
| coderd | object | Primary service responsible for all things Coder! | `{"affinity":{"podAntiAffinity":{"preferredDuringSchedulingIgnoredDuringExecution":[{"podAffinityTerm":{"labelSelector":{"matchExpressions":[{"key":"app.kubernetes.io/name","operator":"In","values":["coderd"]}]},"topologyKey":"kubernetes.io/hostname"},"weight":1}]}},"alternateHostnames":[],"annotations":{},"builtinProviderServiceAccount":{"annotations":{},"labels":{},"migrate":true},"clientTLS":{"secretName":""},"devurlsHost":"","extraEnvs":[],"extraLabels":{},"image":"","imagePullSecret":"","liveness":{"failureThreshold":30,"initialDelaySeconds":30,"periodSeconds":10,"timeoutSeconds":3},"networkPolicy":{"enable":true},"oidc":{"enableRefresh":false,"redirectOptions":{}},"podSecurityContext":{"runAsGroup":1000,"runAsNonRoot":true,"runAsUser":1000,"seccompProfile":{"type":"RuntimeDefault"}},"proxy":{"exempt":"cluster.local","http":"","https":""},"readiness":{"failureThreshold":15,"initialDelaySeconds":10,"periodSeconds":10,"timeoutSeconds":3},"replicas":1,"resources":{"limits":{"cpu":"250m","memory":"512Mi"},"requests":{"cpu":"250m","memory":"512Mi"}},"reverseProxy":{"headers":[],"trustedOrigins":[]},"satellite":{"accessURL":"","enable":false,"primaryURL":""},"scim":{"authSecret":{"key":"secret","name":""},"enable":false},"securityContext":{"allowPrivilegeEscalation":false,"readOnlyRootFilesystem":true,"runAsGroup":1000,"runAsNonRoot":true,"runAsUser":1000,"seccompProfile":{"type":"RuntimeDefault"}},"serviceAnnotations":{},"serviceNodePorts":{"http":null,"https":null},"serviceSpec":{"externalTrafficPolicy":"Local","loadBalancerIP":"","loadBalancerSourceRanges":[],"type":"LoadBalancer"},"superAdmin":{"passwordSecret":{"key":"password","name":""}},"tls":{"devurlsHostSecretName":"","hostSecretName":""},"trustProxyIP":false,"workspaceServiceAccount":{"annotations":{},"labels":{}}}` |
| coderd.affinity | object | Allows specifying an affinity rule for the `coderd` deployment. The default rule prefers to schedule coderd pods on different nodes, which is only applicable if coderd.replicas is greater than 1. | `{"podAntiAffinity":{"preferredDuringSchedulingIgnoredDuringExecution":[{"podAffinityTerm":{"labelSelector":{"matchExpressions":[{"key":"app.kubernetes.io/name","operator":"In","values":["coderd"]}]},"topologyKey":"kubernetes.io/hostname"},"weight":1}]}}` |
| coderd.alternateHostnames | list | A list of hostnames that coderd (including satellites) will allow for OIDC. If this list is not set, all OIDC traffic will go to the configured access URL in the admin settings on the dashboard (or the satellite's primary URL as configured by Helm). | `[]` |
| coderd.annotations | object | Apply annotations to the coderd deployment. https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/ | `{}` |
| coderd.builtinProviderServiceAccount | object | Customize the built-in Kubernetes provider service account. | `{"annotations":{},"labels":{},"migrate":true}` |
| coderd.builtinProviderServiceAccount.annotations | object | A KV mapping of annotations. See: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/ | `{}` |
| coderd.builtinProviderServiceAccount.labels | object | Add labels to the service account used for the built-in provider. | `{}` |
| coderd.builtinProviderServiceAccount.migrate | bool | Will migrate the built-in workspace provider using the coded environment. | `true` |
| coderd.clientTLS | object | Client-side TLS configuration for coderd. | `{"secretName":""}` |
| coderd.clientTLS.secretName | string | Secret containing a PEM encoded cert file. | `""` |
| coderd.devurlsHost | string | Wildcard hostname to allow matching against custom-created dev URLs. Leaving as an empty string results in DevURLs being disabled. | `""` |
| coderd.extraEnvs | list | Add additional environment variables to the coderd deployment containers. Overriding any environment variables that the Helm chart sets automatically is unsupported and will result in undefined behavior. You can find a list of the environment variables we set by default by inspecting the helm template files or by running `kubectl describe` against your existing coderd deployment. https://kubernetes.io/docs/tasks/inject-data-application/define-environment-variable-container/ | `[]` |
| coderd.extraLabels | object | Allows specifying additional labels to pods in the `coderd` deployment (.spec.template.metadata.labels). | `{}` |
| coderd.image | string | Injected by Coder during release. | `""` |
| coderd.imagePullSecret | string | The secret used for pulling the coderd image from a private registry. | `""` |
| coderd.liveness | object | Configure the liveness check for the coderd service. | `{"failureThreshold":30,"initialDelaySeconds":30,"periodSeconds":10,"timeoutSeconds":3}` |
| coderd.networkPolicy | object | Configure the network policy to apply to coderd. | `{"enable":true}` |
| coderd.networkPolicy.enable | bool | Manage a network policy for coderd using Helm. If false, no policies will be created for the Coder control plane. | `true` |
| coderd.podSecurityContext | object | Fields related to the pod's security context (as opposed to the container). Some fields are also present in the container security context, which will take precedence over these values. | `{"runAsGroup":1000,"runAsNonRoot":true,"runAsUser":1000,"seccompProfile":{"type":"RuntimeDefault"}}` |
| coderd.podSecurityContext.runAsGroup | int | Sets the group id of the pod. For security reasons, we recommend using a non-root group. | `1000` |
| coderd.podSecurityContext.runAsNonRoot | bool | Requires that containers in the pod run as an unprivileged user. If setting runAsUser to 0 (root), this will need to be set to false. | `true` |
| coderd.podSecurityContext.runAsUser | int | Sets the user id of the pod. For security reasons, we recommend using a non-root user. | `1000` |
| coderd.podSecurityContext.seccompProfile | object | Sets the seccomp profile for the pod. If set, the container security context setting will take precedence over this value. | `{"type":"RuntimeDefault"}` |
| coderd.proxy | object | Whether Coder should initiate outbound connections using a proxy. | `{"exempt":"cluster.local","http":"","https":""}` |
| coderd.proxy.exempt | string | Bypass the configured proxy rules for this comma-delimited list of hosts or prefixes. This corresponds to the no_proxy environment variable. | `"cluster.local"` |
| coderd.proxy.http | string | Proxy to use for HTTP connections. If unset, coderd will initiate HTTP connections directly. This corresponds to the http_proxy environment variable. | `""` |
| coderd.proxy.https | string | Proxy to use for HTTPS connections. If this is not set, coderd will use the HTTP proxy (if set), otherwise it will initiate HTTPS connections directly. This corresponds to the https_proxy environment variable. | `""` |
| coderd.readiness | object | Configure the readiness check for the coderd service. | `{"failureThreshold":15,"initialDelaySeconds":10,"periodSeconds":10,"timeoutSeconds":3}` |
| coderd.replicas | int | The number of Kubernetes Pod replicas. | `1` |
| coderd.resources | object | Kubernetes resource specification for coderd pods. To unset a value, set it to "". To unset all values, set resources to nil. | `{"limits":{"cpu":"250m","memory":"512Mi"},"requests":{"cpu":"250m","memory":"512Mi"}}` |
| coderd.reverseProxy | object | Whether Coder should trust proxy headers for inbound connections, important for ensuring correct IP addresses when an Ingress Controller, service mesh, or other Layer 7 reverse proxy are deployed in front of Coder. | `{"headers":[],"trustedOrigins":[]}` |
| coderd.reverseProxy.headers | list | A list of trusted headers. | `[]` |
| coderd.reverseProxy.trustedOrigins | list | A list of IPv4 or IPv6 subnets to consider trusted, specified in CIDR format. If hosts are part of a matching network, the configured headers will be trusted; otherwise, coderd will rely on the connecting client IP address. | `[]` |
| coderd.satellite | object | Deploy a satellite to geodistribute access to workspaces for lower latency. | `{"accessURL":"","enable":false,"primaryURL":""}` |
| coderd.satellite.accessURL | string | URL of the satellite that clients will connect to. e.g. https://sydney.coder.myorg.com | `""` |
| coderd.satellite.enable | bool | Run coderd as a satellite pointing to a primary deployment. Satellite enable low-latency access to workspaces all over the world. Read more: https://coder.com/docs/coder/latest/admin/satellites | `false` |
| coderd.satellite.primaryURL | string | URL of the primary Coder deployment. Must be accessible from the satellite and clients. eg. https://coder.myorg.com | `""` |
| coderd.scim.authSecret.key | string | The key of the secret that contains the SCIM auth header. | `"secret"` |
| coderd.scim.authSecret.name | string | Name of a secret that should be used to determine the auth header used for the SCIM server. The secret should be contained in the field `secret`, or the manually specified one. | `""` |
| coderd.scim.enable | bool | Enable SCIM support in coderd. SCIM allows you to automatically provision/deprovision users. If true, authSecret.name must be set. | `false` |
| coderd.securityContext | object | Fields related to the container's security context (as opposed to the pod). Some fields are also present in the pod security context, in which case these values will take precedence. | `{"allowPrivilegeEscalation":false,"readOnlyRootFilesystem":true,"runAsGroup":1000,"runAsNonRoot":true,"runAsUser":1000,"seccompProfile":{"type":"RuntimeDefault"}}` |
| coderd.securityContext.allowPrivilegeEscalation | bool | Controls whether the container can gain additional privileges, such as escalating to root. It is recommended to leave this setting disabled in production. | `false` |
| coderd.securityContext.readOnlyRootFilesystem | bool | Mounts the container's root filesystem as read-only. It is recommended to leave this setting enabled in production. This will override the same setting in the pod | `true` |
| coderd.securityContext.runAsGroup | int | Sets the group id of the pod. For security reasons, we recommend using a non-root group. | `1000` |
| coderd.securityContext.runAsNonRoot | bool | Requires that the coderd and migrations containers run as an unprivileged user. If setting runAsUser to 0 (root), this will need to be set to false. | `true` |
| coderd.securityContext.runAsUser | int | Sets the user id of the pod. For security reasons, we recommend using a non-root user. | `1000` |
| coderd.securityContext.seccompProfile | object | Sets the seccomp profile for the migration and runtime containers. | `{"type":"RuntimeDefault"}` |
| coderd.serviceAnnotations | object | Apply annotations to the coderd service. https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/ | `{}` |
| coderd.serviceNodePorts | object | Allows manually setting static node ports for the coderd service. This is only helpful if static ports are required, and usually should be left alone. By default these are dynamically chosen. | `{"http":null,"https":null}` |
| coderd.serviceNodePorts.http | string | Sets a static 'coderd' service non-TLS nodePort. This should usually be omitted. | `nil` |
| coderd.serviceNodePorts.https | string | Sets a static 'coderd' service TLS nodePort This should usually be omitted. | `nil` |
| coderd.serviceSpec | object | Specification to inject for the coderd service. See: https://kubernetes.io/docs/concepts/services-networking/service/ | `{"externalTrafficPolicy":"Local","loadBalancerIP":"","loadBalancerSourceRanges":[],"type":"LoadBalancer"}` |
| coderd.serviceSpec.externalTrafficPolicy | string | Set the traffic policy for the service. See: https://kubernetes.io/docs/tasks/access-application-cluster/create-external-load-balancer/#preserving-the-client-source-ip | `"Local"` |
| coderd.serviceSpec.loadBalancerIP | string | Set the IP address of the coderd service. | `""` |
| coderd.serviceSpec.loadBalancerSourceRanges | list | Traffic through the LoadBalancer will be restricted to the specified client IPs. This field will be ignored if the cloud provider does not support this feature. | `[]` |
| coderd.serviceSpec.type | string | Set the type of Service. See: https://kubernetes.io/docs/concepts/services-networking/service/#publishing-services-service-types | `"LoadBalancer"` |
| coderd.superAdmin.passwordSecret.key | string | The key of the secret that contains the super admin password. | `"password"` |
| coderd.superAdmin.passwordSecret.name | string | Name of a secret that should be used to determine the password for the super admin account. The password should be contained in the field `password`, or the manually specified one. | `""` |
| coderd.tls | object | TLS configuration for coderd. These options will override dashboard configuration. | `{"devurlsHostSecretName":"","hostSecretName":""}` |
| coderd.tls.devurlsHostSecretName | string | The secret to use for DevURL TLS. | `""` |
| coderd.tls.hostSecretName | string | The secret to use for TLS. | `""` |
| coderd.trustProxyIP | bool | Configures Coder to accept X-Real-IP and X-Forwarded-For headers from any origin. This option is deprecated and will be removed in a future release. Use the coderd.reverseProxy setting instead, which supports configuring an allowlist of trusted origins. | `false` |
| coderd.workspaceServiceAccount | object | Customize the default service account used for workspaces. | `{"annotations":{},"labels":{}}` |
| coderd.workspaceServiceAccount.annotations | object | A KV mapping of annotations. See: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/ | `{}` |
| coderd.workspaceServiceAccount.labels | object | Add labels to the service account used for workspaces. | `{}` |
| envbox | object | Required for running Docker inside containers. See requirements: https://coder.com/docs/coder/latest/admin/workspace-management/cvms | `{"image":""}` |
| envbox.image | string | Injected by Coder during release. | `""` |
| ingress | object | Configure an Ingress to route traffic to Coder services. | `{"annotations":{"nginx.ingress.kubernetes.io/proxy-body-size":"0"},"className":"","enable":false,"host":"","tls":{"enable":false}}` |
| ingress.annotations | object | Additional annotations to add to the Ingress object. The behavior is typically dependent on the Ingress Controller implementation, and useful for managing features like TLS termination. | `{"nginx.ingress.kubernetes.io/proxy-body-size":"0"}` |
| ingress.className | string | The ingressClassName to set on the Ingress. | `""` |
| ingress.enable | bool | A boolean controlling whether to create an Ingress. | `false` |
| ingress.host | string | The hostname to proxy to the Coder installation. The cluster Ingress Controller typically uses server name indication or the HTTP Host header to route traffic. The dev URLs hostname is specified in coderd.devurlsHost. | `""` |
| ingress.tls | object | Configures TLS settings for the Ingress. TLS certificates are specified in coderd.tls.hostSecretName and coderd.tls.devurlsHostSecretName. | `{"enable":false}` |
| ingress.tls.enable | bool | Determines whether the Ingress handles TLS. | `false` |
| logging | object | Configures the logging format and output of Coder. | `{"human":"/dev/stderr","json":"","splunk":{"channel":"","token":"","url":""},"stackdriver":"","verbose":true}` |
| logging.human | string | Location to send logs that are formatted for readability. Set to an empty string to disable. | `"/dev/stderr"` |
| logging.json | string | Location to send logs that are formatted as JSON. Set to an empty string to disable. | `""` |
| logging.splunk | object | Coder can send logs directly to Splunk in addition to file-based output. | `{"channel":"","token":"","url":""}` |
| logging.splunk.token | string | Splunk HEC collector token. | `""` |
| logging.splunk.url | string | Splunk HEC collector endpoint. | `""` |
| logging.stackdriver | string | Location to send logs that are formatted for Google Stackdriver. Set to an empty string to disable. | `""` |
| logging.verbose | bool | Toggles coderd debug logging. | `true` |
| metrics | object | Configure various metrics to gain observability into Coder. | `{"amplitudeKey":""}` |
| metrics.amplitudeKey | string | Enables telemetry pushing to Amplitude. Amplitude records how users interact with Coder, which is used to improve the product. No events store any personal information. Amplitude can be found here: https://amplitude.com/ Keep empty to disable. | `""` |
| postgres.connector | string | Option for configuring database connector type. valid values are: - "postgres" -- default connector - "awsiamrds" -- uses AWS IAM account in environment to authenticate using   IAM to connect to an RDS instance. | `"postgres"` |
| postgres.database | string | Name of the database that Coder will use. You must create this database first. | `""` |
| postgres.default | object | Configure a built-in PostgreSQL deployment. | `{"annotations":{},"enable":true,"image":"","networkPolicy":{"enable":true},"resources":{"limits":{"cpu":"250m","memory":"1Gi"},"requests":{"cpu":"250m","memory":"1Gi","storage":"10Gi"}},"storageClassName":""}` |
| postgres.default.annotations | object | Apply annotations to the default postgres service. https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/ | `{}` |
| postgres.default.enable | bool | Deploys a PostgreSQL instance. We recommend using an external PostgreSQL instance in production. If true, all other values are ignored. | `true` |
| postgres.default.image | string | Injected by Coder during release. | `""` |
| postgres.default.networkPolicy | object | Configure the network policy to apply to the built-in PostgreSQL deployment. | `{"enable":true}` |
| postgres.default.networkPolicy.enable | bool | Manage a network policy for PostgreSQL using Helm. If false, no policies will be created for the built-in database. | `true` |
| postgres.default.resources | object | Kubernetes resource specification for the PostgreSQL pod. To unset a value, set it to "". To unset all values, set resources to nil. | `{"limits":{"cpu":"250m","memory":"1Gi"},"requests":{"cpu":"250m","memory":"1Gi","storage":"10Gi"}}` |
| postgres.default.resources.requests.storage | string | Specifies the size of the volume claim for persisting the database. | `"10Gi"` |
| postgres.default.storageClassName | string | Set the storageClass to store the database. | `""` |
| postgres.host | string | Host of the external PostgreSQL instance. | `""` |
| postgres.passwordSecret | string | Name of an existing secret in the current namespace with the password of the PostgreSQL instance. The password must be contained in the secret field `password`. This should be set to an empty string if the database does not require a password to connect. | `""` |
| postgres.port | string | Port of the external PostgreSQL instance. | `""` |
| postgres.ssl | object | Options for configuring the SSL cert, key, and root cert when connecting to Postgres. | `{"certSecret":{"key":"","name":""},"keySecret":{"key":"","name":""},"rootCertSecret":{"key":"","name":""}}` |
| postgres.ssl.certSecret | object | Secret containing a PEM encoded cert file. | `{"key":"","name":""}` |
| postgres.ssl.certSecret.key | string | Key pointing to a certificate in the secret. | `""` |
| postgres.ssl.certSecret.name | string | Name of the secret. | `""` |
| postgres.ssl.keySecret | object | Secret containing a PEM encoded key file. | `{"key":"","name":""}` |
| postgres.ssl.keySecret.key | string | Key pointing to a certificate in the secret. | `""` |
| postgres.ssl.rootCertSecret | object | Secret containing a PEM encoded root cert file. | `{"key":"","name":""}` |
| postgres.ssl.rootCertSecret.key | string | Key pointing to a certificate in the secret. | `""` |
| postgres.ssl.rootCertSecret.name | string | Name of the secret. | `""` |
| postgres.sslMode | string | Provides variable levels of protection for the PostgreSQL connection. For acceptable values, see:  https://www.postgresql.org/docs/11/libpq-ssl.html | `"require"` |
| postgres.user | string | User of the external PostgreSQL instance. | `""` |
| services | object | Kubernetes Service configuration that applies to Coder services. | `{"annotations":{},"clusterDomainSuffix":".svc.cluster.local","nodeSelector":{"kubernetes.io/arch":"amd64","kubernetes.io/os":"linux"},"tolerations":[],"type":"ClusterIP"}` |
| services.annotations | object | A KV mapping of annotations. See: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/ DEPRECATED -- Please use the annotations value for each object. | `{}` |
| services.clusterDomainSuffix | string | Custom domain suffix for DNS resolution in your cluster. See: https://kubernetes.io/docs/tasks/administer-cluster/dns-custom-nameservers/ | `".svc.cluster.local"` |
| services.nodeSelector | object | See: https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#nodeselector | `{"kubernetes.io/arch":"amd64","kubernetes.io/os":"linux"}` |
| services.tolerations | list | Each element is a toleration object. See: https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/ | `[]` |
| services.type | string | See the following for configurable types: https://kubernetes.io/docs/concepts/services-networking/service/#publishing-services-service-types | `"ClusterIP"` |

## Contributing

Thanks for considering a contribution to this Chart!
Please see [CONTRIBUTING.md](CONTRIBUTING.md) for our conventions and practices.

## Support

If you experience issues, have feedback, or want to ask a question, open an issue or
pull request in this repository. Feel free to [contact us instead](https://coder.com/contact).

## Copyright and License

Copyright (C) 2020-2022 Coder Technologies Inc.

This program is free software: you can redistribute it and/or modify it under
the terms of the GNU General Public License as published by the Free Software
Foundation, either version 3 of the License, or (at your option) any later
version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY
WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
PARTICULAR PURPOSE.  See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with
this program.  If not, see <https://www.gnu.org/licenses/>.
