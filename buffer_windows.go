package main

var (
	pCreateCompatibleDC     = gdi32.NewProc("CreateCompatibleDC")
	pCreateCompatibleBitmap = gdi32.NewProc("CreateCompatibleBitmap")
	pDeleteDC               = gdi32.NewProc("DeleteDC")
	pBitBlt                 = gdi32.NewProc("BitBlt")

	backDC     uintptr
	backBitmap uintptr
	backOld    uintptr
	backW      int32
	backH      int32
)

func ensureBackBuffer(screenDC uintptr, w, h int32) uintptr {
	if w <= 0 || h <= 0 {
		return screenDC
	}
	if backDC != 0 && backW == w && backH == h {
		return backDC
	}
	cleanupBackBuffer()
	backDC, _, _ = pCreateCompatibleDC.Call(screenDC)
	backBitmap, _, _ = pCreateCompatibleBitmap.Call(screenDC, uintptr(w), uintptr(h))
	backOld, _, _ = pSelectObject.Call(backDC, backBitmap)
	backW, backH = w, h
	return backDC
}

func presentBackBuffer(screenDC uintptr, w, h int32) {
	if backDC != 0 {
		pBitBlt.Call(screenDC, 0, 0, uintptr(w), uintptr(h), backDC, 0, 0, SRCCOPY)
	}
}

func cleanupBackBuffer() {
	if backDC == 0 {
		return
	}
	if backOld != 0 {
		pSelectObject.Call(backDC, backOld)
	}
	if backBitmap != 0 {
		pDeleteObject.Call(backBitmap)
	}
	pDeleteDC.Call(backDC)
	backDC, backBitmap, backOld = 0, 0, 0
	backW, backH = 0, 0
}
