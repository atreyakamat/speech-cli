//go:build darwin

package evhotkey

import (
	"context"
	"log"
)

// PushToTalk emits start=true when Alt+S is pressed, and start=false when released.
// macOS support is currently limited.
func PushToTalk(ctx context.Context) (<-chan bool, <-chan error) {
	out := make(chan bool)
	errCh := make(chan error)
	
	log.Printf("[evhotkey] PushToTalk is not yet supported on macOS. Use manual CLI tools.")
	
	go func() {
		<-ctx.Done()
		close(out)
		close(errCh)
	}()

	return out, errCh
}
