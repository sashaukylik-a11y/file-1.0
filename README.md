# SNAPWAVE HORROR V8.2 AUDITED

Native Windows x64 horror-prank / visual locker simulation written in Go + Win32/GDI.

V8.2 is the stability rebuild of the CCTV story. The previous hitch around `Следующий кадр пропущен` was addressed by removing lazy BMP rendering, reducing redundant paints, adding robust back-buffer fallback, and fixing the Win32 cleanup/paint lifecycle.

## Story
Recovered CCTV archive from Sector 09: empty corridors, a distant figure moving closer between missing frames, false recording endings, a four-camera wall, a recovered frame marked `FACE MATCH: VIEWER`, multiple scare waves, then the final SNAPWAVE lock screen.

## Build
Requirements:
- Go 1.22+
- Python 3
- Python packages: `numpy`, `Pillow`

Windows:
```bat
build.bat
```

Linux/macOS cross-build:
```bash
./build.sh
```

Both scripts regenerate the real BMP/WAV resources and run the audit before producing `SNAPWAVE_HORROR_V8_2_AUDITED.exe`.

## Password
`snapwave`

## Behavior
- fullscreen/topmost prank UI
- no autostart
- no registry changes
- no deletion/encryption of files
- no network access
- no global keyboard hooks
- `Alt+F4` remains available

See `AUDIT.md` for the 10-pass audit.
