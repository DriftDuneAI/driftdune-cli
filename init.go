package main

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a starter driftdune.hcl policy file",
	RunE: func(cmd *cobra.Command, args []string) error {
		const stub = `# DriftDune policy (HCL)
model       = "gpt-4"
baseline    = "gpt-3.5"
prompt_suite = "suite.yaml"
alert_threshold = 0.05
`
		if _, err := os.Stat("driftdune.hcl"); !errors.Is(err, os.ErrNotExist) {
			return errors.New("driftdune.hcl already exists")
		}
		return ioutil.WriteFile("driftdune.hcl", []byte(stub), 0644)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
