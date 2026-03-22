# Test Report - 2026-03-22

## Scope
Validation for release-candidate usage guidance and cross-platform readiness.

## Commands run

### Quality and build
- `go test ./...`
- `go vet ./...`
- `go build ./...`
- `GOOS=linux GOARCH=amd64 go build ./...`
- `GOOS=linux GOARCH=arm64 go build ./...`
- `GOOS=windows GOARCH=amd64 go build ./...`
- `GOOS=windows GOARCH=arm64 go build ./...`
- `GOOS=darwin GOARCH=amd64 go build ./...`
- `GOOS=darwin GOARCH=arm64 go build ./...`

### CLI smoke
- `go run ./cmd/speech-cli version`
- `go run ./cmd/speechd version`
- `go run ./cmd/speech-cli init`
- `go run ./cmd/speech-cli config get model`
- `go run ./cmd/speech-cli model list`
- `go run ./cmd/speech-cli doctor`
- `go run ./cmd/speechd start`
- `go run ./cmd/speechd status`
- `go run ./cmd/speechd stop`
- `timeout 5s go run ./cmd/speechd mvp`

### Installer automation
- `chmod +x scripts/install-linux.sh`
- `./scripts/install-linux.sh --import-model /home/atreya/.speech-cli/models/ggml-tiny.en.bin`
- Windows installer script created: `scripts/install-windows.ps1` (not executable in this Linux environment due to missing PowerShell runtime)

## Observed results
- All tests/build commands succeeded.
- Cross-platform compile matrix succeeded for Linux, Windows, and macOS targets.
- Version commands returned `0.1.0`.
- Model list and config commands succeeded.
- Daemon lifecycle commands succeeded on Linux (`status` reported a running PID).
- MVP mode launches and stays running until interrupted (timeout exit expected).
- `doctor` reported local model configured and showed one missing dependency in this environment:
  - `whisper-cli` missing from `PATH`
- Linux installer script completed successfully through build, init, model import, and doctor stages.

## Conclusion
- Codebase is compile-ready for Linux, Windows, and macOS targets.
- Linux runtime path is smoke-tested.
- Full end-to-end transcription on any platform still requires local `whisper-cli` installation and platform dependencies described in the usage guide.
