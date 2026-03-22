# Speech CLI Usage Guide (Linux, Windows, macOS)

This guide shows how to run and use Speech CLI on your system, including offline setup.

## 1. What is required

### Core requirements (all platforms)
- Go 1.26+
- A local Whisper executable (`whisper-cli`) available in your `PATH`
- A local Whisper model file (`ggml-*.bin`)

### Audio capture requirements
- Linux: `arecord` (default) or `ffmpeg`
- Windows: `ffmpeg` (default)
- macOS: `ffmpeg` (recommended)

### Text injection requirements
- Linux Wayland: `wtype`
- Linux X11: `xdotool`
- Windows: built-in PowerShell backend (no extra package)
- macOS: currently falls back to stdout (text print), not native typing

### Floating recording bar requirements
- Linux: install `yad` or `zenity` for on-screen recording bar
- Windows/macOS: floating bar backend is not implemented yet

## 2. Build and quick verification

From project root:

```bash
go test ./...
go build ./...
go run ./cmd/speech-cli version
go run ./cmd/speechd version
```

Expected version output currently: `0.1.0`

## 2.1 One-command setup scripts

Linux:

```bash
bash scripts/install-linux.sh
```

Windows (PowerShell):

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\install-windows.ps1
```

Common script options:
- Linux: `--model base.en`, `--import-model /path/to/model.bin`, `--skip-model-download`
- Windows: `-Model base.en`, `-ImportModel C:\path\model.bin`, `-SkipModelDownload`

## 3. Initial setup (all platforms)

```bash
go run ./cmd/speech-cli init
```

Set your transcription command (default is already set, but explicit is better):

```bash
go run ./cmd/speech-cli config set whisper-cmd 'whisper-cli -m {model} -f {wav} --no-timestamps -l en'
```

## 4. Model setup

### Option A: Download model (internet required once)

```bash
go run ./cmd/speech-cli model set tiny.en
```

### Option B: Import local model (offline-friendly)

```bash
go run ./cmd/speech-cli model import /path/to/ggml-tiny.en.bin
```

List installed models:

```bash
go run ./cmd/speech-cli model list
```

## 5. Linux usage (recommended production target)

### Install dependencies

Ubuntu/Debian example:

```bash
sudo apt update
sudo apt install -y alsa-utils ffmpeg wtype xdotool zenity
```

Arch example:

```bash
sudo pacman -S --needed alsa-utils ffmpeg wtype xdotool zenity
```

### Run diagnostics

```bash
go run ./cmd/speech-cli doctor
```

Fix any `missing:` lines before continuing.

### Start daemon (foreground)

```bash
go run ./cmd/speechd
```

Hold `Alt+S` to record, release to transcribe and type.

Toggle floating recording bar:

```bash
go run ./cmd/speech-cli config set overlay-bar true
go run ./cmd/speech-cli config set overlay-bar false
```

### Start daemon (background process commands)

```bash
go run ./cmd/speechd start
go run ./cmd/speechd status
go run ./cmd/speechd stop
```

### Systemd user service (Linux only)

```bash
go build -o speech-cli ./cmd/speech-cli
go build -o speechd ./cmd/speechd
./speech-cli service install
./speech-cli service status
```

## 6. Windows usage

### Install dependencies
- Install `ffmpeg` and add it to `PATH`
- Install `whisper-cli` and add it to `PATH`
- Provide a local `ggml-*.bin` model via `model set` or `model import`

If ffmpeg does not detect your default microphone, set input device string:

```powershell
$env:SPEECH_FFMPEG_INPUT = "audio=Microphone (Your Device Name)"
```

### Build and run

```powershell
go build -o speech-cli.exe ./cmd/speech-cli
go build -o speechd.exe ./cmd/speechd
.\speech-cli.exe doctor
.\speechd.exe run
```

Hotkey backend on Windows: `Alt+S` (polled via Win32 key state).
Text injection backend on Windows: PowerShell `SendKeys`.

## 7. macOS usage

### Current support level
- Build support: yes
- Basic config/model/doctor commands: yes
- Global hotkey + native text injection: not fully implemented yet

### What you can run now

```bash
go build ./...
go run ./cmd/speech-cli init
go run ./cmd/speech-cli doctor
go run ./cmd/speech-cli model list
```

For full end-user dictation workflow today, Linux and Windows are the reliable targets.

## 8. Offline operation checklist

After first-time provisioning, runtime can work without internet if all are local:
- `whisper-cli` local binary
- local model file configured in `config.json`
- local recorder and injector dependencies installed

Validate with:

```bash
go run ./cmd/speech-cli doctor
go run ./cmd/speech-cli config get model
```

## 9. Troubleshooting

### `missing: whisper-cli`
Install local whisper executable and ensure it is in `PATH`.

### Linux hotkey not working
You likely need `/dev/input` permissions (or root).

### No typing occurs
Check injector availability:
- Linux Wayland: `wtype`
- Linux X11: `xdotool`
- Windows: PowerShell available and focused target window accepts synthetic keys

### Transcription empty
Verify microphone capture and model path.
