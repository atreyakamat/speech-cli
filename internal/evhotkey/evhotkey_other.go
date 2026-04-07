//go:build !windows && !linux && !darwin

package evhotkey

import (
	"context"
	"log"
)

func PushToTalk(ctx context.Context) (<-chan bool, <-chan error) {
	out := make(chan bool)
	errCh := make(chan error)
	log.Printf("[evhotkey] PushToTalk is not supported on this platform.")
	go func() {
		<-ctx.Done()
		close(out)
		close(errCh)
	}()
	return out, errCh
}
