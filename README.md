# speech-cli

System-wide speech input layer (Linux-first MVP).

**🎤 Press Alt+S, speak, release → your words are typed automatically**

## 📖 Documentation
- **New here?** Start with [QUICK_START.md](QUICK_START.md) (5 minutes)
- **Full guide:** [GETTING_STARTED.md](GETTING_STARTED.md) (comprehensive setup)
- **Platform specifics:** [docs/USAGE_GUIDE.md](docs/USAGE_GUIDE.md) (Linux, Windows, macOS)

## One-command setup
- Linux: `bash scripts/install-linux.sh`
- Windows (PowerShell): `powershell -ExecutionPolicy Bypass -File .\\scripts\\install-windows.ps1`

## Tested status
- Build-tested: Linux, Windows, macOS (`amd64` and `arm64` cross-builds)
- Runtime-smoke-tested in this workspace: Linux
- Runtime on Windows/macOS still depends on host-specific audio/input permissions and installed local tools

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

Check version:
```bash
go run ./cmd/speech-cli version
go run ./cmd/speechd version
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

## Model management
Download a model (online):
```bash
./speech-cli model set tiny.en
```

Import a model you already have (offline setup):
```bash
./speech-cli model import /path/to/ggml-base.en.bin
./speech-cli model list
```

After `whisper-cli` and a local model are present, transcription works without internet.

## Linux and Windows notes
- Linux default recorder: `arecord` (you can switch to `ffmpeg` in config)
- Windows default recorder: `ffmpeg`
- For Windows `ffmpeg` capture device naming, set `SPEECH_FFMPEG_INPUT` when needed
- Windows hotkey backend uses Alt+S via Win32 key state polling
- Windows default injection backend uses PowerShell SendKeys
- Both platforms require local `whisper-cli` and a local model file for offline usage
- Floating recording bar on Linux requires `yad` or `zenity`

Enable/disable floating recording bar:
```bash
./speech-cli config set overlay-bar true
./speech-cli config set overlay-bar false
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
