package record

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"
)

type ARecorder struct {
	Cmd      string
	SampleHz int
	Channels int
	StateDir string

	mu   sync.Mutex
	proc *exec.Cmd
	wav  string
}

func (r *ARecorder) Start(ctx context.Context) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.proc != nil {
		return "", fmt.Errorf("already recording")
	}
	if err := os.MkdirAll(r.StateDir, 0o755); err != nil {
		return "", err
	}

	wav := filepath.Join(r.StateDir, fmt.Sprintf("recording-%d.wav", time.Now().UnixNano()))
	args, err := recorderArgs(r.Cmd, r.SampleHz, r.Channels, wav)
	if err != nil {
		return "", err
	}
	cmd := exec.CommandContext(ctx, r.Cmd, args...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return "", err
	}

	r.proc = cmd
	r.wav = wav
	return wav, nil
}

func (r *ARecorder) Stop() (string, error) {
	r.mu.Lock()
	cmd := r.proc
	wav := r.wav
	r.proc = nil
	r.wav = ""
	r.mu.Unlock()
	if cmd == nil {
		return "", fmt.Errorf("not recording")
	}

	// Ask recorder to gracefully finish the wav header when possible.
	if cmd.Process != nil {
		switch recorderKind(r.Cmd) {
		case "arecord":
			_ = cmd.Process.Signal(os.Interrupt)
		default:
			if runtime.GOOS == "windows" {
				_ = cmd.Process.Kill()
			} else {
				_ = cmd.Process.Signal(syscall.SIGTERM)
			}
		}
	}

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	select {
	case err := <-done:
		// arecord often returns non-zero on SIGINT; ignore if file exists.
		if _, statErr := os.Stat(wav); statErr == nil {
			return wav, nil
		}
		return "", err
	case <-time.After(2 * time.Second):
		if cmd.Process != nil {
			_ = cmd.Process.Signal(syscall.SIGTERM)
		}
		select {
		case err := <-done:
			if _, statErr := os.Stat(wav); statErr == nil {
				return wav, nil
			}
			return "", err
		case <-time.After(1 * time.Second):
			if cmd.Process != nil {
				_ = cmd.Process.Kill()
			}
			return "", fmt.Errorf("failed to stop recorder")
		}
	}
}

func recorderKind(cmd string) string {
	return strings.ToLower(filepath.Base(cmd))
}

func recorderArgs(cmd string, sampleHz, channels int, wav string) ([]string, error) {
	switch recorderKind(cmd) {
	case "arecord":
		return []string{
			"-q",
			"-f", "S16_LE",
			"-r", fmt.Sprintf("%d", sampleHz),
			"-c", fmt.Sprintf("%d", channels),
			"-t", "wav",
			wav,
		}, nil
	case "ffmpeg", "ffmpeg.exe":
		format := "alsa"
		input := "default"
		if runtime.GOOS == "windows" {
			format = "dshow"
			input = "audio=default"
		}
		if v := strings.TrimSpace(os.Getenv("SPEECH_FFMPEG_INPUT")); v != "" {
			input = v
		}
		return []string{
			"-y",
			"-hide_banner",
			"-loglevel", "error",
			"-f", format,
			"-i", input,
			"-ac", fmt.Sprintf("%d", channels),
			"-ar", fmt.Sprintf("%d", sampleHz),
			wav,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported recorder: %s (supported: arecord, ffmpeg)", cmd)
	}
}

// WithSignalCancel cancels the context on SIGINT/SIGTERM.
func WithSignalCancel(parent context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		select {
		case <-ctx.Done():
		case <-ch:
			cancel()
		}
		signal.Stop(ch)
	}()
	return ctx, cancel
}
