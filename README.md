# Coder Enterprise Helm

WARNING: The master branch may contain updates for a yet unreleased version of
Coder. The current state of the repo may not represent the latest release.

The Helm package here is only a template for releases: version numbers and
image URIs are injected and bundled into releases uploaded to
[helm.coder.com][helm-repo].

You can pull the official charts with `helm repo add coder
https://helm.coder.com`.

Feel free to open PRs for fixes and improvements you've found while deploying
Coder!

## Values

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
