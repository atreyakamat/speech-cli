package config

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := Default()
	if err := cfg.Validate(); err != nil {
		t.Errorf("default config is invalid: %v", err)
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name:    "valid default",
			cfg:     Default(),
			wantErr: false,
		},
		{
			name: "invalid whisper backend",
			cfg: func() Config {
				c := Default()
				c.Whisper.Backend = "invalid"
				return c
			}(),
			wantErr: true,
		},
		{
			name: "empty whisper command",
			cfg: func() Config {
				c := Default()
				c.Whisper.Backend = "command"
				c.Whisper.Command = ""
				return c
			}(),
			wantErr: true,
		},
		{
			name: "invalid sherpa config",
			cfg: func() Config {
				c := Default()
				c.Whisper.Backend = "sherpa"
				c.STT.Sherpa.ModelPath = ""
				return c
			}(),
			wantErr: true,
		},
		{
			name: "invalid sample rate",
			cfg: func() Config {
				c := Default()
				c.Audio.SampleRate = 0
				return c
			}(),
			wantErr: true,
		},
		{
			name: "invalid inject backend",
			cfg: func() Config {
				c := Default()
				c.Inject.Backend = "non-existent"
				return c
			}(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cfg.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
