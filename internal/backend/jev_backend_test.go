package backend

import (
	"context"
	"errors"
	"math"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sarems/textscan/internal/core"
	jev "github.com/stumble/jev-go"
)

type fakeAsker struct {
	mu        sync.Mutex
	requests  []jev.Request
	scores    map[string]float64
	err       error
	inFlight  atomic.Int32
	maxFlight atomic.Int32
	delay     time.Duration
}

func (f *fakeAsker) Ask(_ context.Context, request jev.Request, _ ...jev.CallOptions) (jev.Response, error) {
	current := f.inFlight.Add(1)
	defer f.inFlight.Add(-1)

	for {
		previous := f.maxFlight.Load()
		if current <= previous || f.maxFlight.CompareAndSwap(previous, current) {
			break
		}
	}

	if f.delay > 0 {
		time.Sleep(f.delay)
	}

	f.mu.Lock()
	f.requests = append(f.requests, request)
	f.mu.Unlock()
	if f.err != nil {
		return jev.Response{}, f.err
	}

	document := request.State.(map[string]any)["document"].(string)
	probability, ok := f.scores[document]
	if !ok {
		return jev.Response{}, errors.New("unexpected document: " + document)
	}

	return jev.Response{
		Answers: map[string]jev.Answer{
			"matches_query": {Type: jev.QuestionNoul, Noul: probability},
		},
	}, nil
}

func TestJevBackendGatesSentencesAndCombinesProbabilities(t *testing.T) {
	const text = "High match. Another high match.\n\nLow match."
	client := &fakeAsker{scores: map[string]float64{
		"High match. Another high match.": 0.9,
		"Low match.":                      0.1,
		"High match.":                     0.8,
		" Another high match.":            0.5,
	}}
	backend := &JevBackend{
		client:          client,
		query:           "Is this about a high-priority topic?",
		riskCutoffValue: 0.5,
		workers:         2,
	}
	document := core.NewDocument("notes.txt", text)

	if err := backend.ScanText([]core.Document{*document}); err != nil {
		t.Fatalf("ScanText() error = %v", err)
	}

	scores := document.Scoring.Items()
	assertProbabilityAt(t, scores, strings.Index(text, "High match."), 0.72)
	assertProbabilityAt(t, scores, strings.Index(text, "Another high match."), 0.45)
	assertProbabilityAt(t, scores, strings.Index(text, "Low match."), 0.1)

	client.mu.Lock()
	defer client.mu.Unlock()
	if len(client.requests) != 4 {
		t.Fatalf("request count = %d, want 4", len(client.requests))
	}
	for _, request := range client.requests {
		question := request.Questions["matches_query"]
		if question.Instructions != backend.query {
			t.Fatalf("question = %#v, want query %q", question.Instructions, backend.query)
		}
	}
}

func TestJevBackendBoundsConcurrentRequests(t *testing.T) {
	const text = "One.\n\nTwo.\n\nThree.\n\nFour."
	client := &fakeAsker{
		scores: map[string]float64{
			"One.":   0,
			"Two.":   0,
			"Three.": 0,
			"Four.":  0,
		},
		delay: 10 * time.Millisecond,
	}
	backend := &JevBackend{
		client:          client,
		query:           "topic",
		riskCutoffValue: 1,
		workers:         2,
	}
	document := core.NewDocument("notes.txt", text)

	if err := backend.ScanText([]core.Document{*document}); err != nil {
		t.Fatalf("ScanText() error = %v", err)
	}
	if got := client.maxFlight.Load(); got > 2 {
		t.Fatalf("maximum concurrent requests = %d, want <= 2", got)
	}
}

func TestJevBackendPropagatesRequestErrors(t *testing.T) {
	client := &fakeAsker{err: errors.New("service unavailable")}
	backend := &JevBackend{
		client:          client,
		query:           "topic",
		riskCutoffValue: 0.5,
		workers:         1,
	}
	document := core.NewDocument("notes.txt", "Paragraph.")

	if err := backend.ScanText([]core.Document{*document}); err == nil ||
		!strings.Contains(err.Error(), "service unavailable") {
		t.Fatalf("ScanText() error = %v, want propagated request error", err)
	}
}

func assertProbabilityAt(t *testing.T, scores []float64, index int, want float64) {
	t.Helper()
	if math.Abs(scores[index]-want) > 1e-9 {
		t.Fatalf("score at %d = %f, want %f", index, scores[index], want)
	}
}
