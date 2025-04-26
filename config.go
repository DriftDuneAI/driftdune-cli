// config.go
package main

import "fmt"

// PromptSpec must match the one in apply.go
type PromptSpec struct {
	Name     string `json:"name"`
	Baseline string `json:"baseline"`
	Current  string `json:"current"`
}

// Config holds your pipeline settings.
// Later you can parse an HCL or YAML file here.
type Config struct {
	Model          string
	Baseline       string
	PromptSuite    string
	AlertThreshold float64
	SlackWebhook   string
}

// loadConfig returns defaults for now.
func loadConfig() (*Config, error) {
	cfg := &Config{
		Model:          "gpt-4",
		Baseline:       "gpt-3.5",
		PromptSuite:    "suite.json", // point this at your suite file
		AlertThreshold: 0.05,
		SlackWebhook:   "", // set via flag or env later
	}
	fmt.Println("✅  [Config] Using default settings")
	return cfg, nil
}
