# SNAPWAVE HORROR V9 LITE/STABLE

V9 is the lightweight stability rebuild after the V8.x performance issues.

## Architecture
- no embedded BMP/PNG images
- no long 121-second WAV
- no DIBSection/back-buffer pipeline
- procedural Win32/GDI visuals
- six short synthesized sound clips
- 100 ms event timer (10 Hz)
- static scenes repaint only on state changes
- normal WM_CLOSE -> DestroyWindow -> cleanup

## Password
`snapwave`

## Safety
No autostart, registry changes, file deletion/encryption, network access, global keyboard hooks, or system-volume control.

The working full source is also supplied in the V9 source ZIP produced with the build.
