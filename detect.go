package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

var (
	baselinePath string
	currentPath  string
)

// detectCmd is the `detect` subcommand
var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Detect semantic drift between two embedding files",
	Run: func(cmd *cobra.Command, args []string) {
		runDetect(baselinePath, currentPath)
	},
}

func init() {
	// 1) Register detectCmd under the root
	rootCmd.AddCommand(detectCmd)

	// 2) Define its flags
	detectCmd.Flags().StringVarP(&baselinePath, "baseline", "b", "baseline.json", "path to baseline embedding JSON")
	detectCmd.Flags().StringVarP(&currentPath, "current", "c", "current.json", "path to current embedding JSON")
}

// runDetect drives the pipeline: read two JSON files, compute drift, write report
func runDetect(baselinePath, currentPath string) {
	log.Println("INFO: Stage 1: Reading baseline embeddings")
	baseline, err := readEmbedding(baselinePath)
	if err != nil {
		handleError("reading baseline", err)
	}
	log.Println("✅ Stage 1 succeeded")

	log.Println("INFO: Stage 2: Reading current embeddings")
	current, err := readEmbedding(currentPath)
	if err != nil {
		handleError("reading current", err)
	}
	log.Println("✅ Stage 2 succeeded")

	log.Println("INFO: Stage 3: Computing cosine-distance drift score")
	score, err := cosineDistance(baseline, current)
	if err != nil {
		handleError("computing drift", err)
	}
	log.Printf("✅ Stage 3 succeeded (drift_score=%.4f)\n", score)

	log.Println("INFO: Stage 4: Writing report.json")
	report := struct {
		DriftScore float64 `json:"drift_score"`
		Timestamp  string  `json:"timestamp"`
	}{
		DriftScore: score,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	}
	out, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		handleError("marshaling report", err)
	}
	if err := ioutil.WriteFile("report.json", out, 0644); err != nil {
		handleError("writing report.json", err)
	}
	log.Println("✅ Stage 4 succeeded")
	log.Println("All done: 4/4 stages succeeded")
}

// readEmbedding loads a JSON array of floats from disk
func readEmbedding(path string) ([]float64, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var vec []float64
	if err := json.Unmarshal(data, &vec); err != nil {
		return nil, err
	}
	return vec, nil
}

// cosineDistance computes 1 - cosine_similarity(a,b)
func cosineDistance(a, b []float64) (float64, error) {
	if len(a) != len(b) {
		return 0, errors.New("vector lengths mismatch")
	}
	var dot, magA, magB float64
	for i := range a {
		dot += a[i] * b[i]
		magA += a[i] * a[i]
		magB += b[i] * b[i]
	}
	if magA == 0 || magB == 0 {
		return 0, errors.New("zero-magnitude vector")
	}
	return 1 - dot/(math.Sqrt(magA)*math.Sqrt(magB)), nil
}

// handleError logs ❌, prints debug hints, then exits
func handleError(stage string, err error) {
	log.Printf("❌ Stage failed (%s): %v\n", stage, err)
	log.Printf("DEBUG: cwd=%s\n", mustGetwd())
	fmt.Printf("SUGGESTION: run `ls -la` and `cat %s` to inspect files.\n", filepath.Base(stage+".json"))
	os.Exit(1)
}

// mustGetwd returns the current working directory or "unknown"
func mustGetwd() string {
	if cwd, err := os.Getwd(); err == nil {
		return cwd
	}
	return "unknown"
}
