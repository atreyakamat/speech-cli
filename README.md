# speech-cli

System-wide speech input layer (Linux-first MVP).

## What you get
### Starter-scope MVP (matches the snippet)
- `speechd` (default): press **Alt+S** → records **5s** → runs **whisper.cpp CLI** → types into focused window
- Whisper binary expected at `./bin/whisper`
- Model expected at `./models/ggml-base.en.bin`

Run (default = hold-to-talk daemon):
```bash
go run ./cmd/speechd
```

Starter fixed-5s mode still exists:
```bash
go run ./cmd/speechd mvp
```

### Daemon mode (default)
- `speechd` / `speechd run`: **push-to-talk hold** (Alt+S) using config-driven whisper command and model path

## Install / build
```bash
go build ./...
```

## Whisper.cpp setup (for starter-scope MVP)
```bash
git clone https://github.com/ggerganov/whisper.cpp
cd whisper.cpp && make
cp main /path/to/speech-cli/bin/whisper
bash ./models/download-ggml-model.sh base.en
cp models/ggml-base.en.bin /path/to/speech-cli/models/
```

## Dependencies
- Hotkey (Linux evdev): may require `sudo` or `/dev/input` permissions
- Audio: `ffmpeg` (preferred) or `arecord`
- Injection: `wtype` (Wayland) or `xdotool` (X11)

## Daemon mode config
```bash
./speech-cli init
./speech-cli config set model ~/.speech-cli/models/ggml-tiny.en.bin
./speech-cli config set whisper-cmd 'whisper-cli -m {model} -f {wav} --no-timestamps -l en'

sudo ./speechd run
```

## Troubleshooting
```bash
./speech-cli doctor
```

## User service (systemd)
```bash
./speech-cli service install
./speech-cli service status
./speech-cli service logs
```
