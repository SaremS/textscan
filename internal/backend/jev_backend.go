package backend

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/sarems/textscan/internal/core"
	jev "github.com/stumble/jev-go"
)

const defaultWorkers = 4

type jevAsker interface {
	Ask(context.Context, jev.Request, ...jev.CallOptions) (jev.Response, error)
}

type JevBackend struct {
	client          jevAsker
	query           string
	riskCutoffValue float64
	workers         int
}

func NewJevBackend(
	provider string,
	apiKey string,
	query string,
	riskCutoffValue float32,
) (*JevBackend, error) {
	if strings.TrimSpace(query) == "" {
		return nil, fmt.Errorf("query must not be empty")
	}
	if riskCutoffValue < 0 || riskCutoffValue > 1 {
		return nil, fmt.Errorf("risk cutoff must be between 0 and 1: %f", riskCutoffValue)
	}

	var config jev.Config
	switch provider {
	case "vercel":
		config = jev.Config{Provider: jev.ProviderVercel, APIKey: apiKey}
	case "typesafe":
		config = jev.Config{Provider: jev.ProviderTypeSafe, APIKey: apiKey}
	default:
		return nil, fmt.Errorf("provider must be 'vercel' or 'typesafe', got %q", provider)
	}

	client, err := jev.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &JevBackend{
		client:          client,
		query:           query,
		riskCutoffValue: float64(riskCutoffValue),
		workers:         defaultWorkers,
	}, nil
}

func (j *JevBackend) ScanText(documents []core.Document) error {
	for index := range documents {
		if err := j.scanDocument(&documents[index]); err != nil {
			return err
		}
	}

	return nil
}

type scoredSpan struct {
	span        textSpan
	probability float64
}

func (j *JevBackend) scanDocument(document *core.Document) error {
	paragraphs := splitParagraphs(document.Text)
	paragraphScores, err := j.scoreSpans(paragraphs, document.Text)
	if err != nil {
		return fmt.Errorf("score paragraphs: %w", err)
	}

	var sentenceTasks []sentenceTask
	for _, paragraphScore := range paragraphScores {
		if paragraphScore.probability < j.riskCutoffValue {
			if err := document.Scoring.InitInRange(
				paragraphScore.span.start,
				paragraphScore.span.end,
				paragraphScore.probability,
			); err != nil {
				return err
			}
			continue
		}

		sentences, err := splitSentences(document.Text, paragraphScore.span)
		if err != nil {
			return err
		}
		for _, sentence := range sentences {
			sentenceTasks = append(sentenceTasks, sentenceTask{
				span:                 sentence,
				paragraphProbability: paragraphScore.probability,
			})
		}
	}

	sentenceScores, err := j.scoreSentenceTasks(sentenceTasks, document.Text)
	if err != nil {
		return fmt.Errorf("score sentences: %w", err)
	}

	for _, sentenceScore := range sentenceScores {
		if err := document.Scoring.InitInRange(
			sentenceScore.span.start,
			sentenceScore.span.end,
			sentenceScore.paragraphProbability,
		); err != nil {
			return err
		}
		if err := document.Scoring.MultiplyInRange(
			sentenceScore.span.start,
			sentenceScore.span.end,
			sentenceScore.sentenceProbability,
		); err != nil {
			return err
		}
	}

	return nil
}

func (j *JevBackend) scoreSpans(spans []textSpan, text core.Text) ([]scoredSpan, error) {
	results := make([]scoredSpan, len(spans))
	err := runBounded(len(spans), j.workers, func(index int) error {
		probability, err := j.ask(string(text[spans[index].start:spans[index].end]))
		if err != nil {
			return err
		}
		results[index] = scoredSpan{span: spans[index], probability: probability}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return results, nil
}

type sentenceTask struct {
	span                 textSpan
	paragraphProbability float64
}

type scoredSentence struct {
	span                 textSpan
	paragraphProbability float64
	sentenceProbability  float64
}

func (j *JevBackend) scoreSentenceTasks(
	tasks []sentenceTask,
	text core.Text,
) ([]scoredSentence, error) {
	results := make([]scoredSentence, len(tasks))
	err := runBounded(len(tasks), j.workers, func(index int) error {
		probability, err := j.ask(string(text[tasks[index].span.start:tasks[index].span.end]))
		if err != nil {
			return err
		}
		results[index] = scoredSentence{
			span:                 tasks[index].span,
			paragraphProbability: tasks[index].paragraphProbability,
			sentenceProbability:  probability,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (j *JevBackend) ask(document string) (float64, error) {
	response, err := j.client.Ask(context.Background(), jev.Request{
		State: map[string]any{"document": document},
		Questions: map[string]jev.Question{
			"matches_query": jev.Noul(j.query),
		},
	})
	if err != nil {
		return 0, err
	}

	answer, ok := response.Answers["matches_query"]
	if !ok {
		return 0, fmt.Errorf("JEV response did not include matches_query")
	}

	return answer.Noul, nil
}

func runBounded(count, workers int, work func(int) error) error {
	if count == 0 {
		return nil
	}
	if workers < 1 {
		return fmt.Errorf("worker count must be positive")
	}
	if workers > count {
		workers = count
	}

	jobs := make(chan int)
	var group sync.WaitGroup
	var firstError error
	var errorMutex sync.Mutex

	for range workers {
		group.Add(1)
		go func() {
			defer group.Done()
			for index := range jobs {
				if err := work(index); err != nil {
					errorMutex.Lock()
					if firstError == nil {
						firstError = err
					}
					errorMutex.Unlock()
				}
			}
		}()
	}

	for index := range count {
		jobs <- index
	}
	close(jobs)
	group.Wait()

	return firstError
}
