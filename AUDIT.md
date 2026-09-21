# SNAPWAVE V8.2 audit

V8.2 passed 10 audit passes: formatting/source hygiene, timeline coverage, critical transition tests, scare table integrity, Windows go vet, Win64 GUI build, asset integrity, GDI lifetime/back-buffer fallback, lifecycle/input cleanup, and static safety/reproducibility checks.

Key fixes: BMPs are predecoded into cached DIBSections before the story, static screens do not repaint continuously, back-buffer failure falls back to direct rendering, BeginPaint/EndPaint is paired, WM_CLOSE destroys the window normally, and the password screen receives explicit focus.
