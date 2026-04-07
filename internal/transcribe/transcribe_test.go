package transcribe

import (
	"context"
	"os"
	"testing"
)

func TestTsPrefix(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"[00:00:00.000 -> 00:00:05.000] hello world", "hello world"},
		{"[00:00:05.000] goodbye", "goodbye"},
		{"no timestamp", "no timestamp"},
		{"  [00:00:00.000] leading space", "  [00:00:00.000] leading space"}, // should not match if it's not at the start
	}

	for _, tt := range tests {
		got := tsPrefix.ReplaceAllString(tt.in, "")
		if got != tt.want {
			t.Errorf("tsPrefix mismatch for %q: got %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestRun(t *testing.T) {
	// We can't easily run real whisper-cli in unit tests, 
	// but we can test the command execution logic if we mock it.
	// However, Run is quite simple. Let's test it with a simple echo command.
	
	ctx := context.Background()
	tmp, _ := os.CreateTemp("", "test.wav")
	defer os.Remove(tmp.Name())
	
	// mock whisper-cli output
	tpl := "echo '[00:00:00.000] Result from {wav} with {model}'"
	got, err := Run(ctx, tpl, tmp.Name(), "model.bin")
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	
	want := "Result from " + tmp.Name() + " with model.bin"
	if got != want {
		t.Errorf("Run output = %q, want %q", got, want)
	}
}
