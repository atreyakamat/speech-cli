//go:build !windows

package evhotkey

import (
	"context"
	"fmt"
	"strings"

	"github.com/MarinX/keylogger"
)

// PushToTalk emits start=true when Alt+S is pressed, and start=false when released.
// Linux-only: reads from /dev/input/event*.
func PushToTalk(ctx context.Context) (<-chan bool, <-chan error) {
	out := make(chan bool, 4)
	errCh := make(chan error, 1)

	keyboard := keylogger.FindKeyboardDevice()
	if keyboard == "" {
		errCh <- fmt.Errorf("no keyboard device found (try setting up /dev/input permissions)")
		close(out)
		return out, errCh
	}

	k, err := keylogger.New(keyboard)
	if err != nil {
		errCh <- err
		close(out)
		return out, errCh
	}

	go func() {
		defer k.Close()
		defer close(out)

		altDown := false
		sDown := false
		recording := false

		events := k.Read()
		for {
			select {
			case <-ctx.Done():
				return
			case e, ok := <-events:
				if !ok {
					return
				}
				if e.Type != keylogger.EvKey {
					continue
				}

				key := strings.ToUpper(e.KeyString())
				if key == "L_ALT" || key == "R_ALT" {
					if e.KeyPress() {
						altDown = true
					}
					if e.KeyRelease() {
						altDown = false
					}
				}
				if key == "S" {
					if e.KeyPress() {
						sDown = true
					}
					if e.KeyRelease() {
						sDown = false
					}
				}

				want := altDown && sDown
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
