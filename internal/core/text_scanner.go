package core

type TextScanner struct {
	reader   TextReader
	backend  Backend
	frontend AsyncFrontend
}

func NewTextScanner(reader TextReader,
	backend Backend,
	frontend AsyncFrontend) *TextScanner {
	return &TextScanner{
		reader,
		backend,
		frontend,
	}
}

func (s *TextScanner) Scan(questions QuestionSet) {
	textItems, err := s.reader.ParseText(questions)
	if err != nil {
		panic(err)
	}

	probsRangeChan := s.backend.ScanText(textItems)

	s.frontend.DisplayOutput(probsRangeChan)
}
