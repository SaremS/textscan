package cmd

import (
	"fmt"
	"os"

	"github.com/sarems/textscan/internal/backend"
	"github.com/sarems/textscan/internal/core"
	"github.com/sarems/textscan/internal/frontend"
	"github.com/sarems/textscan/internal/input"
	"github.com/spf13/cobra"
)

var (
	provider   string
	apiKey     string
	riskCutoff float32
)

func Execute() {
	if err := newRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "textscan <file> <query>",
		Short: "Scan text for a topic with JEV",
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			return runScan(args[0], args[1])
		},
	}

	rootCmd.Flags().StringVar(&provider, "provider", "vercel", "Backend provider: vercel or typesafe")
	rootCmd.Flags().StringVar(&apiKey, "api-key", "", "API key (defaults to JEV_API_KEY)")
	rootCmd.Flags().Float32Var(&riskCutoff, "risk-cutoff", 0.95, "Paragraph risk cutoff between 0 and 1")

	return rootCmd
}

func runScan(filename, query string) error {
	key := apiKey
	if key == "" {
		key = os.Getenv("JEV_API_KEY")
	}
	if key == "" {
		return fmt.Errorf("API key is required: use --api-key or set JEV_API_KEY")
	}

	reader, err := input.NewSingleFileReader(filename)
	if err != nil {
		return fmt.Errorf("create text reader: %w", err)
	}
	scanBackend, err := backend.NewJevBackend(provider, key, query, riskCutoff)
	if err != nil {
		return fmt.Errorf("create backend: %w", err)
	}

	scanner := core.NewScanner(reader, scanBackend, frontend.NewConsoleOutputFrontend())
	return scanner.Scan()
}
