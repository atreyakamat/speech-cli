package daemon

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/atreya/speech-cli/internal/config"
	"github.com/atreya/speech-cli/internal/evhotkey"
	"github.com/atreya/speech-cli/internal/inject"
	"github.com/atreya/speech-cli/internal/overlay"
	"github.com/atreya/speech-cli/internal/record"
	"github.com/atreya/speech-cli/internal/transcribe"
	"github.com/atreya/speech-cli/internal/util"
)

func Run(ctx context.Context, cfg config.Config) error {
	stateDir, err := config.StateDir()
	if err != nil {
		return err
	}

	rec := &record.ARecorder{
		Cmd:      cfg.Audio.Recorder,
		SampleHz: cfg.Audio.SampleRate,
		Channels: cfg.Audio.Channels,
		StateDir: stateDir,
	}
	bar := overlay.New(cfg.UI.ShowRecordingBar)
	defer bar.Stop()

	hot, errCh := evhotkey.PushToTalk(ctx)
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	default:
	}

	var wav string
	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-errCh:
			if err != nil {
				return err
			}
		case start, ok := <-hot:
			if !ok {
				return fmt.Errorf("hotkey channel closed")
			}
			if start {
				log.Printf("[speechd] recording... (Alt+S held)")
				bar.StartRecording()
				w, err := rec.Start(ctx)
				if err != nil {
					bar.Stop()
					log.Printf("[speechd] start record error: %v", err)
					continue
				}
				wav = w
				continue
			}

			// stop
			if wav == "" {
				continue
			}
			bar.Stop()
			log.Printf("[speechd] processing...")
			w, err := rec.Stop()
			if err != nil {
				log.Printf("[speechd] stop record error: %v", err)
				wav = ""
				continue
			}
			wav = ""

			if cfg.Whisper.ModelPath == "" && cfg.Whisper.Backend == "command" {
				log.Printf("[speechd] whisper model_path not set")
				continue
			}

			var tr transcribe.Transcriber
			if cfg.Whisper.Backend == "sherpa" {
				tr = &transcribe.SherpaTranscriber{
					ModelPath:  cfg.STT.Sherpa.ModelPath,
					TokensPath: cfg.STT.Sherpa.TokensPath,
					NumThreads: cfg.STT.Sherpa.NumThreads,
				}
			} else {
				tr = &transcribe.CommandTranscriber{
					Command:   cfg.Whisper.Command,
					ModelPath: cfg.Whisper.ModelPath,
				}
			}

			// Short timeout to avoid hanging forever.
			tctx, cancel := context.WithTimeout(ctx, 60*time.Second)
			text, err := tr.Transcribe(tctx, w)
			cancel()
			if err != nil {
				log.Printf("[speechd] transcribe error: %v", err)
				continue
			}
			text = util.CleanTranscript(text)
			if text == "" {
				log.Printf("[speechd] empty transcript")
				continue
			}

			if err := inject.Type(ctx, inject.Backend(cfg.Inject.Backend), text); err != nil {
				log.Printf("[speechd] inject error: %v", err)
			}
		}
	}
}
