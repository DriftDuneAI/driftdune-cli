// apply.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// PromptSpec defines one prompt’s metadata and file paths.
//type PromptSpec struct {
//	Name     string `json:"name"`
//	Baseline string `json:"baseline"`
//	Current  string `json:"current"`
//}

// PromptResult holds the result for one prompt
type PromptResult struct {
	Name       string  `json:"name"`
	DriftScore float64 `json:"drift_score"`
	Passed     bool    `json:"passed"`
}

// ApplyReport is the full report schema
type ApplyReport struct {
	Timestamp     string         `json:"timestamp"`
	TotalPrompts  int            `json:"total_prompts"`
	PassedCount   int            `json:"passed_count"`
	FailedCount   int            `json:"failed_count"`
	Threshold     float64        `json:"alert_threshold"`
	PromptDetails []PromptResult `json:"prompt_details"`
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Run drift analysis, send alerts, and write apply_report.json",
	RunE:  runApplyCmd,
}

func init() {
	rootCmd.AddCommand(applyCmd)
}

func runApplyCmd(cmd *cobra.Command, args []string) error {
	log.Println("INFO: Stage 1: Loading config")
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("config error: %w", err)
	}
	log.Println("✅ Stage 1 succeeded")

	log.Println("INFO: Stage 2: Loading prompt suite")
	specs, err := readPromptSuite(cfg.PromptSuite)
	if err != nil {
		return fmt.Errorf("reading prompt suite: %w", err)
	}
	log.Printf("✅ Stage 2 succeeded (%d prompts)\n", len(specs))

	log.Println("INFO: Stage 3: Running drift analysis")
	results := analyzePrompts(specs, cfg.AlertThreshold)
	log.Printf("✅ Stage 3 succeeded (passed=%d, failed=%d)\n", results.PassedCount, results.FailedCount)

	if results.FailedCount > 0 && cfg.SlackWebhook != "" {
		log.Println("INFO: Stage 4: Sending alert")
		if err := sendAlert(results, cfg.SlackWebhook); err != nil {
			return fmt.Errorf("alert error: %w", err)
		}
		log.Println("✅ Stage 4 succeeded")
	} else {
		log.Println("INFO: Stage 4: No alerts to send")
	}

	log.Println("INFO: Stage 5: Writing apply_report.json")
	if err := writeApplyReport(results); err != nil {
		return fmt.Errorf("writing report: %w", err)
	}
	log.Println("✅ Stage 5 succeeded")

	log.Printf("All done: 5/5 stages succeeded\n")
	return nil
}

// readPromptSuite loads a JSON array of PromptSpec from disk.
func readPromptSuite(path string) ([]PromptSpec, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var specs []PromptSpec
	if err := json.Unmarshal(data, &specs); err != nil {
		return nil, err
	}
	return specs, nil
}

// analyzePrompts computes a drift score for each prompt and applies the threshold.
func analyzePrompts(specs []PromptSpec, threshold float64) ApplyReport {
	report := ApplyReport{
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
		Threshold:    threshold,
		TotalPrompts: len(specs),
	}
	for _, s := range specs {
		base, err := readEmbedding(s.Baseline)
		if err != nil {
			log.Fatalf("Failed reading %s: %v", s.Baseline, err)
		}
		curr, err := readEmbedding(s.Current)
		if err != nil {
			log.Fatalf("Failed reading %s: %v", s.Current, err)
		}
		score, err := cosineDistance(base, curr)
		if err != nil {
			log.Fatalf("Drift error for %s: %v", s.Name, err)
		}
		passed := score <= threshold
		if passed {
			report.PassedCount++
		} else {
			report.FailedCount++
		}
		report.PromptDetails = append(report.PromptDetails, PromptResult{
			Name:       s.Name,
			DriftScore: math.Round(score*1e4) / 1e4, // 4-decimal precision
			Passed:     passed,
		})
	}
	return report
}

// sendAlert posts a summary of failures to a Slack (or generic) webhook.
func sendAlert(report ApplyReport, webhookURL string) error {
	var names []string
	for _, pr := range report.PromptDetails {
		if !pr.Passed {
			names = append(names, pr.Name)
		}
	}
	payload := map[string]string{
		"text": fmt.Sprintf(
			"DriftDune Alert: %d/%d prompts exceeded threshold %.4f:\n• %s",
			report.FailedCount,
			report.TotalPrompts,
			report.Threshold,
			strings.Join(names, "\n• "),
		),
	}
	body, _ := json.Marshal(payload)

	// retry up to 3 times
	for i := 0; i < 3; i++ {
		resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(body))
		if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil
		}
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("failed to send alert after retries")
}

// writeApplyReport writes the final JSON to disk.
func writeApplyReport(r ApplyReport) error {
	out, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile("apply_report.json", out, 0644)
}

// The utility functions below can be left as before in detect.go / main.go:

// readEmbedding, cosineDistance, etc.
// … (unchanged from your previous detect.go) …
