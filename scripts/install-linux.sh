#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

MODEL_NAME="tiny.en"
IMPORT_MODEL=""
SKIP_MODEL_DOWNLOAD="false"

usage() {
  cat <<'EOF'
Usage: scripts/install-linux.sh [options]

Options:
  --model <name>          Model to download with "speech-cli model set" (default: tiny.en)
  --import-model <path>   Import an existing local model file instead of downloading
  --skip-model-download   Skip model setup step
  -h, --help              Show help

Examples:
  scripts/install-linux.sh
  scripts/install-linux.sh --model base.en
  scripts/install-linux.sh --import-model ~/.speech-cli/models/ggml-tiny.en.bin
  scripts/install-linux.sh --skip-model-download
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --model)
      MODEL_NAME="${2:-}"
      shift 2
      ;;
    --import-model)
      IMPORT_MODEL="${2:-}"
      shift 2
      ;;
    --skip-model-download)
      SKIP_MODEL_DOWNLOAD="true"
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown argument: $1" >&2
      usage
      exit 2
      ;;
  esac
done

echo "[1/6] Building binaries"
go build -o speech-cli ./cmd/speech-cli
go build -o speechd ./cmd/speechd

echo "[2/6] Initializing config"
./speech-cli init

WHISPER_BIN="whisper-cli"
if command -v whisper-cli >/dev/null 2>&1; then
  WHISPER_BIN="whisper-cli"
elif [[ -x "$ROOT_DIR/bin/whisper-cli" ]]; then
  WHISPER_BIN="$ROOT_DIR/bin/whisper-cli"
fi

echo "[3/6] Setting whisper command"
./speech-cli config set whisper-cmd "$WHISPER_BIN -m {model} -f {wav} --no-timestamps -l en"

echo "[4/6] Configuring model"
if [[ -n "$IMPORT_MODEL" ]]; then
  ./speech-cli model import "$IMPORT_MODEL"
elif [[ "$SKIP_MODEL_DOWNLOAD" == "true" ]]; then
  echo "Skipped model setup (you can run ./speech-cli model set tiny.en later)."
else
  ./speech-cli model set "$MODEL_NAME"
fi

echo "[5/6] Running diagnostics"
./speech-cli doctor

echo "[6/6] Done"
echo "Run daemon in foreground: ./speechd run"
echo "Run daemon in background: ./speechd start && ./speechd status"
