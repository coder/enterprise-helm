<!-- DO NOT EDIT. THIS IS GENERATED FROM README.md.gotmpl -->

# Coder Helm Chart

[![build](https://github.com/cdr/enterprise-helm/actions/workflows/build.yml/badge.svg?event=push)](https://github.com/cdr/enterprise-helm/actions/workflows/build.yml)
[![Twitter Follow](https://img.shields.io/twitter/follow/CoderHQ?label=%40CoderHQ&style=social)](https://twitter.com/coderhq)

Coder moves developer workspaces to your cloud and centralizes their creation and management. Keep developers in flow with the power of the cloud and a superior developer experience.

The Coder chart is the best way to operate Coder on Kubernetes. It contains all the required components, and can scale to large deployments.

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
| coderd | object | Primary service responsible for all things Coder! | `{"builtinProviderServiceAccount":{"annotations":{},"labels":{}},"devurlsHost":"","image":"","oidc":{"enableRefresh":false,"redirectOptions":{}},"podSecurityContext":{"runAsNonRoot":true,"runAsUser":1000,"seccompProfile":{"type":"RuntimeDefault"}},"replicas":1,"resources":{"limits":{"cpu":"250m","memory":"512Mi"},"requests":{"cpu":"250m","memory":"512Mi"}},"satellite":{"accessURL":"","enable":false,"primaryURL":""},"securityContext":{"allowPrivilegeEscalation":false,"readOnlyRootFilesystem":true,"seccompProfile":{"type":"RuntimeDefault"}},"serviceAnnotations":{},"serviceNodePorts":{"http":null,"https":null},"serviceSpec":{"externalTrafficPolicy":"Local","loadBalancerIP":"","loadBalancerSourceRanges":[],"type":"LoadBalancer"},"superAdmin":{"passwordSecret":{"key":"password","name":""}},"tls":{"devurlsHostSecretName":"","hostSecretName":""},"trustProxyIP":false}` |
| coderd.builtinProviderServiceAccount | object | Customize the built-in Kubernetes provider service account. | `{"annotations":{},"labels":{}}` |
| coderd.builtinProviderServiceAccount.annotations | object | A KV mapping of annotations. See: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/ | `{}` |
| coderd.builtinProviderServiceAccount.labels | object | Add labels to the service account used for the built-in provider. | `{}` |
| coderd.devurlsHost | string | Wildcard hostname to allow matching against custom-created dev URLs. Leaving as an empty string results in DevURLs being disabled. | `""` |
| coderd.image | string | Injected by Coder during release. | `""` |
| coderd.podSecurityContext | object | Fields related to the pod's security context (as opposed to the container). Some fields are also present in the container security context, which will take precedence over these values. | `{"runAsNonRoot":true,"runAsUser":1000,"seccompProfile":{"type":"RuntimeDefault"}}` |
| coderd.podSecurityContext.runAsNonRoot | bool | Requires that containers in the pod run as a non-privileged user. | `true` |
| coderd.podSecurityContext.runAsUser | int | Sets the user id of the pod. This must not be set to root (uid 0). | `1000` |
| coderd.podSecurityContext.seccompProfile | object | Sets the seccomp profile for the pod. If set, the container security context setting will take precedence over this value. | `{"type":"RuntimeDefault"}` |
| coderd.replicas | int | The number of Kubernetes Pod replicas. | `1` |
| coderd.resources | object | Kubernetes resource specification for coderd pods. To unset a value, set it to "". To unset all values, set resources to nil. | `{"limits":{"cpu":"250m","memory":"512Mi"},"requests":{"cpu":"250m","memory":"512Mi"}}` |
| coderd.satellite | object | Deploy a satellite to geodistribute access to workspaces for lower latency. | `{"accessURL":"","enable":false,"primaryURL":""}` |
| coderd.satellite.accessURL | string | URL of the satellite that clients will connect to. e.g. https://sydney.coder.myorg.com | `""` |
| coderd.satellite.enable | bool | Run coderd as a satellite pointing to a primary deployment. Satellite enable low-latency access to workspaces all over the world. Read more: TODO: Link to docs. | `false` |
| coderd.satellite.primaryURL | string | URL of the primary Coder deployment. Must be accessible from the satellite and clients. eg. https://coder.myorg.com | `""` |
| coderd.securityContext | object | Fields related to the container's security context (as opposed to the pod). Some fields are also present in the pod security context, in which case these values will take precedence. | `{"allowPrivilegeEscalation":false,"readOnlyRootFilesystem":true,"seccompProfile":{"type":"RuntimeDefault"}}` |
| coderd.securityContext.allowPrivilegeEscalation | bool | Controls whether the container can gain additional privileges, such as escalating to root. It is recommended to leave this setting disabled in production. | `false` |
| coderd.securityContext.readOnlyRootFilesystem | bool | Mounts the container's root filesystem as read-only. It is recommended to leave this setting enabled in production. This will override the same setting in the pod | `true` |
| coderd.securityContext.seccompProfile | object | Sets the seccomp profile for the migration and runtime containers. | `{"type":"RuntimeDefault"}` |
| coderd.serviceAnnotations | object | Extra annotations to apply to the coderd service. | `{}` |
| coderd.serviceNodePorts | object | Allows manually setting static node ports for the coderd service. This is only helpful if static ports are required, and usually should be left alone. By default these are dynamically chosen. | `{"http":null,"https":null}` |
| coderd.serviceNodePorts.http | string | Sets a static 'coderd' service non-TLS nodePort. This should usually be omitted. | `nil` |
| coderd.serviceNodePorts.https | string | Sets a static 'coderd' service TLS nodePort This should usually be omitted. | `nil` |
| coderd.serviceSpec | object | Specification to inject for the coderd service. See: https://kubernetes.io/docs/concepts/services-networking/service/ | `{"externalTrafficPolicy":"Local","loadBalancerIP":"","loadBalancerSourceRanges":[],"type":"LoadBalancer"}` |
| coderd.serviceSpec.externalTrafficPolicy | string | Set the traffic policy for the service. See: https://kubernetes.io/docs/tasks/access-application-cluster/create-external-load-balancer/#preserving-the-client-source-ip | `"Local"` |
| coderd.serviceSpec.loadBalancerIP | string | Set the external IP address of the Ingress service. | `""` |
| coderd.serviceSpec.loadBalancerSourceRanges | list | Traffic through the LoadBalancer will be restricted to the specified client IPs. This field will be ignored if the cloud provider does not support this feature. | `[]` |
| coderd.serviceSpec.type | string | Set the type of Service. See: https://kubernetes.io/docs/concepts/services-networking/service/#publishing-services-service-types | `"LoadBalancer"` |
| coderd.superAdmin.passwordSecret.key | string | The key of the secret that contains the super admin password. | `"password"` |
| coderd.superAdmin.passwordSecret.name | string | Name of a secret that should be used to determine the password for the super admin account. The password should be contained in the field `password`, or the manually specified one. | `""` |
| coderd.tls | object | TLS configuration for coderd. These options will override dashboard configuration. | `{"devurlsHostSecretName":"","hostSecretName":""}` |
| coderd.tls.devurlsHostSecretName | string | The secret to use for DevURL TLS. | `""` |
| coderd.tls.hostSecretName | string | The secret to use for TLS. | `""` |
| coderd.trustProxyIP | bool | Whether Coder should trust X-Real-IP and/or X-Forwarded-For headers from your reverse proxy. This should only be turned on if you're using a reverse proxy that sets both of these headers. This is always enabled if the Nginx ingress is deployed. | `false` |
| envbox | object | Required for running Docker inside containers. See requirements: https://coder.com/docs/coder/v1.19/admin/workspace-management/cvms | `{"image":""}` |
| envbox.image | string | Injected by Coder during release. | `""` |
| logging | object | Configures the logging format and output of Coder. | `{"human":"/dev/stderr","json":"","splunk":{"channel":"","token":"","url":""},"stackdriver":""}` |
| logging.human | string | Location to send logs that are formatted for readability. Set to an empty string to disable. | `"/dev/stderr"` |
| logging.json | string | Location to send logs that are formatted as JSON. Set to an empty string to disable. | `""` |
| logging.splunk | object | Coder can send logs directly to Splunk in addition to file-based output. | `{"channel":"","token":"","url":""}` |
| logging.splunk.token | string | Splunk HEC collector token. | `""` |
| logging.splunk.url | string | Splunk HEC collector endpoint. | `""` |
| logging.stackdriver | string | Location to send logs that are formatted for Google Stackdriver. Set to an empty string to disable. | `""` |
| metrics | object | Configure various metrics to gain observability into Coder. | `{"amplitudeKey":""}` |
| metrics.amplitudeKey | string | Enables telemetry pushing to Amplitude. Keep empty to disable | `""` |
| postgres.database | string | Name of the database that Coder will use. You must create this database first. | `""` |
| postgres.default | object | Configure a built-in PostgreSQL deployment. | `{"enable":true,"image":"","resources":{"limits":{"cpu":"250m","memory":"1Gi"},"requests":{"cpu":"250m","memory":"1Gi","storage":"10Gi"}},"storageClassName":""}` |
| postgres.default.enable | bool | Deploys a PostgreSQL instance. We recommend using an external PostgreSQL instance in production. If true, all other values are ignored. | `true` |
| postgres.default.image | string | Injected by Coder during release. | `""` |
| postgres.default.resources | object | Kubernetes resource specification for the PostgreSQL pod. To unset a value, set it to "". To unset all values, set resources to nil. | `{"limits":{"cpu":"250m","memory":"1Gi"},"requests":{"cpu":"250m","memory":"1Gi","storage":"10Gi"}}` |
| postgres.default.resources.requests.storage | string | Specifies the size of the volume claim for persisting the database. | `"10Gi"` |
| postgres.default.storageClassName | string | Set the storageClass to store the database. | `""` |
| postgres.host | string | Host of the external PostgreSQL instance. | `""` |
| postgres.passwordSecret | string | Name of an existing secret in the current namespace with the password of the PostgreSQL instance. The password must be contained in the secret field `password`. This should be set to an empty string if the database does not require a password to connect. | `""` |
| postgres.port | string | Port of the external PostgreSQL instance. | `""` |
| postgres.sslMode | string | Provides variable levels of protection for the PostgreSQL connection. For acceptable values, see:  https://www.postgresql.org/docs/9.1/libpq-ssl.html | `"require"` |
| postgres.user | string | User of the external PostgreSQL instance. | `""` |
| services | object | Kubernetes Service configuration that applies to Coder services. | `{"annotations":{},"clusterDomainSuffix":".svc.cluster.local","nodeSelector":{"kubernetes.io/arch":"amd64","kubernetes.io/os":"linux"},"tolerations":[],"type":"ClusterIP"}` |
| services.annotations | object | A KV mapping of annotations. See: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/ | `{}` |
| services.clusterDomainSuffix | string | Custom domain suffix for DNS resolution in your cluster. See: https://kubernetes.io/docs/tasks/administer-cluster/dns-custom-nameservers/ | `".svc.cluster.local"` |
| services.nodeSelector | object | See: https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#nodeselector | `{"kubernetes.io/arch":"amd64","kubernetes.io/os":"linux"}` |
| services.tolerations | list | Each element is a toleration object. See: https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/ | `[]` |
| services.type | string | See the following for configurable types: https://kubernetes.io/docs/concepts/services-networking/service/#publishing-services-service-types | `"ClusterIP"` |

## Contributing

Templates for Coder services live in the `templates` directory.
Helm compiles templates with `values.yaml` when deploying.

`README.md` is generated from `README.md.gotmpl` to ensure values are correct. Regenerate the readme:

```shell-session
$ make README.md
```

Deprecation notices should be added to `templates/NOTES.txt`.

## Support

If you experience issues, have feedback, or want to ask a question, open an issue or
pull request in this repository. Feel free to [contact us instead](https://coder.com/contact).

## Copyright and License

Copyright (C) 2020-2021 Coder Technologies Inc.

This program is free software: you can redistribute it and/or modify it under
the terms of the GNU General Public License as published by the Free Software
Foundation, either version 3 of the License, or (at your option) any later
version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY
WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
PARTICULAR PURPOSE.  See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with
this program.  If not, see <https://www.gnu.org/licenses/>.
