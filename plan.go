package main

import (
	"fmt"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/spf13/cobra"
)

type Policy struct {
	Model          string  `hcl:"model"`
	Baseline       string  `hcl:"baseline"`
	PromptSuite    string  `hcl:"prompt_suite"`
	AlertThreshold float64 `hcl:"alert_threshold"`
}

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Show what would happen (dry-run)",
	RunE: func(cmd *cobra.Command, args []string) error {
		var pol Policy
		if err := hclsimple.DecodeFile("driftdune.hcl", nil, &pol); err != nil {
			return err
		}
		fmt.Println("🔍 DriftDune Plan:")
		fmt.Printf(" • Model:         %s\n", pol.Model)
		fmt.Printf(" • Baseline:      %s\n", pol.Baseline)
		fmt.Printf(" • Prompt Suite:  %s\n", pol.PromptSuite)
		fmt.Printf(" • Threshold:     %.3f\n", pol.AlertThreshold)
		fmt.Println("No real changes—just showing you the policy.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
}
