package config

import (
	"fmt"
	"strings"
)

func (c *Config) Validate() error {
	if c.Whisper.Backend != "command" && c.Whisper.Backend != "sherpa" {
		return fmt.Errorf("invalid whisper backend: %s (must be 'command' or 'sherpa')", c.Whisper.Backend)
	}

	if c.Whisper.Backend == "command" && c.Whisper.Command == "" {
		return fmt.Errorf("whisper command cannot be empty when backend is 'command'")
	}

	if c.Whisper.Backend == "sherpa" {
		if c.STT.Sherpa.ModelPath == "" {
			return fmt.Errorf("sherpa model_path cannot be empty when backend is 'sherpa'")
		}
		if c.STT.Sherpa.TokensPath == "" {
			return fmt.Errorf("sherpa tokens_path cannot be empty when backend is 'sherpa'")
		}
	}

	if c.Audio.SampleRate <= 0 {
		return fmt.Errorf("invalid audio sample_rate: %d", c.Audio.SampleRate)
	}

	if c.Audio.Channels <= 0 {
		return fmt.Errorf("invalid audio channels: %d", c.Audio.Channels)
	}

	validInject := []string{"auto", "wtype", "xdotool", "powershell", "applescript", "stdout"}
	isValidInject := false
	for _, v := range validInject {
		if string(c.Inject.Backend) == v {
			isValidInject = true
			break
		}
	}
	if !isValidInject {
		return fmt.Errorf("invalid inject backend: %s (must be one of: %s)", c.Inject.Backend, strings.Join(validInject, ", "))
	}

	return nil
}
