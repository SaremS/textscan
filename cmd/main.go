package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	backend "github.com/sarems/textscan/internal/backend"
	core "github.com/sarems/textscan/internal/core"
	frontend "github.com/sarems/textscan/internal/frontend"
	input "github.com/sarems/textscan/internal/input"
)

var (
	provider          string
	apiKey            string
	maxDepth          uint8
	probabilityCutoff float64
	questions         []string
	filetype          string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "textscan <file-or-directory|->",
		Short: "Find things in text with JEV",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runScan(args[0])
		},
	}

	rootCmd.Flags().StringVar(
		&provider,
		"provider",
		"vercel",
		"Backend provider: vercel or typesafe",
	)

	rootCmd.Flags().StringVar(
		&apiKey,
		"api-key",
		"",
		"API key (defaults to JEV_API_KEY)",
	)

	rootCmd.Flags().Uint8Var(
		&maxDepth,
		"max-depth",
		3,
		"Maximum scan depth",
	)

	rootCmd.Flags().Float64Var(
		&probabilityCutoff,
		"probability-cutoff",
		0.75,
		"Risk cutoff value between 0 and 1",
	)

	rootCmd.Flags().StringSliceVarP(
		&questions,
		"question",
		"q",
		nil,
		"Search question for JEV Noul type",
	)
	rootCmd.Flags().StringVar(
		&filetype,
		"filetype",
		"",
		"Required stdin file type: go, js, py, or txt",
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runScan(filename string) error {
	var reader core.TextReader
	var err error
	if filename == "-" {
		reader, err = input.NewStdinReader(os.Stdin, filetype)
	} else {
		reader, err = input.NewFileReader(filename)
	}
	if err != nil {
		return fmt.Errorf("create input reader: %w", err)
	}

	key := apiKey
	if key == "" {
		key = os.Getenv("JEV_API_KEY")
	}

	if key == "" {
		return fmt.Errorf(
			"API key is required: use --api-key or set JEV_API_KEY",
		)
	}

	scanBackend, err := backend.NewJevCodeBackend(
		provider,
		key,
		maxDepth,
		probabilityCutoff,
	)
	if err != nil {
		return fmt.Errorf("create backend: %w", err)
	}

	outputFrontend := frontend.NewAsyncConsoleOutputFrontend(probabilityCutoff)

	scanner := core.NewTextScanner(
		reader,
		scanBackend,
		outputFrontend,
	)

	scanner.Scan(questions)

	return nil
}
