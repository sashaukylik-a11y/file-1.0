# SNAPWAVE HORROR V8

Native Windows x64 horror-prank / visual locker simulation written in Go + Win32/GDI.

## V8 story
V8 is a complete story rewrite. There is no swing scene. The horror is built around a recovered CCTV archive from Sector 09: empty corridors, a distant figure that moves closer between missing frames, a fake recording end, a four-camera monitor wall, a recovered frame marked `FACE MATCH: VIEWER`, multiple false endings, then the final SNAPWAVE lock screen.

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

`tools/generate_assets.py` deterministically generates the 2D horror BMP frames and a synchronized ~121 second stereo soundtrack. These generated assets are intentionally gitignored, so all editable source stays small while the final executable contains the real audiovisual resources and is larger than 10 MB.

## Password
`snapwave`

## Safety / behavior
- no autostart
- no registry changes
- no file deletion/encryption
- no network access
- no global keyboard hooks
- `Alt+F4` remains available

The application is a fullscreen/topmost visual prank only.
