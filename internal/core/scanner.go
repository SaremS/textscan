package core

type Scanner struct {
	reader   TextReader
	backend  ScannerBackend
	frontend OutputFrontend
}

func NewScanner(
	reader TextReader,
	backend ScannerBackend,
	frontend OutputFrontend,
) *Scanner {
	return &Scanner{
		reader:   reader,
		backend:  backend,
		frontend: frontend,
	}
}

func (s *Scanner) Scan() error {
	documents, err := s.reader.ReadText()
	if err != nil {
		return err
	}

	if err := s.backend.ScanText(documents); err != nil {
		return err
	}

	return s.frontend.DisplayOutput(documents)
}
