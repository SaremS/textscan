package backend

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/sarems/textscan/internal/core"

	sitter "github.com/smacker/go-tree-sitter"
	jev "github.com/stumble/jev-go"
)

type JevCodeBackend struct {
	client                 *jev.Client
	maxDepth               uint8   //only go `maxDepth` deep in AST
	probabilityCutoffValue float64 //don't continue AST traversal if risk at node is less
}

func NewJevCodeBackend(provider, apiKey string, maxDepth uint8, probabilityCutOffValue float64) (*JevCodeBackend, error) {
	if maxDepth > 10 {
		return nil, fmt.Errorf("maxDepth should be less than 10, but is %d", maxDepth)
	}

	if probabilityCutOffValue < 0.0 || probabilityCutOffValue > 1.0 {
		return nil, fmt.Errorf("riskCutOffValue should be >=0, <=1, but is %f", probabilityCutOffValue)
	}

	switch provider {
	case "vercel":
		client, err := jev.NewClient(jev.Config{
			Provider: jev.ProviderVercel,
			APIKey:   apiKey,
		},
		)

		if err != nil {
			return nil, err
		}

		return &JevCodeBackend{
			client:                 client,
			maxDepth:               maxDepth,
			probabilityCutoffValue: probabilityCutOffValue,
		}, nil
	case "typesafe":
		client, err := jev.NewClient(jev.Config{
			APIKey: apiKey,
		},
		)

		if err != nil {
			return nil, err
		}

		return &JevCodeBackend{
			client:                 client,
			maxDepth:               maxDepth,
			probabilityCutoffValue: probabilityCutOffValue,
		}, nil

	default:
		return nil, fmt.Errorf("Provider must be 'vercel' or 'typesafe' but got '%s'", provider)

	}
}

func (j JevCodeBackend) ScanText(textItems []core.TextItem) <-chan core.ProbsRange {
	probsRangeChan := make(chan core.ProbsRange, 5)

	go func() {
		var wg sync.WaitGroup

		for i := range textItems {
			item := &textItems[i]

			wg.Go(func() {
				if err := j.doScore(item, probsRangeChan); err != nil {
					// handle error
				}
			})
		}

		wg.Wait()
		close(probsRangeChan)
	}()

	return probsRangeChan
}

func (j JevCodeBackend) inferLanguage(item *core.TextItem) (*core.ProgrammingLanguage, error) {
	ext := filepath.Ext(string(item.Identifier))
	var language core.ProgrammingLanguage

	switch ext {
	case ".go":
		language = core.Go
	case ".js":
		language = core.JavaScript
	case ".py":
		language = core.Python
	default:
		language = ""
	}

	return &language, nil
}

func makeProbMap(items []string) map[string]float64 {
	set := make(map[string]float64, len(items))

	for _, item := range items {
		set[item] = 1.0
	}

	return set
}

func (j JevCodeBackend) doScore(item *core.TextItem, probsRangeChan chan<- core.ProbsRange) error {
	text := item.Text
	language, err := j.inferLanguage(item)
	if err != nil {
		return err
	}

	tree, err := parseCodeFromText(text, *language)
	if err != nil {
		return err
	}

	rootNode := tree.RootNode()

	j.scoreRecursively(
		item,
		makeProbMap(item.Questions),
		rootNode,
		0,
		probsRangeChan,
	)

	return nil
}

type workResult struct {
	newChildNode          *sitter.Node
	questionProbabilities map[string]float64
}

type workResultOrError struct {
	result *workResult
	err    error
}

func (j JevCodeBackend) scoreRecursively(
	item *core.TextItem,
	parentProbabilities map[string]float64,
	node *sitter.Node,
	currentDepth uint8,
	probsRangeChan chan<- core.ProbsRange,
) {
	workChannel := j.doWork(item, item.Text, node, parentProbabilities)

	dropCounts := make(map[string]int)
	totalCount := 0

	for resultItem := range workChannel {
		totalCount++
		if resultItem.err != nil {
			continue
		}

		result := resultItem.result
		currentProbabilities := result.questionProbabilities

		newProbabilities := multiplyProbilities(parentProbabilities, currentProbabilities)

		remainingProbabilities := j.updateCutoff(newProbabilities)

		if currentDepth+1 <= j.maxDepth && len(remainingProbabilities) > 0 && result.newChildNode.NamedChildCount() > 0 {
			j.scoreRecursively(item, remainingProbabilities, result.newChildNode, currentDepth+1, probsRangeChan)
		} else {
			j.emitResult(item, remainingProbabilities, result.newChildNode, probsRangeChan)
		}

		for k := range parentProbabilities {
			if _, ok := remainingProbabilities[k]; !ok {
				dropCounts[k] = dropCounts[k] + 1
			}
		}
	}

	//emit parent probabilities if they were dropped in all children
	for k, count := range dropCounts {
		if count == totalCount && parentProbabilities[k] < 1.0 {
			singleEmitProbabilities := map[string]float64{k: parentProbabilities[k]}
			j.emitResult(item, singleEmitProbabilities, node, probsRangeChan)
		}
	}
}

func (j JevCodeBackend) updateCutoff(newProbabilities map[string]float64) map[string]float64 {
	remainingProbabilities := make(map[string]float64, len(newProbabilities))

	for k, prob := range newProbabilities {
		if prob >= j.probabilityCutoffValue {
			remainingProbabilities[k] = prob
		}
	}
	return remainingProbabilities
}

func (j JevCodeBackend) emitResult(item *core.TextItem, probabilities map[string]float64, node *sitter.Node, probsRangeChan chan<- core.ProbsRange) {
	for k, prob := range probabilities {
		probsRangeChan <- core.ProbsRange{
			Identifier:      item.Identifier,
			Start:           int(node.StartByte()),
			End:             int(node.EndByte()),
			Probability:     prob,
			StartLineOffset: item.NewLineOffsets[int(node.StartByte())],
			EndLineOffset:   item.NewLineOffsets[int(node.EndByte())-1],
			TextSnippet:     item.Text[int(node.StartByte()):int(node.EndByte())],
			Question:        k,
		}

	}
}

func (j JevCodeBackend) doWork(item *core.TextItem, text core.Text, node *sitter.Node, parentProbabilities map[string]float64) <-chan workResultOrError {

	resultChan := make(chan workResultOrError, 10)
	var wg sync.WaitGroup

	questions := make(map[string]jev.Question, len(item.Questions))
	for k, _ := range parentProbabilities {
		questions[k] = jev.Noul(k)
	}

	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(i)

		startByte := int(child.StartByte())
		endByte := int(child.EndByte())
		targetCode := text[startByte:endByte]

		wg.Go(func() {
			response, err := j.client.Ask(context.Background(), jev.Request{
				State:     map[string]any{"document": targetCode},
				Questions: questions,
			})
			if err != nil {
				resultChan <- workResultOrError{nil, err}
				return
			}

			questionProbabilities := make(map[string]float64, len(item.Questions))
			for k, v := range response.Answers {
				questionProbabilities[k] = v.Noul
			}

			result := workResult{
				newChildNode:          child,
				questionProbabilities: questionProbabilities,
			}

			resultChan <- workResultOrError{&result, nil}
		})
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	return resultChan
}

func multiplyProbilities(left, right map[string]float64) map[string]float64 {
	result := make(map[string]float64, len(left))

	for k, v := range left {
		if rightValue, ok := right[k]; ok {
			result[k] = v * rightValue
		} else {
			result[k] = v
		}
	}

	return result
}
