# Release Checklist

## Preflight
- [x] go test ./...
- [x] go build ./...
- [x] GOOS=linux GOARCH=amd64 go build ./...
- [x] GOOS=windows GOARCH=amd64 go build ./...

## Functional smoke
- [x] speech-cli version
- [x] speechd version
- [x] speech-cli init
- [x] speech-cli doctor
- [x] speech-cli model list

## Linux runtime verification
- [x] Verify /dev/input permissions for hotkey capture
- [x] Verify recorder dependency is installed (arecord or ffmpeg)
- [x] Verify injector dependency is installed (wtype or xdotool)
- [x] Verify whisper-cli is installed and available in PATH
- [x] Verify local model file is present and configured

## Windows runtime verification
- [x] Verify ffmpeg is installed and in PATH
- [x] Verify whisper-cli is installed and available in PATH
- [x] Set SPEECH_FFMPEG_INPUT if default audio device string differs
- [x] Verify local model file is present and configured
- [x] Validate hotkey and text injection behavior on target desktop environment

## Offline readiness
- [x] Import a local model with speech-cli model import <path>
- [x] Disable network and run transcription flow with local model + local whisper-cli
- [x] Confirm no network dependency in daemon runtime path

## Tag and release
- [x] Update CHANGELOG.md for final release notes
- [x] Create git tag: v0.1.0
- [x] Build release binaries (Linux + Windows)
- [x] Attach binaries and checksums to release
