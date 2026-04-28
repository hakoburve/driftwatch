// Package main is the entry point for the driftwatch CLI tool.
// It loads configuration, fetches live service state, and reports drift.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/driftwatch/internal/config"
	"github.com/yourorg/driftwatch/internal/drift"
)

const defaultConfigPath = "driftwatch.yaml"

func main() {
	var (
		configPath = flag.String("config", defaultConfigPath, "path to driftwatch config file")
		verbose    = flag.Bool("verbose", false, "show all keys, including those in sync")
		exitCode   = flag.Bool("exit-code", false, "exit with code 1 if drift is detected")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "driftwatch — detect configuration drift between live services and declared state\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n  driftwatch [flags]\n\nFlags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		os.Exit(2)
	}

	hasDrift := false

	for _, svc := range cfg.Services {
		live, err := fetchLiveConfig(svc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[%s] failed to fetch live config: %v\n", svc.Name, err)
			hasDrift = true
			continue
		}

		results := drift.Detect(svc.Expected, live)
		summary := drift.Summary(results)

		if summary.DriftedCount == 0 && summary.MissingCount == 0 && summary.ExtraCount == 0 {
			fmt.Printf("[%s] ✓ in sync (%d keys checked)\n", svc.Name, summary.TotalKeys)
			if *verbose {
				printResults(svc.Name, results)
			}
			continue
		}

		hasDrift = true
		fmt.Printf("[%s] ✗ drift detected — drifted:%d missing:%d extra:%d (total:%d)\n",
			svc.Name, summary.DriftedCount, summary.MissingCount, summary.ExtraCount, summary.TotalKeys)
		printResults(svc.Name, results)
	}

	if *exitCode && hasDrift {
		os.Exit(1)
	}
}

// fetchLiveConfig retrieves the running configuration for a service.
// Currently reads from the file path specified in svc.LiveSource;
// future versions will support HTTP endpoints and secret managers.
func fetchLiveConfig(svc config.Service) (map[string]string, error) {
	if svc.LiveSource == "" {
		return nil, fmt.Errorf("no live_source configured")
	}

	data, err := os.ReadFile(svc.LiveSource)
	if err != nil {
		return nil, fmt.Errorf("reading live source %q: %w", svc.LiveSource, err)
	}

	parsed, err := config.ParseEnvFile(data)
	if err != nil {
		return nil, fmt.Errorf("parsing live source %q: %w", svc.LiveSource, err)
	}

	return parsed, nil
}

// printResults prints per-key drift details for a service.
func printResults(svcName string, results []drift.Result) {
	for _, r := range results {
		switch r.Status {
		case drift.StatusMatch:
			fmt.Printf("  [%s]   ok       %s\n", svcName, r.Key)
		case drift.StatusDrifted:
			fmt.Printf("  [%s]   drifted  %s: expected=%q live=%q\n", svcName, r.Key, r.Expected, r.Live)
		case drift.StatusMissing:
			fmt.Printf("  [%s]   missing  %s (not present in live config)\n", svcName, r.Key)
		case drift.StatusExtra:
			fmt.Printf("  [%s]   extra    %s (not declared in expected config)\n", svcName, r.Key)
		}
	}
}
