package core

type DocumentIdentifier string
type Text string

type Document struct {
	Identifier DocumentIdentifier
	Text       Text
	Scoring    *ProbabilityScoring
}

func NewDocument(identifier DocumentIdentifier, text Text) *Document {
	return &Document{
		Identifier: identifier,
		Text:       text,
		Scoring:    newProbabilityScoring(len(text)),
	}
}

type TextReader interface {
	ReadText() ([]Document, error)
}

type ScannerBackend interface {
	ScanText(documents []Document) error
}

type OutputFrontend interface {
	DisplayOutput(documents []Document) error
}
