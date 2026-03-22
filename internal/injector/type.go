package injector

import (
	"context"

	"github.com/atreya/speech-cli/internal/inject"
)

func Type(text string) {
	_ = inject.Type(context.Background(), inject.BackendAuto, text)
}
