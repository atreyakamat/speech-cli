//go:build darwin

package hotkey

import (
	"log"
)

// Listen triggers callback when Alt+S is pressed.
// macOS support is currently limited to manual triggering via CLI.
func Listen(callback func()) error {
	log.Printf("[hotkey] global hotkeys are not yet supported on macOS. Please use 'speechd run' or 'speech-cli' for manual interaction.")
	return nil
}
