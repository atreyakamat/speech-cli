package stt

import (
	"bytes"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/atreya/speech-cli/internal/util"
)

var tsPrefix = regexp.MustCompile(`^\[[0-9:.\s\-\>]+\]\s*`)

// Transcribe uses a local whisper.cpp binary (./bin/whisper) and model (./models/...).
// This matches the starter scope.
func Transcribe(file string) string {
	bin := "./bin/whisper"
	model := "models/ggml-base.en.bin"
	if v := os.Getenv("SPEECH_WHISPER_BIN"); v != "" {
		bin = v
	}
	if v := os.Getenv("SPEECH_WHISPER_MODEL"); v != "" {
		model = v
	}

	cmd := exec.Command(bin, "-m", model, "-f", file, "-nt")
	var out bytes.Buffer
	cmd.Stdout = &out
	_ = cmd.Run()

	result := strings.TrimSpace(out.String())
	lines := strings.Split(result, "\n")
	kept := make([]string, 0, len(lines))
	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" {
			continue
		}
		ln = tsPrefix.ReplaceAllString(ln, "")
		kept = append(kept, ln)
	}
	return util.CleanTranscript(strings.Join(kept, " "))
}
