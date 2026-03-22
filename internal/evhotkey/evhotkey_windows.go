//go:build windows

package evhotkey

import (
	"context"
	"syscall"
	"time"
)

const (
	vkMenu = 0x12 // Alt
	vkS    = 0x53 // S
)

var (
	user32            = syscall.NewLazyDLL("user32.dll")
	procGetAsyncState = user32.NewProc("GetAsyncKeyState")
)

func keyDown(vk int) bool {
	ret, _, _ := procGetAsyncState.Call(uintptr(vk))
	return (ret & 0x8000) != 0
}

// PushToTalk emits start=true when Alt+S is pressed, and start=false when released.
func PushToTalk(ctx context.Context) (<-chan bool, <-chan error) {
	out := make(chan bool, 4)
	errCh := make(chan error, 1)

	go func() {
		defer close(out)
		t := time.NewTicker(25 * time.Millisecond)
		defer t.Stop()

		recording := false
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				want := keyDown(vkMenu) && keyDown(vkS)
				if want && !recording {
					recording = true
					out <- true
				}
				if !want && recording {
					recording = false
					out <- false
				}
			}
		}
	}()

	return out, errCh
}
