# 📋 Project TODOs & Future Improvements

This document outlines planned improvements, potential features, and architectural enhancements for `speech-cli`.

## 🚀 High Priority (Short Term)
- [ ] **Native Go Audio Capture**: Move away from external dependencies like `arecord` and `ffmpeg` for recording. Implement native audio capture using `portaudio` or `oto` for better reliability.
- [ ] **Integrated STT Engine**: Provide an optional built-in STT engine (using `whisper.go` or `sherpa-onnx` bindings) so users don't have to manage external `whisper-cli` binaries.
- [ ] **Improved Installation Scripts**: Combine Linux and Windows installation logic into a more robust Go-based installer that can detect and install missing dependencies.
- [ ] **Cross-Platform Service Integration**: 
    - [ ] Linux: Better `systemd` user unit management.
    - [ ] Windows: Proper Windows Service integration using `winsvc`.
- [ ] **Configuration UI**: A simple TUI or web-based configuration dashboard for managing hotkeys, models, and recording settings.

## 🛠️ Feature Enhancements
- [ ] **Support for More STT Engines**: 
    - [ ] Add support for `faster-whisper`.
    - [ ] Add support for cloud-based providers (OpenAI, Anthropic, Google) as fallback options.
- [ ] **Context-Aware Transcription**: Allow users to provide a custom dictionary or "context" to improve transcription accuracy for technical terms or names.
- [ ] **Advanced Hotkey Support**: Support for custom key combinations and global hotkey management across all desktop environments.
- [ ] **Improved Overlay Bar**: 
    - [ ] Create a more polished, cross-platform floating status bar (maybe using `Fyne` or `Wails`).
    - [ ] Show real-time volume levels during recording.
- [ ] **Auto-Update**: Implement a mechanism for `speech-cli` to check for and apply updates from GitHub Releases.

## 🏗️ Architectural & DX Improvements
- [ ] **Unit and Integration Tests**: Increase test coverage, especially for the audio recording and transcription logic.
- [ ] **Better Logging & Error Handling**: Implement structured logging (e.g., using `slog`) and provide more actionable error messages in `speech-cli doctor`.
- [ ] **CI/CD Enhancements**: 
    - [ ] Automated linting and security scanning in GitHub Actions.
    - [ ] Automated testing on multiple OS environments in CI.
- [ ] **Refactor Daemon Logic**: Decouple the transcription engine from the daemon to allow for easier testing and swapping of STT backends.

## 📝 Documentation
- [ ] Add more platform-specific troubleshooting guides.
- [ ] Create a dedicated documentation site (e.g., using MkDocs or Docusaurus).
- [ ] Provide a "Common Issues" section in the README.
