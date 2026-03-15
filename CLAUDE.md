# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**CanI.TF** is a static website comparing features between OpenTofu and Terraform. It exposes both an HTML comparison table and a JSON API at `/tools.json`.

## Development Workflow

Adding or updating features requires two steps:

1. Edit `tools/tools.yaml` — this is the source of truth (feature-centric: each feature lists per-tool data under a `tools:` map)
2. Run `go run main.go` to validate and regenerate `data/tools.json` and `static/tools.json`

To preview the site locally:
```bash
hugo server --buildDrafts --disableFastRender
```

Hugo version is pinned via `.tool-versions` (hugo 0.135.0). Install with `asdf install`.

## Architecture

- **`tools/tools.yaml`** — YAML source of truth. Feature-centric: a top-level `tools:` map holds metadata per tool, and a top-level `features:` list holds each feature with per-tool data under a `tools:` map. Omitting a tool key means that tool doesn't support the feature.
- **`main.go`** — reads `tools/tools.yaml`, defines the JSON schema, validates the data, and writes `data/tools.json` + `static/tools.json`.
- **`data/tools.json`** — consumed by Hugo templates at build time.
- **`static/tools.json`** — served directly as the public API.
- **`layouts/`** — Hugo templates. `_default/baseof.html` is the shell; `index.html` renders the comparison table by iterating the tools data.

## Data Schema

Each tool entry has: `version`, `versionURL`, `license`, `licenseURL`, `registry`, and `features[]`. Each feature has `name` (required), and optional `version`, `url`, and `featureRequestURL`. Feature names are matched by string equality across tools to build the comparison table — a feature present in one tool but not the other renders as `-`. A `featureRequestURL` with no `version` renders as a `?` link.

## Deployment

Pushing to `main` triggers `.github/workflows/hugo.yml`, which builds with Hugo and deploys to GitHub Pages automatically.
