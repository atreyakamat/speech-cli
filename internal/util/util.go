package util

import (
	"errors"
	"os/exec"
	"strings"
)

func LookPath(cmd string) (string, bool) {
	p, err := exec.LookPath(cmd)
	return p, err == nil
}

func ExpandPlaceholders(tpl string, repl map[string]string) string {
	out := tpl
	for k, v := range repl {
		out = strings.ReplaceAll(out, "{"+k+"}", v)
	}
	return out
}

func CleanTranscript(s string) string {
	// Minimal normalization: trim and collapse whitespace.
	fields := strings.Fields(strings.TrimSpace(s))
	return strings.Join(fields, " ")
}

var ErrNotConfigured = errors.New("not configured")
