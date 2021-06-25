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
| coderd | object | Primary service responsible for all things Coder! | `{"image":"","replica":{"enable":false,"primaryURL":""},"replicas":1,"resources":{"limits":{"cpu":"250m","memory":"512Mi"},"requests":{"cpu":"250m","memory":"512Mi"}},"securityContext":{"readOnlyRootFilesystem":true},"serviceSpec":{"loadBalancerIP":"","loadBalancerSourceRanges":[],"type":"LoadBalancer"},"tls":{"devurlsHostSecretName":"","hostSecretName":""}}` |
| coderd.image | string | Injected by Coder during release. | `""` |
| coderd.replica | object | Deploy a replica to geodistribute access to workspaces for lower latency. | `{"enable":false,"primaryURL":""}` |
| coderd.replica.enable | bool | Run coderd as a replica pointing to a primary deployment. | `false` |
| coderd.replica.primaryURL | string | URL of the primary deployment. eg. us-east.coder.myorg.com | `""` |
| coderd.replicas | int | The number of Kubernetes Pod replicas. | `1` |
| coderd.resources | object | Kubernetes resource specification for coderd pods. To unset a value, set it to "". To unset all values, set resources to nil. | `{"limits":{"cpu":"250m","memory":"512Mi"},"requests":{"cpu":"250m","memory":"512Mi"}}` |
| coderd.securityContext | object | Fields related to the container's security context (as opposed to the pod). | `{"readOnlyRootFilesystem":true}` |
| coderd.serviceSpec | object | Specification to inject for the coderd service. See: https://kubernetes.io/docs/concepts/services-networking/service/ | `{"loadBalancerIP":"","loadBalancerSourceRanges":[],"type":"LoadBalancer"}` |
| coderd.tls | object | TLS configuration for coderd. These options will override dashboard configuration. | `{"devurlsHostSecretName":"","hostSecretName":""}` |
| coderd.tls.devurlsHostSecretName | string | The secret to use for DevURL TLS. | `""` |
| coderd.tls.hostSecretName | string | The secret to use for TLS. | `""` |
| envbox | object | Required for running Docker inside containers. See requirements: https://coder.com/docs/coder/v1.19/admin/workspace-management/cvms | `{"image":""}` |
| envbox.image | string | Injected by Coder during release. | `""` |
| logging | object | Configures the logging format and output of Coder. | `{"human":"/dev/stderr","json":"","splunk":{"channel":"","token":"","url":""},"stackdriver":""}` |
| logging.human | string | Location to send logs that are formatted for readability. Set to an empty string to disable. | `"/dev/stderr"` |
| logging.json | string | Location to send logs that are formatted as JSON. Set to an empty string to disable. | `""` |
| logging.splunk | object | Coder can send logs directly to Splunk in addition to file-based output. | `{"channel":"","token":"","url":""}` |
| logging.splunk.token | string | Splunk HEC collector token. | `""` |
| logging.splunk.url | string | Splunk HEC collector endpoint. | `""` |
| logging.stackdriver | string | Location to send logs that are formatted for Google Stackdriver. Set to an empty string to disable. | `""` |
| postgres.database | string | Name of the database that Coder will use. You must create this database first. | `""` |
| postgres.default | object | Configure a built-in PostgreSQL deployment. | `{"enable":true,"image":"","resources":{"limits":{"cpu":"250m","memory":"1Gi"},"requests":{"cpu":"250m","memory":"1Gi","storage":"10Gi"}}}` |
| postgres.default.enable | bool | Deploys a PostgreSQL instance. We recommend using an external PostgreSQL instance in production. If true, all other values are ignored. | `true` |
| postgres.default.image | string | Injected by Coder during release. | `""` |
| postgres.default.resources | object | Kubernetes resource specification for the PostgreSQL pod. To unset a value, set it to "". To unset all values, set resources to nil. | `{"limits":{"cpu":"250m","memory":"1Gi"},"requests":{"cpu":"250m","memory":"1Gi","storage":"10Gi"}}` |
| postgres.default.resources.requests.storage | string | Specifies the size of the volume claim for persisting the database. | `"10Gi"` |
| postgres.host | string | Host of the external PostgreSQL instance. | `""` |
| postgres.passwordSecret | string | Name of an existing secret in the current namespace with the password of the PostgreSQL instance. The password must be contained in the secret field `password`. This should be set to an empty string if the database does not require a password to connect. | `""` |
| postgres.port | string | Port of the external PostgreSQL instance. | `""` |
| postgres.sslMode | string | Provides variable levels of protection for the PostgreSQL connection. For acceptable values, see:  https://www.postgresql.org/docs/9.1/libpq-ssl.html | `"require"` |
| postgres.user | string | User of the external PostgreSQL instance. | `""` |
| services | object | Kubernetes Service configuration that applies to Coder services.  | `{"annotations":{},"nodeSelector":{"kubernetes.io/arch":"amd64","kubernetes.io/os":"linux"},"tolerations":[],"type":"ClusterIP"}` |
| services.annotations | object | A KV mapping of annotations. See: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/ | `{}` |
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
