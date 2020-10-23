# Coder Enterprise Helm

WARNING: The master branch may contain updates for a yet unreleased version of
Coder. The current state of the repo may not represent the latest release.

The Helm package here is only a template for releases: version numbers and
image URIs are injected and bundled into releases uploaded to
[helm.coder.com][helm-repo].

You can pull the official charts with `helm repo add coder https://helm.coder.com`.

Feel free to open PRs for fixes and improvements you've found while deploying
Coder!

## Values

| Key                                    | Type   | Description                                                                                                                                                                                                                                                                                                                                  | Default                                                                                                                       |
| -------------------------------------- | ------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------- |
| cemanager.image                        | string | Injected during releases.                                                                                                                                                                                                                                                                                                                    | `""`                                                                                                                          |
| cemanager.replicas                     | int    | The number of replicas to run of the manager.                                                                                                                                                                                                                                                                                                | `1`                                                                                                                           |
| cemanager.resources                    | object | Kubernetes resource request and limits for cemanager pods. To unset a value, set it to "". To unset all values, you can provide a values.yaml file which sets resources to nil. See values.yaml for an example.                                                                                                                              | `{"limits":{"cpu":"250m","memory":"512Mi"},"requests":{"cpu":"250m","memory":"512Mi"}}`                                       |
| certs                                  | object | Describes CAs that should be added to Coder services. These certs are NOT added to environments.                                                                                                                                                                                                                                             | `{"secret":{"key":"","name":""}}`                                                                                             |
| certs.secret.key                       | string | The key in the secret pointing to the certificate bundle.                                                                                                                                                                                                                                                                                    | `""`                                                                                                                          |
| certs.secret.name                      | string | The name of the secret.                                                                                                                                                                                                                                                                                                                      | `""`                                                                                                                          |
| clusterDomainSuffix                    | string | If you've set a custom default domain for your cluster, you may need to remove or change this DNS suffix for service resolution to work correctly.                                                                                                                                                                                           | `".svc.cluster.local"`                                                                                                        |
| devurls.host                           | string | Should be a wildcard hostname to allow matching against custom-created dev URLs. Leaving as an empty string results in devurls being disabled. Example: "\*.devurls.coder.com".                                                                                                                                                              | `""`                                                                                                                          |
| dockerd.image                          | string | Injected during releases.                                                                                                                                                                                                                                                                                                                    | `""`                                                                                                                          |
| envbuilder.image                       | string | Injected during releases.                                                                                                                                                                                                                                                                                                                    | `""`                                                                                                                          |
| environments.tolerations               | list   | Tolerations are applied to all user environments. Each element is a regular pod toleration object. To set service tolerations see serviceTolerations. See values.yaml for an example.                                                                                                                                                        | `[]`                                                                                                                          |
| envproxy.image                         | string | Injected during releases.                                                                                                                                                                                                                                                                                                                    | `""`                                                                                                                          |
| envproxy.replicas                      | int    | The number of replicas to run of the envproxy.                                                                                                                                                                                                                                                                                               | `1`                                                                                                                           |
| envproxy.resources                     | object | Kubernetes resource request and limits for envproxy pods. To unset a value, set it to "". To unset all values, you can provide a values.yaml file which sets resources to nil. See values.yaml for an example.                                                                                                                               | `{"limits":{"cpu":"250m","memory":"512Mi"},"requests":{"cpu":"250m","memory":"512Mi"}}`                                       |
| envproxy.terminationGracePeriodSeconds | int    | Amount of seconds to wait before shutting down the environment proxy if there are still open connections. This is set very long intentionally so developers do not deal with disconnects during deployments.                                                                                                                                 | `14400`                                                                                                                       |
| imagePullPolicy                        | string | Sets the policy for pulling a container image across all services.                                                                                                                                                                                                                                                                           | `"Always"`                                                                                                                    |
| ingress.additionalAnnotations          | list   | Additional annotations to be used when creating the ingress. These can be used to specify certificate issuers or other cloud provider specific integrations. Annotations are provided as strings e.g. [ "mykey:myvalue", "mykey2:myvalue2" ]                                                                                                 | `[]`                                                                                                                          |
| ingress.host                           | string | The hostname to use for accessing the platform. This can be left blank and the user can still access the platform from the external IP or a DNS name that resolves to the external IP address.                                                                                                                                               | `""`                                                                                                                          |
| ingress.podSecurityPolicyName          | string | The name of the pod security policy the built in ingress controller should abide. It should be noted that the ingress controller requires the `NET_BIND_SERVICE` capability, privilege escalation, and access to privileged ports to successfully deploy.                                                                                    | `""`                                                                                                                          |
| ingress.tls                            | object | TLS options for the ingress. The hosts used for the tls configuration come from the ingress.host and the devurls.host variables. If those don't exist, then the TLS configuration will be ignored.                                                                                                                                           | `{"devurlsHostSecretName":"","enable":false,"hostSecretName":""}`                                                             |
| ingress.tls.devurlsHostSecretName      | string | The secret to use for the devurls.host hostname.                                                                                                                                                                                                                                                                                             | `""`                                                                                                                          |
| ingress.tls.enable                     | bool   | Enables the tls configuration.                                                                                                                                                                                                                                                                                                               | `false`                                                                                                                       |
| ingress.tls.hostSecretName             | string | The secret to use for the ingress.host hostname.                                                                                                                                                                                                                                                                                             | `""`                                                                                                                          |
| ingress.useDefault                     | bool   | If set to true will deploy an nginx ingress that will allow you to access Coder from an external IP address, but if your kubernetes cluster is configured to provision external IP addresses. If you would like to bring your own ingress and hook Coder into that instead, set this value to false.                                         | `true`                                                                                                                        |
| logging.human                          | string | Where to send logs that are formatted for readability by a human. Set to an empty string to disable.                                                                                                                                                                                                                                         | `"/dev/stderr"`                                                                                                               |
| logging.json                           | string | Where to send logs that are formatted as JSON. Set to an empty string to disable.                                                                                                                                                                                                                                                            | `""`                                                                                                                          |
| logging.stackdriver                    | string | Where to send logs that are formatted for Google Stackdriver. Set to an empty string to disable.                                                                                                                                                                                                                                             | `""`                                                                                                                          |
| namespaceWhitelist                     | list   | A list of additional namespaces that environments may be deploy to.                                                                                                                                                                                                                                                                          | `[]`                                                                                                                          |
| podSecurityPolicyName                  | string | The name of the pod security policy to apply to all Coder services and user environments. The optional ingress has its own field for pod security policy as well.                                                                                                                                                                            | `""`                                                                                                                          |
| postgres.database                      | string | The name of the database that coder will use. It must exist before Coder is installed.                                                                                                                                                                                                                                                       | `""`                                                                                                                          |
| postgres.host                          | string | The host of the external postgres instance.                                                                                                                                                                                                                                                                                                  | `""`                                                                                                                          |
| postgres.passwordSecret                | string | The name of an existing secret in the current namespace with the password to the Postgres instance. The password must be contained in the secret field `password`. This should be set to an empty string if the database does not require a password to connect.                                                                             | `""`                                                                                                                          |
| postgres.port                          | string | The port of the external postgres instance.                                                                                                                                                                                                                                                                                                  | `""`                                                                                                                          |
| postgres.sslMode                       | string | Determines how the connection is made to the database. The acceptable values are: `disable`, `allow`, `prefer`, `require`, `verify-ca`, and `verify-full`.                                                                                                                                                                                   | `"require"`                                                                                                                   |
| postgres.useDefault                    | bool   | Deploys an internal Postgres instance alongside the platform. It is not recommended to run the internal Postgres instance in production. If true, all other values are ignored.                                                                                                                                                              | `true`                                                                                                                        |
| postgres.user                          | string | the user of the external postgres instance.                                                                                                                                                                                                                                                                                                  | `""`                                                                                                                          |
| serviceTolerations                     | list   | Tolerations are applied to all Coder managed services. Each element is a toleration object. To set user environment tolerations see environments.tolerations. See values.yaml for an example.                                                                                                                                                | `[]`                                                                                                                          |
| serviceType                            | string | See the following for the different serviceType options and their use: https://kubernetes.io/docs/concepts/services-networking/service/#publishing-services-service-types                                                                                                                                                                    | `"ClusterIP"`                                                                                                                 |
| ssh.enable                             | bool   | Enables accessing environments via SSH.                                                                                                                                                                                                                                                                                                      | `true`                                                                                                                        |
| storageClassName                       | string | Sets the storage class for all Coder services and user environments. By default the storageClassName is not specified and thus the default StorageClass is used. If storageClassName is not specified and a default StorageClass does not exist, then the deployment will fail. The storageClass MUST support the ReadWriteOnce access mode. | `""`                                                                                                                          |
| timescale                              | object | Contains configuration for the internal database. It is not recommended to run this service in production. See the `postgres` section for connecting to an external Postgres database.                                                                                                                                                       | `{"image":"","resources":{"limits":{"cpu":"250m","memory":"1Gi"},"requests":{"cpu":"250m","memory":"1Gi","storage":"10Gi"}}}` |
| timescale.image                        | string | Injected during releases.                                                                                                                                                                                                                                                                                                                    | `""`                                                                                                                          |
| timescale.resources                    | object | Kubernetes resource request and limits for the timescale pod. To unset a value, set it to "". To unset all values, you can provide a values.yaml file which sets resources to nil. See values.yaml for an example.                                                                                                                           | `{"limits":{"cpu":"250m","memory":"1Gi"},"requests":{"cpu":"250m","memory":"1Gi","storage":"10Gi"}}`                          |
| timescale.resources.requests.storage   | string | Specifies the size of the volume claim for persisting the database.                                                                                                                                                                                                                                                                          | `"10Gi"`                                                                                                                      |

## Contributing

All of the Helm charts for Coder services are contained in the
[`templates`][template-folder] folder.

Template values and documentation are in the [`values.yaml`][values-file] file.
The values table in the README is parsed from that file.

This README file is generated from [`README.md.gotmpl`][readme-template-file].
When adding content to the README make sure to update that file.

When adding a value or updating the README template make sure to document it
and run [`gen-readme.sh`][gen-readme-file]. This will update and format the
values table and additional docs.

[helm-repo]: https://helm.coder.com/
[template-folder]: https://github.com/cdr/enterprise-helm/tree/master/templates
[values-file]: https://github.com/cdr/enterprise-helm/blob/master/values.yaml
[readme-template-file]: https://github.com/cdr/enterprise-helm/blob/master/README.md.gotmpl
[gen-readme-file]: https://github.com/cdr/enterprise-helm/blob/master/gen-readme.sh
