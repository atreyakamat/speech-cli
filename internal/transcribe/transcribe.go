package transcribe

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/atreya/speech-cli/internal/util"
)

var tsPrefix = regexp.MustCompile(`^\[[0-9:.\s\-\>]+\]\s*`)

func Run(ctx context.Context, commandTpl, wavPath, modelPath string) (string, error) {
	cmdline := util.ExpandPlaceholders(commandTpl, map[string]string{
		"wav":   wavPath,
		"model": modelPath,
	})
	sh := exec.CommandContext(ctx, "/bin/sh", "-c", cmdline)
	var stdout, stderr bytes.Buffer
	sh.Stdout = &stdout
	sh.Stderr = &stderr
	if err := sh.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("transcribe failed: %s", msg)
	}

	// whisper.cpp often prints timestamps per line; strip and join.
	lines := strings.Split(stdout.String(), "\n")
	var kept []string
	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" {
			continue
		}
		ln = tsPrefix.ReplaceAllString(ln, "")
		kept = append(kept, ln)
	}
	return util.CleanTranscript(strings.Join(kept, " ")), nil
}
