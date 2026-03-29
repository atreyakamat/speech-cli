package util

import "testing"

func TestExpandPlaceholders(t *testing.T) {
	tests := []struct {
		tpl  string
		repl map[string]string
		want string
	}{
		{"{wav} -m {model}", map[string]string{"wav": "file.wav", "model": "base.bin"}, "file.wav -m base.bin"},
		{"{wav} and {wav}", map[string]string{"wav": "file.wav"}, "file.wav and file.wav"},
		{"no placeholders", map[string]string{"wav": "file.wav"}, "no placeholders"},
	}

	for _, tt := range tests {
		if got := ExpandPlaceholders(tt.tpl, tt.repl); got != tt.want {
			t.Errorf("ExpandPlaceholders(%q) = %q, want %q", tt.tpl, got, tt.want)
		}
	}
}

func TestCleanTranscript(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"  hello   world  ", "hello world"},
		{"\nhello\tworld\r", "hello world"},
		{"", ""},
	}

	for _, tt := range tests {
		if got := CleanTranscript(tt.in); got != tt.want {
			t.Errorf("CleanTranscript(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
