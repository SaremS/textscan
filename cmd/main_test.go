package cmd

import "testing"

func TestRootCommandUsesFileAndQueryArguments(t *testing.T) {
	command := newRootCmd()

	if command.Use != "textscan <file> <query>" {
		t.Fatalf("Use = %q", command.Use)
	}
	if command.Flags().Lookup("max-depth") != nil {
		t.Fatal("max-depth flag must not be present")
	}
	if command.Flags().Lookup("risk-cutoff") == nil {
		t.Fatal("risk-cutoff flag must be present")
	}
	if err := command.Args(command, []string{"notes.txt"}); err == nil {
		t.Fatal("one argument should be rejected")
	}
	if err := command.Args(command, []string{"notes.txt", "topic"}); err != nil {
		t.Fatalf("two arguments should be accepted: %v", err)
	}
}
