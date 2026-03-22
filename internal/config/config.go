package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	Whisper WhisperConfig `json:"whisper"`
	Inject  InjectConfig  `json:"inject"`
	Audio   AudioConfig   `json:"audio"`
	UI      UIConfig      `json:"ui"`
}

type WhisperConfig struct {
	Command   string `json:"command"` // shell command; supports {wav} and {model}
	ModelPath string `json:"model_path"`
	Language  string `json:"language"`
}

type InjectConfig struct {
	Backend string `json:"backend"` // auto|wtype|xdotool|stdout
}

type AudioConfig struct {
	Recorder   string `json:"recorder"` // arecord
	SampleRate int    `json:"sample_rate"`
	Channels   int    `json:"channels"`
}

type UIConfig struct {
	ShowRecordingBar bool `json:"show_recording_bar"`
}

func Default() Config {
	recorder := "arecord"
	if runtime.GOOS == "windows" {
		recorder = "ffmpeg"
	}

	return Config{
		Whisper: WhisperConfig{
			Command:   "whisper-cli -m {model} -f {wav} --no-timestamps -l en",
			ModelPath: "models/ggml-tiny.en.bin",
			Language:  "en",
		},
		Inject: InjectConfig{Backend: "auto"},
		Audio:  AudioConfig{Recorder: recorder, SampleRate: 16000, Channels: 1},
		UI:     UIConfig{ShowRecordingBar: true},
	}
}

func ConfigPath() (string, error) {
	d, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, "speech-cli", "config.json"), nil
}

func StateDir() (string, error) {
	h, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(h, ".local", "state", "speech-cli"), nil
}

func ModelsDir() (string, error) {
	h, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(h, ".speech-cli", "models"), nil
}

func Load() (Config, error) {
	p, err := ConfigPath()
	if err != nil {
		return Default(), err
	}
	b, err := os.ReadFile(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Default(), nil
		}
		return Default(), err
	}
	cfg := Default()
	if err := json.Unmarshal(b, &cfg); err != nil {
		return Default(), err
	}
	// Env overrides (simple, no extra deps)
	if v := os.Getenv("SPEECH_WHISPER_CMD"); v != "" {
		cfg.Whisper.Command = v
	}
	if v := os.Getenv("SPEECH_MODEL"); v != "" {
		cfg.Whisper.ModelPath = v
	}
	if v := os.Getenv("SPEECH_INJECT_BACKEND"); v != "" {
		cfg.Inject.Backend = v
	}
	return cfg, nil
}

func Save(cfg Config) error {
	p, err := ConfigPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, b, 0o644)
}
