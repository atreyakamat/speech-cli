# Changelog

All notable changes to this project will be documented in this file.

## [0.1.0] - 2026-03-22
- Added semantic version command to both binaries:
  - speech-cli version
  - speechd version
- Fixed default speechd behavior to launch daemon mode when run without arguments.
- Fixed systemd service install to target the speechd daemon binary (not speech-cli).
- Added offline model workflow:
  - speech-cli model import <path>
  - speech-cli model list
- Made model download cross-platform by removing curl dependency from model download path.
- Added OS-specific speechd process control implementations so Linux and Windows both compile.
- Added platform-aware default recorder selection (Linux: arecord, Windows: ffmpeg).
- Improved recorder compatibility with both arecord and ffmpeg command modes.
- Added Windows hotkey implementations for both daemon and MVP paths (Alt+S).
- Added Windows text injection backend via PowerShell SendKeys.
- Added explicit Windows message for unsupported systemd service management commands.
