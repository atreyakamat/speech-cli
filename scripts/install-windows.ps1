param(
    [string]$Model = "tiny.en",
    [string]$ImportModel = "",
    [switch]$SkipModelDownload
)

$ErrorActionPreference = "Stop"

$RootDir = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
Set-Location $RootDir

Write-Host "[1/6] Building binaries"
go build -o speech-cli.exe ./cmd/speech-cli
go build -o speechd.exe ./cmd/speechd

Write-Host "[2/6] Initializing config"
./speech-cli.exe init

$whisperBin = "whisper-cli"
if (Get-Command whisper-cli -ErrorAction SilentlyContinue) {
    $whisperBin = "whisper-cli"
} elseif (Test-Path ".\bin\whisper-cli.exe") {
    $whisperBin = (Resolve-Path ".\bin\whisper-cli.exe").Path
} elseif (Test-Path ".\bin\whisper-cli") {
    $whisperBin = (Resolve-Path ".\bin\whisper-cli").Path
}

Write-Host "[3/6] Setting whisper command"
./speech-cli.exe config set whisper-cmd "$whisperBin -m {model} -f {wav} --no-timestamps -l en"

Write-Host "[4/6] Configuring model"
if ($ImportModel -ne "") {
    ./speech-cli.exe model import $ImportModel
} elseif ($SkipModelDownload.IsPresent) {
    Write-Host "Skipped model setup (you can run .\speech-cli.exe model set tiny.en later)."
} else {
    ./speech-cli.exe model set $Model
}

Write-Host "[5/6] Running diagnostics"
./speech-cli.exe doctor

Write-Host "[6/6] Done"
Write-Host "Run daemon in foreground: .\speechd.exe run"
Write-Host "Run daemon background helpers: .\speechd.exe start ; .\speechd.exe status"
