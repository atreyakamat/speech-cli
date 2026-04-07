//go:build linux

package hotkey

import (
	"fmt"
	"strings"

	"github.com/MarinX/keylogger"
)

// Listen triggers callback when Alt+S is pressed.
// Linux-only: reads from /dev/input/event* (may require root/permissions).
func Listen(callback func()) error {
	keyboard := keylogger.FindKeyboardDevice()
	if keyboard == "" {
		return fmt.Errorf("no keyboard device found")
	}
	k, err := keylogger.New(keyboard)
	if err != nil {
		return err
	}

	go func() {
		defer k.Close()
		events := k.Read()
		altDown := false
		for e := range events {
			if e.Type != keylogger.EvKey {
				continue
			}
			key := strings.ToUpper(e.KeyString())
			switch key {
			case "L_ALT", "R_ALT":
				if e.KeyPress() {
					altDown = true
				}
				if e.KeyRelease() {
					altDown = false
				}
			case "S":
				if altDown && e.KeyPress() {
					callback()
				}
			}
		}
	}()

	return nil
}
