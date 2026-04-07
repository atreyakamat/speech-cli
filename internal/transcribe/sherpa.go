package transcribe

import (
	"context"
	"fmt"
	"os"

	sherpa "github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx"
)

type SherpaTranscriber struct {
	ModelPath  string
	TokensPath string
	NumThreads int
}

func (s *SherpaTranscriber) Transcribe(ctx context.Context, wavPath string) (string, error) {
	if _, err := os.Stat(s.ModelPath); os.IsNotExist(err) {
		return "", fmt.Errorf("sherpa model not found at %s", s.ModelPath)
	}
	if _, err := os.Stat(s.TokensPath); os.IsNotExist(err) {
		return "", fmt.Errorf("sherpa tokens not found at %s", s.TokensPath)
	}

	config := sherpa.OfflineRecognizerConfig{
		ModelConfig: sherpa.OfflineModelConfig{
			SenseVoice: sherpa.OfflineSenseVoiceModelConfig{
				Model: s.ModelPath,
			},
			Tokens:     s.TokensPath,
			NumThreads: s.NumThreads,
			Debug:      0,
		},
	}

	recognizer := sherpa.NewOfflineRecognizer(&config)
	if recognizer == nil {
		return "", fmt.Errorf("failed to create sherpa recognizer")
	}
	defer sherpa.DeleteOfflineRecognizer(recognizer)

	wave := sherpa.ReadWave(wavPath)
	if wave == nil {
		return "", fmt.Errorf("failed to read wave file: %s", wavPath)
	}

	stream := sherpa.NewOfflineStream(recognizer)
	defer sherpa.DeleteOfflineStream(stream)

	stream.AcceptWaveform(wave.SampleRate, wave.Samples)

	recognizer.Decode(stream)
	result := stream.GetResult()

	return result.Text, nil
}
