package stt

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/atreya/speech-cli/internal/util"
)

var tsPrefix = regexp.MustCompile(`^\[[0-9:.\s\-\>]+\]\s*`)

// Transcribe uses a local transcription engine.
// Defaults to whisper.cpp (./bin/whisper) but supports SenseVoice if configured.
func Transcribe(file string) string {
	bin := "./bin/whisper"
	model := "models/ggml-tiny.en.bin"

	// Support SenseVoice if it exists in bin
	if _, ok := util.LookPath("./bin/sensevoice"); ok {
		bin = "./bin/sensevoice"
		// SenseVoice model is usually handled by the script itself
	}

	if v := os.Getenv("SPEECH_WHISPER_BIN"); v != "" {
		bin = v
	}
	if v := os.Getenv("SPEECH_STT_BIN"); v != "" {
		bin = v
	}
	if v := os.Getenv("SPEECH_WHISPER_MODEL"); v != "" {
		model = v
	}

	// For SenseVoice/Sherpa-ONNX wrapper, we just pass the file
	var cmd *exec.Cmd
	if strings.Contains(bin, "sensevoice") {
		cmd = exec.Command(bin, file)
	} else {
		cmd = exec.Command(bin, "-m", model, "-f", file, "-nt")
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr // pipe stderr for debugging
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "[stt] execution error: %v\n", err)
	}

	result := strings.TrimSpace(out.String())
	lines := strings.Split(result, "\n")
	kept := make([]string, 0, len(lines))
	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" || strings.HasPrefix(ln, "DEBUG:") || strings.HasPrefix(ln, "INFO:") {
			continue
		}
		ln = tsPrefix.ReplaceAllString(ln, "")
		kept = append(kept, ln)
	}
	return util.CleanTranscript(strings.Join(kept, " "))
}

func init() {
	// Ensure bin exists
	_ = os.MkdirAll("bin", 0755)
}
