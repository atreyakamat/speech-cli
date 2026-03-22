package inject

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/atreya/speech-cli/internal/util"
)

type Backend string

const (
	BackendAuto   Backend = "auto"
	BackendWType  Backend = "wtype"
	BackendXDo    Backend = "xdotool"
	BackendPS     Backend = "powershell"
	BackendStdout Backend = "stdout"
)

func Type(ctx context.Context, backend Backend, text string) error {
	if text == "" {
		return nil
	}
	if backend == BackendAuto {
		if runtime.GOOS == "windows" {
			backend = BackendPS
		}

		if os.Getenv("WAYLAND_DISPLAY") != "" {
			if _, ok := util.LookPath("wtype"); ok {
				backend = BackendWType
			}
		}
		if backend == BackendAuto {
			if _, ok := util.LookPath("xdotool"); ok {
				backend = BackendXDo
			} else if _, ok := util.LookPath("wtype"); ok {
				backend = BackendWType
			} else {
				backend = BackendStdout
			}
		}
	}

	switch backend {
	case BackendWType:
		// wtype --delay 0 -- "text"
		cmd := exec.CommandContext(ctx, "wtype", "--delay", "0", "--", text)
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
		return cmd.Run()
	case BackendXDo:
		cmd := exec.CommandContext(ctx, "xdotool", "type", "--delay", "0", "--clearmodifiers", "--", text)
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
		return cmd.Run()
	case BackendPS:
		if runtime.GOOS != "windows" {
			return fmt.Errorf("powershell backend is only available on windows")
		}
		escaped := strings.ReplaceAll(text, "'", "''")
		script := "Add-Type -AssemblyName System.Windows.Forms; [System.Windows.Forms.SendKeys]::SendWait('" + escaped + "')"
		cmd := exec.CommandContext(ctx, "powershell", "-NoProfile", "-Command", script)
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
		return cmd.Run()
	case BackendStdout:
		fmt.Fprintln(os.Stdout, text)
		return nil
	default:
		return fmt.Errorf("unknown inject backend: %s", backend)
	}
}
