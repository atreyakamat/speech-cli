//go:build !windows && !linux && !darwin

package hotkey

import (
	"log"
)

func Listen(callback func()) error {
	log.Printf("[hotkey] global hotkeys are not supported on this platform.")
	return nil
}
