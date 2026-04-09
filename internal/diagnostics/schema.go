// Package diagnostics provides introspection tools for the CLI, including
// command schema inspection that shows how each CLI command maps to backend
// API endpoints and whether an extension override is configured.
package diagnostics

import (
	"fmt"
	"sort"

	"github.com/openedx/cli/internal/config"
	"github.com/openedx/cli/internal/registry"
)

// CommandSchema describes how a command is resolved for display.
type CommandSchema struct {
	CommandKey   string   `json:"command"`
	PublicMethod string   `json:"public_method"`
	PublicPath   string   `json:"public_path"`
	HasExtension bool     `json:"has_extension"`
	ExtensionURL string   `json:"extension_url,omitempty"`
	RequiredArgs []string `json:"required_args"`
	OutputModel  string   `json:"output_model"`
}

// GetCommandSchema looks up a single command in the registry and checks
// whether an extension override exists in the provided extensions map.
// Returns an error if the command key is not found in the registry.
func GetCommandSchema(cmdKey string, extensions map[string]config.ExtensionMapping) (*CommandSchema, error) {
	reg := registry.LatestRegistry()

	meta, ok := reg[cmdKey]
	if !ok {
		return nil, fmt.Errorf("unknown command: %s", cmdKey)
	}

	schema := &CommandSchema{
		CommandKey:   meta.Key,
		PublicMethod: meta.Method,
		PublicPath:   meta.Path,
		RequiredArgs: meta.RequiredArgs,
		OutputModel:  meta.OutputModel,
	}

	if ext, found := extensions[cmdKey]; found {
		schema.HasExtension = true
		schema.ExtensionURL = ext.URL
	}

	// Ensure required_args is never nil in JSON output (empty slice instead).
	if schema.RequiredArgs == nil {
		schema.RequiredArgs = []string{}
	}

	return schema, nil
}

// GetAllCommandSchemas returns schemas for all registered commands, sorted
// alphabetically by command key for deterministic output.
func GetAllCommandSchemas(extensions map[string]config.ExtensionMapping) ([]CommandSchema, error) {
	reg := registry.LatestRegistry()

	schemas := make([]CommandSchema, 0, len(reg))
	for key, meta := range reg {
		schema := CommandSchema{
			CommandKey:   meta.Key,
			PublicMethod: meta.Method,
			PublicPath:   meta.Path,
			RequiredArgs: meta.RequiredArgs,
			OutputModel:  meta.OutputModel,
		}

		if ext, found := extensions[key]; found {
			schema.HasExtension = true
			schema.ExtensionURL = ext.URL
		}

		if schema.RequiredArgs == nil {
			schema.RequiredArgs = []string{}
		}

		schemas = append(schemas, schema)
	}

	// Sort for deterministic output.
	sort.Slice(schemas, func(i, j int) bool {
		return schemas[i].CommandKey < schemas[j].CommandKey
	})

	return schemas, nil
}
