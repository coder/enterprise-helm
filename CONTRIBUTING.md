# How to contribute

This guide describes our conventions and practices.
If you would like to contribute to this project but don't know how, please consider adding unit tests, examples, or documentation.

## Project structure

THis project is organized as follows:

* [.github](.github): configuration for GitHub, including builds
* [examples](examples): a collection of common installation tasks presented as examples, with comments; these are also checked as part of the build process
* [scripts](scripts): a collection of scripts, used mainly for building the Chart
* [templates](templates): the source templates for the Chart
* [tests](tests): unit tests, written in Go

## Tests

We test changes to the Helm Chart in the following ways:

1. We have a suite of **unit tests**, which are written in Go.
   This test suite ensures that the Kubernetes objects that Helm will try to create match what we expect.
   Unit tests are the best way to ensure that a given feature continues to work correctly, particularly when we make changes to the Chart.
   When implementing or modifying a feature, please try to add corresponding unit test coverage.

1. We render all of the **examples** into Kubernetes objects.
   These tests check that the Kubernetes objects that Helm generates for our example files pass static code analysis checks.
   We use tools like `kube-linter` to ensure that we are following best practices when possible.

1. We perform a default installation using [kind test clusters](https://kind.sigs.k8s.io/).
   This runs a minimal Kubernetes installation inside containers, then installs Coder into that cluster.
   This ensures that the Kubernetes API server considers the manifest to be valid.

### Unit testing guide

The unit tests load the Chart from the current directory using [Helm's Go SDK](https://pkg.go.dev/helm.sh/helm/v3).
The [`values.go`](./tests/values.go) file defines the acceptable Helm values with the `CoderValues` type representing the root of the file.
Each test invokes `LoadChart`, which loads the unpacked Chart (templates, values, and settings defined in [Chart.yaml](Chart.yaml)).
Using the resulting `*Chart` type, you can invoke the `Render` function, which accepts a function that will receive the current Chart values (as `*CoderValues`), and can modify those values in place.
`MustRender` is a convenience method, similar to `Render`, but is less configurable and will `panic` on failure.

Here is an [example of a simple test](./tests/deployment_test.go), annotated with what it is doing:

```go
func TestDeployment(t *testing.T) {
    // Run this test in parallel with other tests.
    // Since each Chart is a separate instance, we can safely run tests in parallel.
	t.Parallel()

    // Load the chart from the root directory.
    // This requires the *testing.T object so that it can output errors and stop
    // the tests if failures occur.
	chart := LoadChart(t)

    // Create a sub-test called Labels
	t.Run("Labels", func(t *testing.T) {
		var (
            // The Helm chart will apply some labels by default, according
            // to the release name, namespace, component, and other details.
            // This map shows the expected set of default labels. 
			expectedLabels = map[string]string{
				"app":                         "coderd",
				"app.kubernetes.io/name":      "coder",
				"app.kubernetes.io/part-of":   "coder",
				"app.kubernetes.io/component": "coderd",
				"app.kubernetes.io/instance":  "coder",
				"coder.deployment":            "coderd",
			}
            // extraLabels is a map of additional labels that we will add,
            // which corresponds to the `coderd.extraLabels` setting.
			extraLabels = map[string]string{
				"foo": "bar",
			}

            // MustRender will provide the callback with a clean copy of
            // CoderValues, that the test will modify in order to change
            // the ExtraLabels setting. MustRender then returns the
            // resulting set of objects. It is safe to run MustRender
            // multiple times on the same Chart object.
			objs = chart.MustRender(t, func(cv *CoderValues) {
				cv.Coderd.ExtraLabels = extraLabels
			})

            // MustFindDeployment is a utility function that finds and
            // returns the Deployment kind object with the name "coderd"
			coderd = MustFindDeployment(t, objs, "coderd")
		)

		for k, v := range extraLabels {
			if _, found := expectedLabels[k]; !found {
				expectedLabels[k] = v
			}
		}

        // We compare our expected labels with the labels from the
        // resulting deployment to ensure that they match.
		require.EqualValues(t, expectedLabels, coderd.Spec.Template.Labels)
	})
}
```

## Documentation

If you add or modify the settings in the [`values.yaml`](values.yaml) file, add a comment to that file which describes the purpose and behavior of each field.
We use `helm-docs` to generate the [`README.md`](README.md) file based on a template ([`README.md.gotmpl`](README.md.gotmpl)) file, which you can run using `make README.md` from the root of the project.

## Feature lifecycle

Please label experimental features in the documentation, so that users know what to expect for long-term support.

Deprecation notices or other important notes should be added to [`templates/NOTES.txt`](https://helm.sh/docs/chart_template_guide/notes_files/), as these are visible when running `helm install` or `helm upgrade`.
