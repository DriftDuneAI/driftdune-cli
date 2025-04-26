package main

import (
	"log"

	"github.com/spf13/cobra"
)

// rootCmd is the base command for driftdune
var rootCmd = &cobra.Command{
	Use:   "driftdune",
	Short: "DriftDune CLI: detect & remediate semantic drift",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("ERROR: %v", err)
	}
}
