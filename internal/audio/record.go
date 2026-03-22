package audio

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/atreya/speech-cli/internal/util"
)

// Record records fixed-duration audio and returns a WAV file path.
// Prefers ffmpeg (as in the starter scope); falls back to arecord.
func Record(seconds int) string {
	file := "/tmp/speech.wav"
	_ = os.Remove(file)

	if _, ok := util.LookPath("ffmpeg"); ok {
		cmd := exec.Command(
			"ffmpeg",
			"-y",
			"-f", "alsa",
			"-i", "default",
			"-t", fmt.Sprintf("%d", seconds),
			file,
		)
		_ = cmd.Run()
		return file
	}

	if _, ok := util.LookPath("arecord"); ok {
		cmd := exec.Command(
			"arecord",
			"-q",
			"-f", "S16_LE",
			"-r", "16000",
			"-c", "1",
			"-t", "wav",
			"-d", fmt.Sprintf("%d", seconds),
			file,
		)
		_ = cmd.Run()
		return file
	}

	return file
}
