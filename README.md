# cani.tf

As time passes, OpenTofu and Terraform become more distant from each other. [CanI.TF](https://cani.tf) helps us to understand their differences quickly.

## Add new features

To add new features, edit the JSON in the tool's variable inside `main.go`.
The entry should include the feature's name, the released version, and an optionable URL to its documentation.
If the entry is valid for both tools, use the same name.
Run `go run main.go` to validate the changes against the schema and replace the source of truth `data/tools.json`.

## Local development

```
$ # Install hugo
$ asdf plugin add hugo
$ asdf install
$ # hugo server --buildDrafts --disableFastRender
$ # HTML in layouts/_default
```
