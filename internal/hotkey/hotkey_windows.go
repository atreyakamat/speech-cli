//go:build windows

package hotkey

import (
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

// Listen triggers callback when Alt+S is pressed.
func Listen(callback func()) error {
	go func() {
		t := time.NewTicker(25 * time.Millisecond)
		defer t.Stop()

		last := false
		for range t.C {
			pressed := keyDown(vkMenu) && keyDown(vkS)
			if pressed && !last {
				callback()
			}
			last = pressed
		}
	}()
	return nil
}
