package core

type Identifier string
type Text string
type QuestionSet []string
type ProgrammingLanguage string

const (
	Go         ProgrammingLanguage = "go"
	JavaScript ProgrammingLanguage = "javascript"
	Python     ProgrammingLanguage = "python"
)

type NewLineOffset struct {
	LineStartPos int
	LineEndPos   int
	LineNumber   int
}

type ProbsRange struct {
	Identifier      Identifier
	Start           int
	End             int
	Probability     float64
	StartLineOffset NewLineOffset
	EndLineOffset   NewLineOffset
	TextSnippet     Text
	Question        string
}

type TextItem struct {
	Identifier     Identifier
	Text           Text
	Questions      QuestionSet
	NewLineOffsets []NewLineOffset
	Scoring        TextScoring
	ProbsRanges    []ProbsRange
}

func NewTextItem(identifier Identifier, text Text, questions QuestionSet) (*TextItem, error) {
	scoring, err := newTextScoring(len(questions), []byte(text))
	if err != nil {
		return nil, err
	}

	offsets := getNewLineOffsets(text)
	var probsRanges []ProbsRange

	return &TextItem{
		Identifier:     identifier,
		Text:           text,
		Questions:      questions,
		NewLineOffsets: offsets,
		Scoring:        *scoring,
		ProbsRanges:    probsRanges,
	}, nil
}

func getNewLineOffsets(text Text) []NewLineOffset {
	result := make([]NewLineOffset, len(text))
	currentLine := 1
	currentNewLineOffset := 0

	for pos := range len(text) {
		offset := NewLineOffset{
			LineStartPos: currentNewLineOffset,
			LineEndPos:   -1,
			LineNumber:   currentLine,
		}
		result[pos] = offset
		if text[pos] == '\n' {
			currentNewLineOffset = pos
			currentLine = currentLine + 1
		}
	}

	currentNewLineOffset = len(text)

	for pos := len(text) - 1; pos >= 0; pos-- {
		if text[pos] == '\n' {
			currentNewLineOffset = pos
		}
		result[pos].LineEndPos = currentNewLineOffset
	}

	return result
}

type TextReader interface {
	ParseText(QuestionSet) ([]TextItem, error)
}

type Backend interface {
	ScanText(textItems []TextItem) <-chan ProbsRange
}

type Frontend interface {
	DisplayOutput(textItems []TextItem) error
}

type AsyncFrontend interface {
	DisplayOutput(probsRangeChan <-chan ProbsRange)
}
