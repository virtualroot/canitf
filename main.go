package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/qri-io/jsonschema"
	"gopkg.in/yaml.v3"
)

var schemaData = []byte(`{
  "$schema": "https://json-schema.org/draft/2019-09/schema",
  "type": "object",
  "properties": {
    "opentofu": { "$ref": "#/$defs/tool" },
    "terraform": { "$ref": "#/$defs/tool" }
  },
  "required": ["opentofu", "terraform"],
  "$defs": {
    "tool": {
      "type": "object",
      "properties": {
        "version": {
          "type": "string",
          "maxLength": 4
        },
        "license": {
          "type": "string"
        },
        "licenseURL": {
          "type": "string",
          "format": "uri"
        },
        "registry": {
          "type": "string",
          "format": "uri"
        },
        "features": {
          "type": "array",
          "items": {
            "type": "object",
            "properties": {
              "name": {
                "type": "string"
              },
              "url": {
                "type": "string",
                "format": "uri"
              },
              "version": {
                "type": "string",
                "maxLength": 4
              },
              "featureRequestURL": {
                "type": "string",
                "format": "uri"
              }
            },
            "required": ["name"]
          }
        }
      },
      "required": ["version", "versionURL", "license", "licenseURL", "registry", "features"]
    }
  }
}`)

type Feature struct {
	Name              string `json:"name" yaml:"name"`
	Version           string `json:"version,omitempty" yaml:"version,omitempty"`
	URL               string `json:"url,omitempty" yaml:"url,omitempty"`
	FeatureRequestURL string `json:"featureRequestURL,omitempty" yaml:"featureRequestURL,omitempty"`
}

type Tool struct {
	Version    string    `json:"version" yaml:"version"`
	VersionURL string    `json:"versionURL" yaml:"versionURL"`
	License    string    `json:"license" yaml:"license"`
	LicenseURL string    `json:"licenseURL" yaml:"licenseURL"`
	Registry   string    `json:"registry" yaml:"registry"`
	Features   []Feature `json:"features" yaml:"features"`
}

type ToolMeta struct {
	Version    string `yaml:"version"`
	VersionURL string `yaml:"versionURL"`
	License    string `yaml:"license"`
	LicenseURL string `yaml:"licenseURL"`
	Registry   string `yaml:"registry"`
}

type FeatureToolData struct {
	Version           string `yaml:"version,omitempty"`
	URL               string `yaml:"url,omitempty"`
	FeatureRequestURL string `yaml:"featureRequestURL,omitempty"`
}

type FeatureEntry struct {
	Name  string                      `yaml:"name"`
	Tools map[string]*FeatureToolData `yaml:"tools"`
}

type ToolsFile struct {
	Tools    map[string]ToolMeta `yaml:"tools"`
	Features []FeatureEntry      `yaml:"features"`
}

func main() {
	data, err := os.ReadFile("tools/tools.yaml")
	if err != nil {
		panic("read tools/tools.yaml: " + err.Error())
	}
	var tf ToolsFile
	if err := yaml.Unmarshal(data, &tf); err != nil {
		panic("parse tools/tools.yaml: " + err.Error())
	}

	tools := make(map[string]Tool, len(tf.Tools))
	for name, meta := range tf.Tools {
		tools[name] = Tool{
			Version:    meta.Version,
			VersionURL: meta.VersionURL,
			License:    meta.License,
			LicenseURL: meta.LicenseURL,
			Registry:   meta.Registry,
		}
	}
	for _, f := range tf.Features {
		for toolName, fd := range f.Tools {
			t := tools[toolName]
			t.Features = append(t.Features, Feature{Name: f.Name, Version: fd.Version, URL: fd.URL, FeatureRequestURL: fd.FeatureRequestURL})
			tools[toolName] = t
		}
	}

	toolsJSON, err := json.MarshalIndent(tools, "", "  ")
	if err != nil {
		panic("marshal tools: " + err.Error())
	}

	ctx := context.Background()
	rs := &jsonschema.Schema{}
	if err := json.Unmarshal(schemaData, rs); err != nil {
		panic("unmarshal schema: " + err.Error())
	}

	errs, err := rs.ValidateBytes(ctx, toolsJSON)
	if err != nil {
		panic(err)
	}
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Println(e.Error())
		}
		os.Exit(1)
	}

	output := make(map[string]interface{}, len(tools)+1)
	for k, v := range tools {
		output[k] = v
	}
	output["lastUpdated"] = time.Now().UTC().Format("2006-01-02")

	outputJSON, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		panic("marshal output: " + err.Error())
	}

	// File for Hugo to template the table
	if err := os.WriteFile("data/tools.json", outputJSON, 0644); err != nil {
		panic(err)
	}

	// allow access to https://cani.tf/tools.json
	if err := os.WriteFile("static/tools.json", outputJSON, 0644); err != nil {
		panic(err)
	}
}
