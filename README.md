# cani.tf

As time passes, OpenTofu and Terraform become more distant from each other. [CanI.TF](https://cani.tf) helps us to understand their differences quickly.

## Add new features

The `tools/tools.yaml` file is the source of truth. Then run `go run main.go` to validate and regenerate `data/tools.json` and `static/tools.json`.

### `tools/tools.yaml` structure

The file has two top-level keys:

#### `tools`

Metadata for each tool. Displayed in the header rows of the comparison table.

```yaml
tools:
  opentofu:
    version: "1.11"
    versionURL: https://github.com/opentofu/opentofu/releases/latest
    license: MPL-2.0
    licenseURL: https://github.com/opentofu/opentofu/blob/main/LICENSE
    registry: https://search.opentofu.org/
```

#### `features`

Each tool entry under a feature supports three optional fields:

| Field               | Description                                                                 |
|---------------------|-----------------------------------------------------------------------------|
| `version`           | The version that introduced the feature (e.g. `"1.9"`).                    |
| `url`               | URL to the feature's documentation. Renders the version as a link.         |
| `featureRequestURL` | URL to an open issue or RFC tracking the feature. Renders as a `?` link when the feature is not yet supported (i.e. no `version`). |

**Feature supported by both tools:**
```yaml
- name: State encryption
  tools:
    opentofu:
      version: "1.7"
      url: https://opentofu.org/docs/language/state/encryption/
    terraform:
      version: "1.9"
      url: https://developer.hashicorp.com/terraform/language/state/encryption
```

**Feature supported by one tool only (renders `-` for the other):**
```yaml
- name: State encryption
  tools:
    opentofu:
      version: "1.7"
      url: https://opentofu.org/docs/language/state/encryption/
```

**Feature not yet supported, with a tracking link:**
```yaml
- name: Some upcoming feature
  tools:
    opentofu:
      version: "1.9"
      url: https://opentofu.org/docs/...
    terraform:
      featureRequestURL: https://github.com/hashicorp/terraform/issues/123
```

## Local development

```
$ # Install hugo
$ asdf plugin add hugo
$ asdf install
$ go run main.go
$ hugo server --buildDrafts --disableFastRender
$ # HTML templates are in layouts/_default
```
