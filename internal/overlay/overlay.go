package overlay

import (
	"os"
	"os/exec"
	"runtime"
	"sync"

	"github.com/atreya/speech-cli/internal/util"
)

type Indicator struct {
	enabled bool

	mu  sync.Mutex
	cmd *exec.Cmd
}

func New(enabled bool) *Indicator {
	return &Indicator{enabled: enabled}
}

func (i *Indicator) StartRecording() {
	if !i.enabled {
		return
	}

	i.mu.Lock()
	defer i.mu.Unlock()
	if i.cmd != nil {
		return
	}

	cmd := indicatorCommand()
	if cmd == nil {
		return
	}
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return
	}
	i.cmd = cmd
}

func (i *Indicator) Stop() {
	i.mu.Lock()
	cmd := i.cmd
	i.cmd = nil
	i.mu.Unlock()
	if cmd == nil || cmd.Process == nil {
		return
	}
	_ = cmd.Process.Kill()
	_, _ = cmd.Process.Wait()
}

func indicatorCommand() *exec.Cmd {
	if runtime.GOOS != "linux" {
		return nil
	}

	if _, ok := util.LookPath("yad"); ok {
		return exec.Command(
			"yad",
			"--progress",
			"--pulsate",
			"--no-buttons",
			"--undecorated",
			"--skip-taskbar",
			"--on-top",
			"--sticky",
			"--title", "Speech CLI",
			"--text", "Recording...",
			"--width", "320",
		)
	}

	if _, ok := util.LookPath("zenity"); ok {
		return exec.Command(
			"zenity",
			"--progress",
			"--pulsate",
			"--no-cancel",
			"--title", "Speech CLI",
			"--text", "Recording...",
			"--width", "320",
		)
	}

	return nil
}
