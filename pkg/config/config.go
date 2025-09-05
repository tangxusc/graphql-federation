// Package config provides configuration-related functionality
package config

// PluginConfig represents the configuration for the GraphQL plugin
type PluginConfig struct {
	MockEnable bool `json:"mockEnable" yaml:"mockEnable"`

	// Schema configuration
	Schema SchemaConfig `json:"schema" yaml:"schema"`

	// Logging configuration
	Log LogConfig `json:"log" yaml:"log"`
}

// SchemaConfig represents the GraphQL schema configuration
type SchemaConfig struct {
	Path          string `json:"path" yaml:"path"`
	Introspection bool   `json:"introspection" yaml:"introspection"`
}

// LogConfig represents the logging configuration
type LogConfig struct {
	Level  string `json:"level" yaml:"level"`
	Format string `json:"format" yaml:"format"`
}
