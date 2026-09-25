// for quick testing
package main

import (
	"os"

	backend "github.com/sarems/textscan/internal/backend"
	core "github.com/sarems/textscan/internal/core"
	frontend "github.com/sarems/textscan/internal/frontend"
	input "github.com/sarems/textscan/internal/input"
)

func main() {
	cutoffProbability := 0.7
	reader, err := input.NewFileReader("example.txt")
	if err != nil {
		panic(err)
	}

	questions := core.QuestionSet([]string{
		"Is this related to Germany?",
		"Is this related to Iceland?",
		"Is this related to AI?",
	})

	scanBackend, err := backend.NewJevCodeBackend(
		"vercel",
		os.Getenv("JEV_API_KEY"),
		1,
		cutoffProbability,
	)
	if err != nil {
		panic(err)
	}

	outputFrontend := frontend.NewAsyncConsoleOutputFrontend(cutoffProbability)

	scanner := core.NewTextScanner(reader, scanBackend, outputFrontend)
	scanner.Scan(questions)
}
