package record

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
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
	cmd := exec.CommandContext(ctx, r.Cmd,
		"-q",
		"-f", "S16_LE",
		"-r", fmt.Sprintf("%d", r.SampleHz),
		"-c", fmt.Sprintf("%d", r.Channels),
		"-t", "wav",
		wav,
	)
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

	// Ask arecord to gracefully finish the wav header.
	if cmd.Process != nil {
		_ = cmd.Process.Signal(os.Interrupt)
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
