# Speech-CLI - 5-Minute Quick Start

## What You Need (Linux)

```bash
# Install everything needed
sudo apt install -y alsa-utils ffmpeg wtype xdotool zenity git build-essential

# Install Whisper.cpp
cd ~ && git clone https://github.com/ggerganov/whisper.cpp
cd whisper.cpp && make
sudo cp main /usr/local/bin/whisper-cli
```

## Install Speech-CLI

```bash
# From the speech-cli directory:
bash scripts/install-linux.sh

# Or manual setup:
go build -o speech-cli ./cmd/speech-cli
go build -o speechd ./cmd/speechd
./speech-cli init
./speech-cli model set tiny.en
```

## Run It

```bash
# Foreground (for testing):
sudo ./speechd

# Or as a service:
./speech-cli service install
./speech-cli service start
```

## Use It

1. **Focus any text field** (browser, editor, chat app, etc.)
2. **Press and hold Alt+S**
3. **Speak clearly**
4. **Release Alt+S**
5. **Your text appears!**

## Troubleshooting

```bash
# Check everything is working:
./speech-cli doctor

# View configuration:
./speech-cli config get

# View logs (if running as service):
./speech-cli service logs
```

## Common Issues

- **Hotkey not working?** Run with `sudo` or add yourself to `input` group
- **No text appears?** Install `wtype` (Wayland) or `xdotool` (X11)
- **Empty transcription?** Check microphone with `arecord -d 3 test.wav && aplay test.wav`

---

**For full documentation, see [GETTING_STARTED.md](GETTING_STARTED.md)**
