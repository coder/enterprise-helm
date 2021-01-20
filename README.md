# Workspace providers [BETA]

This chart provides all of the things necessary to run an alpha workspace
provider on a second cluster.

This is only needed for the alpha. Envproxy will become a subchart of the main
Enterprise chart in the beta release of this feature. You probably don't need
this.

## Values

Almost all values are identical to the main chart, but there some changes. A
lot of values that are unneeded were removed, but it should still be safe to
provide them anyways.

#### Added values:
- `envproxy.accessURL` (required if `ingress.host` not set)
- `envproxy.clusterAddress` (required)
- `cemanager.accessURL` (required)
- `cemanager.token` (required, set by coder-cli)
