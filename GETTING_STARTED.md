# Speech-CLI - Getting Started Guide

## 🎤 What is Speech-CLI?

Speech-CLI is a **system-wide voice-to-text daemon** that lets you dictate text into any application using a simple hotkey. Press and hold **Alt+S**, speak, release, and your words are automatically transcribed and typed into the focused window.

**Key Features:**
- 🔒 **100% Offline** - No internet required after initial setup
- 🚀 **Fast & Lightweight** - Uses local Whisper models
- 🖥️ **Cross-Platform** - Linux (fully supported), Windows, macOS
- ⌨️ **Universal** - Works with any text input field
- 🎯 **Privacy-First** - All processing happens locally on your machine

---

## 📋 Quick Start (Linux)

### 1. Install Dependencies

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install -y alsa-utils ffmpeg wtype xdotool zenity git build-essential
```

**Arch Linux:**
```bash
sudo pacman -S --needed alsa-utils ffmpeg wtype xdotool zenity git base-devel
```

**Fedora:**
```bash
sudo dnf install -y alsa-utils ffmpeg wtype xdotool zenity git gcc make
```

### 2. Install Whisper.cpp

```bash
# Clone and build whisper.cpp
cd ~
git clone https://github.com/ggerganov/whisper.cpp
cd whisper.cpp
make

# Make the binary available in PATH
sudo cp main /usr/local/bin/whisper-cli
```

### 3. One-Command Setup

From the speech-cli directory:

```bash
bash scripts/install-linux.sh
```

This script will:
- Build speech-cli binaries
- Initialize configuration
- Download a Whisper model (tiny.en by default)
- Verify dependencies
- Set up the systemd service

### 4. Start Using It!

**Option A: Run in foreground (for testing)**
```bash
sudo ./speechd
```

**Option B: Install as a system service**
```bash
./speech-cli service install
./speech-cli service start
./speech-cli service status
```

**Usage:**
1. Focus any text input field
2. Press and hold **Alt+S**
3. Speak clearly
4. Release **Alt+S**
5. Your transcribed text appears!

---

## 🪟 Windows Setup

### 1. Install Prerequisites

1. **Install Go** (1.26 or later): https://go.dev/dl/
2. **Install FFmpeg**: https://ffmpeg.org/download.html
   - Add to PATH environment variable
3. **Install Whisper.cpp**:
   ```powershell
   # Using prebuilt binary or compile from source
   # Ensure whisper-cli.exe is in PATH
   ```

### 2. One-Command Setup

Open PowerShell as Administrator:

```powershell
cd path\to\speech-cli
powershell -ExecutionPolicy Bypass -File .\scripts\install-windows.ps1
```

### 3. Manual Setup (if script doesn't work)

```powershell
# Build binaries
go build -o speech-cli.exe ./cmd/speech-cli
go build -o speechd.exe ./cmd/speechd

# Initialize config
.\speech-cli.exe init

# Download model
.\speech-cli.exe model set tiny.en

# Verify setup
.\speech-cli.exe doctor

# Run daemon
.\speechd.exe
```

### 4. Configure FFmpeg Audio Input (if needed)

If your microphone isn't detected automatically:

```powershell
# List audio devices
ffmpeg -list_devices true -f dshow -i dummy

# Set the device name
$env:SPEECH_FFMPEG_INPUT = "audio=Microphone (Your Device Name)"
```

---

## 🍎 macOS Setup

**Current Status:** Builds successfully, but text injection is limited (stdout only). Full hotkey + typing support is in development.

### What Works Now

```bash
# Build and test
go build ./...
go run ./cmd/speech-cli init
go run ./cmd/speech-cli doctor
go run ./cmd/speech-cli model set tiny.en
```

---

## ⚙️ Configuration

### View Current Configuration

```bash
./speech-cli config get
```

### Common Settings

**Change Whisper Model:**
```bash
# Download a different model
./speech-cli model set base.en

# Or import a local model
./speech-cli model import /path/to/ggml-model.bin

# List available models
./speech-cli model list
```

**Toggle Recording Overlay (Linux only):**
```bash
./speech-cli config set overlay-bar true
./speech-cli config set overlay-bar false
```

**Change Audio Recorder (Linux):**
```bash
# Use ffmpeg instead of arecord
./speech-cli config set recorder ffmpeg
```

**Custom Whisper Command:**
```bash
./speech-cli config set whisper-cmd 'whisper-cli -m {model} -f {wav} --no-timestamps -l en'
```

**View Config File:**
```bash
cat ~/.speech-cli/config.json
```

---

## 🔍 Troubleshooting

### Run Diagnostics

```bash
./speech-cli doctor
```

This checks:
- ✓ Configuration file exists
- ✓ Whisper binary is available
- ✓ Model file exists
- ✓ Audio recorder is installed
- ✓ Text injection tool is available
- ✓ Overlay tools are installed (Linux)

### Common Issues

**"missing: whisper-cli"**
- Install whisper.cpp and ensure `whisper-cli` is in PATH
- Or set `SPEECH_WHISPER_BIN` environment variable

**Hotkey Not Working (Linux)**
- Run with `sudo` (requires `/dev/input` access)
- Or add your user to the `input` group:
  ```bash
  sudo usermod -aG input $USER
  # Log out and back in
  ```

**No Text Appears**
- Linux Wayland: Install `wtype`
- Linux X11: Install `xdotool`
- Windows: Ensure PowerShell is available
- Check that the target window can receive keyboard input

**Empty Transcriptions**
- Test your microphone: `arecord -d 5 test.wav` (Linux)
- Verify model path: `./speech-cli config get model`
- Check Whisper works: `whisper-cli -m model.bin -f test.wav`

**Permission Denied (Linux)**
- Hotkey capture requires root or input group membership
- Service installation requires sudo

---

## 📚 Advanced Usage

### Environment Variables

```bash
# Custom Whisper binary location
export SPEECH_WHISPER_BIN=/custom/path/whisper

# Custom model location
export SPEECH_WHISPER_MODEL=/custom/path/model.bin

# Alternative STT engine
export SPEECH_STT_BIN=/path/to/custom-stt

# FFmpeg input device (Windows)
export SPEECH_FFMPEG_INPUT="audio=Microphone Name"
```

### Service Management (Linux)

```bash
# Install service
./speech-cli service install

# Start/stop/restart
./speech-cli service start
./speech-cli service stop
./speech-cli service restart

# Check status
./speech-cli service status

# View logs
./speech-cli service logs
./speech-cli service logs -f  # Follow mode
```

### Different Hotkeys

Currently hardcoded to **Alt+S**. To change, edit the source code in `internal/hotkey/`.

### Using Different Models

Speech-CLI supports any Whisper GGML model:

**Available models (smallest to largest):**
- `tiny.en` - Fastest, ~75 MB, English only
- `base.en` - Good balance, ~142 MB, English only  
- `small.en` - Better accuracy, ~466 MB, English only
- `tiny` - Multilingual, ~75 MB
- `base` - Multilingual, ~142 MB
- `small` - Multilingual, ~466 MB

```bash
# Download and set
./speech-cli model set small.en

# Manual download from Hugging Face
wget https://huggingface.co/ggerganov/whisper.ggml/resolve/main/ggml-small.en.bin
./speech-cli model import ggml-small.en.bin
```

---

## 🏗️ Building from Source

### Requirements
- Go 1.26 or later
- Git

### Build Commands

```bash
# Clone repository
git clone <repository-url>
cd speech-cli

# Build all binaries
go build ./...

# Build specific binaries
go build -o speech-cli ./cmd/speech-cli
go build -o speechd ./cmd/speechd

# Run tests
go test ./...

# Cross-compile for other platforms
GOOS=windows GOARCH=amd64 go build -o speechd.exe ./cmd/speechd
GOOS=darwin GOARCH=arm64 go build -o speechd ./cmd/speechd
```

---

## 🎯 Use Cases

- **Coding faster** - Dictate comments, documentation, commit messages
- **Writing emails** - Quick message composition
- **Accessibility** - Hands-free text input
- **Note-taking** - Capture thoughts quickly
- **Chat applications** - Voice-to-text in messaging apps
- **Form filling** - Speed up data entry

---

## 🔐 Privacy & Offline Operation

✅ **Completely offline after setup:**
1. Whisper runs locally (no API calls)
2. Audio never leaves your machine
3. No telemetry or tracking
4. Open source - audit the code yourself

**Internet required only for:**
- Initial model download (can be done manually)
- Whisper.cpp installation (can be done offline)

**After setup, disconnect from internet and everything still works!**

---

## 📖 Project Structure

```
speech-cli/
├── cmd/
│   ├── speech-cli/    # CLI tool for config/model management
│   └── speechd/       # Main daemon process
├── internal/
│   ├── audio/         # Audio recording abstraction
│   ├── config/        # Configuration management
│   ├── daemon/        # Main daemon loop
│   ├── hotkey/        # Platform-specific hotkey capture
│   ├── inject/        # Text injection backends
│   ├── overlay/       # Recording indicator UI
│   ├── record/        # Audio recording implementations
│   └── stt/           # Speech-to-text interface
├── models/            # Whisper model files
├── docs/              # Documentation
└── scripts/           # Installation scripts
```

---

## 🤝 Contributing

Contributions welcome! Areas for improvement:
- macOS native text injection support
- Additional hotkey options
- More STT engine integrations
- Performance optimizations
- Documentation improvements

---

## 📄 License

See LICENSE file in repository.

---

## 🆘 Support

- **Report issues**: Use GitHub Issues
- **Documentation**: See `docs/USAGE_GUIDE.md`
- **Release notes**: See `CHANGELOG.md`

---

## 🎊 Quick Reference Card

```
┌─────────────────────────────────────────────────────┐
│  SPEECH-CLI - Quick Reference                       │
├─────────────────────────────────────────────────────┤
│  HOTKEY:   Alt+S (hold to record, release to type)  │
│                                                      │
│  SETUP:    bash scripts/install-linux.sh            │
│                                                      │
│  RUN:      sudo ./speechd                           │
│                                                      │
│  SERVICE:  ./speech-cli service install             │
│            ./speech-cli service start               │
│                                                      │
│  CONFIG:   ./speech-cli config get                  │
│            ./speech-cli config set <key> <value>    │
│                                                      │
│  MODELS:   ./speech-cli model list                  │
│            ./speech-cli model set tiny.en           │
│                                                      │
│  HELP:     ./speech-cli doctor                      │
│            ./speech-cli --help                      │
└─────────────────────────────────────────────────────┘
```

**That's it! Happy dictating! 🎤✨**
