package transcribe

import (
	"context"
)

type Transcriber interface {
	Transcribe(ctx context.Context, wavPath string) (string, error)
}

type CommandTranscriber struct {
	Command   string
	ModelPath string
}

func (ct *CommandTranscriber) Transcribe(ctx context.Context, wavPath string) (string, error) {
	return Run(ctx, ct.Command, wavPath, ct.ModelPath)
}
