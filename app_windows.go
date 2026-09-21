package main

import (
	"math"
	"syscall"
	"time"
	"unsafe"
)

var (
	hwnd          uintptr
	started       time.Time
	mode          int // 0 story, 1 lock, 2 death, 3 success
	password      string
	successSince  time.Time
	lastVisualKey int64 = -1
	lastDynamic         = time.Time{}
	lastBlink     int64 = -1
	closing       bool
)

const correctPassword = "snapwave"

func drawLock(hdc uintptr, w, h int32) {
	fill(hdc, RECT{0, 0, w, h}, rgb(4, 4, 5))
	fill(hdc, RECT{0, 0, 5, h}, rgb(145, 8, 16))
	fill(hdc, RECT{28, 72, w - 28, 73}, rgb(28, 28, 31))
	fill(hdc, RECT{28, h - 54, w - 28, h - 53}, rgb(20, 20, 23))
	drawText(hdc, "SYSTEM STATUS // COMPROMISED", RECT{28, 22, w - 28, 58}, 14, rgb(110, 110, 112), 500, "Consolas", DT_LEFT|DT_VCENTER|DT_SINGLELINE)
	drawText(hdc, "EXIT: ALT+F4", RECT{28, 22, w - 28, 58}, 12, rgb(70, 70, 72), 400, "Consolas", DT_RIGHT|DT_VCENTER|DT_SINGLELINE)

	pw := int32(math.Min(float64(w)*.66, 900))
	ph := int32(math.Min(float64(h)*.58, 540))
	if pw < 640 { pw = int32(math.Min(float64(w)*.88, 640)) }
	if ph < 420 { ph = int32(math.Min(float64(h)*.78, 420)) }
	px := (w - pw) / 2
	py := (h - ph) / 2
	fill(hdc, RECT{px, py, px + pw, py + ph}, rgb(8, 8, 10))
	fill(hdc, RECT{px, py, px + pw, py + 2}, rgb(78, 9, 14))
	drawText(hdc, "SESSION 09 // SNAPWAVE", RECT{px + 34, py + 24, px + pw - 34, py + 54}, 11, rgb(83, 83, 86), 600, "Consolas", DT_LEFT|DT_VCENTER|DT_SINGLELINE)
	drawText(hdc, "ВАС ПОЙМАЛ @SNAPWAVE", RECT{px + 40, py + 68, px + pw - 40, py + 155}, 48, rgb(237, 237, 239), 800, "Segoe UI", DT_CENTER|DT_VCENTER|DT_SINGLELINE)
	drawText(hdc, "ACCESS TO THIS SESSION HAS BEEN RESTRICTED", RECT{px + 40, py + 150, px + pw - 40, py + 195}, 13, rgb(107, 107, 110), 500, "Consolas", DT_CENTER|DT_VCENTER|DT_SINGLELINE)
	drawText(hdc, "PASSWORD: snapwave", RECT{px + 40, py + 214, px + pw - 40, py + 246}, 14, rgb(150, 150, 152), 600, "Consolas", DT_CENTER|DT_VCENTER|DT_SINGLELINE)

	ix := px + int32(float64(pw)*.11)
	iw := int32(float64(pw) * .60)
	iy := py + 286
	ih := int32(60)
	fill(hdc, RECT{ix, iy, ix + iw, iy + ih}, rgb(13, 13, 16))
	fill(hdc, RECT{ix, iy + ih - 2, ix + iw, iy + ih}, rgb(102, 9, 15))
	shown := ""
	for range password { shown += "•" }
	if int(time.Since(started).Milliseconds()/500)%2 == 0 { shown += "|" }
	drawText(hdc, shown, RECT{ix + 18, iy, ix + iw - 18, iy + ih}, 24, rgb(235, 235, 238), 500, "Consolas", DT_LEFT|DT_VCENTER|DT_SINGLELINE)

	bx := ix + iw + 14
	bw := pw - (bx - px) - int32(float64(pw)*.11)
	fill(hdc, RECT{bx, iy, bx + bw, iy + ih}, rgb(117, 7, 15))
	drawText(hdc, "UNLOCK", RECT{bx, iy, bx + bw, iy + ih}, 14, rgb(245, 245, 245), 700, "Consolas", DT_CENTER|DT_VCENTER|DT_SINGLELINE)
	drawText(hdc, "LOCAL VISUAL SIMULATION // NO SYSTEM CHANGES", RECT{px + 34, py + ph - 54, px + pw - 34, py + ph - 20}, 10, rgb(66, 66, 69), 500, "Consolas", DT_CENTER|DT_VCENTER|DT_SINGLELINE)
}

func drawDeath(hdc uintptr, w, h int32) {
	fill(hdc, RECT{0, 0, w, h}, rgb(24, 0, 2))
	fill(hdc, RECT{0, 0, w, 5}, rgb(155, 6, 16))
	drawText(hdc, "FATAL ERROR", RECT{0, int32(float64(h) * .31), w, int32(float64(h) * .49)}, 90, rgb(255, 33, 42), 900, "Segoe UI", DT_CENTER|DT_VCENTER|DT_SINGLELINE)
	drawText(hdc, "YOU SHOULD HAVE LISTENED", RECT{0, int32(float64(h) * .50), w, int32(float64(h) * .59)}, 20, rgb(185, 79, 85), 600, "Consolas", DT_CENTER|DT_VCENTER|DT_SINGLELINE)
	drawText(hdc, "ALT+F4 // EXIT", RECT{0, int32(float64(h) * .72), w, int32(float64(h) * .79)}, 11, rgb(94, 31, 35), 500, "Consolas", DT_CENTER|DT_VCENTER|DT_SINGLELINE)
}

func drawSuccess(hdc uintptr, w, h int32) {
	fill(hdc, RECT{0, 0, w, h}, rgb(2, 2, 2))
	drawText(hdc, "AUTHENTICATION ACCEPTED", RECT{0, int32(float64(h) * .39), w, int32(float64(h) * .48)}, 16, rgb(105, 105, 108), 500, "Consolas", DT_CENTER|DT_VCENTER|DT_SINGLELINE)
	drawText(hdc, "SNAPWAVE RELEASED THIS SESSION", RECT{0, int32(float64(h) * .48), w, int32(float64(h) * .62)}, 44, rgb(230, 230, 232), 700, "Segoe UI", DT_CENTER|DT_VCENTER|DT_SINGLELINE)
}

func paintStory(hdc uintptr, w, h int32, t float64) {
	if sc, ok := scareAt(t); ok {
		fill(hdc, RECT{0, 0, w, h}, rgb(0, 0, 0))
		if !drawScare(hdc, sc.face, 0, 0, w, h) {
			drawText(hdc, "SIGNAL LOST", RECT{0, 0, w, h}, 54, rgb(176, 8, 16), 900, "Consolas", DT_CENTER|DT_VCENTER|DT_SINGLELINE)
		}
		return
	}
	switch sceneAt(t) {
	case sceneStory:
		drawVignette(hdc, w, h)
		a, b := storyText(t)
		drawText(hdc, a, RECT{int32(float64(w) * .1), int32(float64(h) * .36), int32(float64(w) * .9), int32(float64(h) * .5)}, 38, rgb(226, 226, 228), 650, "Segoe UI", DT_CENTER|DT_VCENTER|DT_WORDBREAK)
		drawText(hdc, b, RECT{int32(float64(w) * .15), int32(float64(h) * .51), int32(float64(w) * .85), int32(float64(h) * .64)}, 20, rgb(108, 110, 112), 400, "Segoe UI", DT_CENTER|DT_VCENTER|DT_WORDBREAK)
	case sceneCorridor0: drawCorridor(hdc, w, h, t, 0)
	case sceneCorridor1: drawCorridor(hdc, w, h, t, 1)
	case sceneCorridor2: drawCorridor(hdc, w, h, t, 2)
	case sceneBlackA, sceneBlackB, sceneBlackC, sceneBlackD:
		fill(hdc, RECT{0, 0, w, h}, 0)
	case sceneFalseEnd1: drawFalseEnd(hdc, w, h, t, 1)
	case sceneMonitor: drawMonitorGrid(hdc, w, h, t)
	case sceneRecovered: drawRecoveredFrame(hdc, w, h, t)
	case sceneAfterCameraText:
		fill(hdc, RECT{0, 0, w, h}, 0)
		drawText(hdc, "последний кадр был снят после отключения камер.", RECT{0, int32(float64(h) * .4), w, int32(float64(h) * .56)}, 27, rgb(176, 176, 178), 600, "Segoe UI", DT_CENTER|DT_VCENTER|DT_WORDBREAK)
	case sceneFalseEnd2: drawFalseEnd(hdc, w, h, t, 2)
	case sceneSnapwave:
		fill(hdc, RECT{0, 0, w, h}, 0)
		drawText(hdc, "@SNAPWAVE", RECT{0, int32(float64(h) * .31), w, int32(float64(h) * .59)}, 98, rgb(239, 239, 241), 900, "Segoe UI", DT_CENTER|DT_VCENTER|DT_SINGLELINE)
	case sceneLock: drawLock(hdc, w, h)
	}
}

func paint(hdc uintptr, rc RECT) {
	w := rc.Right - rc.Left
	h := rc.Bottom - rc.Top
	if w <= 0 || h <= 0 { return }
	switch mode {
	case 1: drawLock(hdc, w, h)
	case 2: drawDeath(hdc, w, h)
	case 3: drawSuccess(hdc, w, h)
	default: paintStory(hdc, w, h, time.Since(started).Seconds())
	}
}

func checkPassword() {
	if password == correctPassword {
		mode = 3
		successSince = time.Now()
		stopSoundtrack()
	} else {
		mode = 2
	}
	pInvalidateRect.Call(hwnd, 0, 0)
}

func renderDue(now time.Time) bool {
	if mode == 1 {
		blink := now.Sub(started).Milliseconds() / 500
		if blink != lastBlink {
			lastBlink = blink
			return true
		}
		return false
	}
	if mode == 2 || mode == 3 { return false }
	t := now.Sub(started).Seconds()
	if t >= 118.2 {
		mode = 1
		c, _, _ := pLoadCursorW.Call(0, IDC_ARROW)
		pSetCursor.Call(c)
		pSetFocus.Call(hwnd)
		return true
	}
	key := visualKey(t)
	if _, ok := scareAt(t); ok {
		if key != lastVisualKey {
			lastVisualKey = key
			return true
		}
		return false
	}
	s := sceneAt(t)
	if dynamicScene(s) {
		if lastDynamic.IsZero() || now.Sub(lastDynamic) >= 50*time.Millisecond {
			lastDynamic = now
			lastVisualKey = key
			return true
		}
		return false
	}
	if key != lastVisualKey {
		lastVisualKey = key
		return true
	}
	return false
}

func handlePaint(h uintptr) {
	var ps PAINTSTRUCT
	dc, _, _ := pBeginPaint.Call(h, uintptr(unsafe.Pointer(&ps)))
	if dc == 0 { return }
	defer pEndPaint.Call(h, uintptr(unsafe.Pointer(&ps)))
	defer func() {
		if recover() != nil {
			var rc RECT
			pGetClientRect.Call(h, uintptr(unsafe.Pointer(&rc)))
			fill(dc, RECT{0, 0, rc.Right, rc.Bottom}, rgb(2, 2, 2))
		}
	}()
	var rc RECT
	pGetClientRect.Call(h, uintptr(unsafe.Pointer(&rc)))
	w, hh := rc.Right-rc.Left, rc.Bottom-rc.Top
	buf := ensureBackBuffer(dc, w, hh)
	paint(buf, RECT{0, 0, w, hh})
	if buf == backDC && backDC != 0 { presentBackBuffer(dc, w, hh) }
}

func wndProc(h uintptr, msg uint32, wp, lp uintptr) (ret uintptr) {
	switch msg {
	case WM_PAINT:
		handlePaint(h); return 0
	case WM_ERASEBKGND:
		return 1
	case WM_TIMER:
		now := time.Now()
		if mode == 3 && !successSince.IsZero() && now.Sub(successSince) > 2200*time.Millisecond {
			if !closing { closing = true; pDestroyWindow.Call(h) }
			return 0
		}
		if renderDue(now) { pInvalidateRect.Call(h, 0, 0) }
		return 0
	case WM_SETCURSOR:
		if mode == 1 { c, _, _ := pLoadCursorW.Call(0, IDC_ARROW); pSetCursor.Call(c) } else { pSetCursor.Call(0) }
		return 1
	case WM_CHAR:
		if mode == 1 {
			ch := rune(wp)
			switch {
			case ch == 8 && len(password) > 0: password = password[:len(password)-1]
			case ch == 13: checkPassword()
			case ch >= 32 && ch < 127 && len(password) < 32: password += string(ch)
			}
			pInvalidateRect.Call(h, 0, 0)
		}
		return 0
	case WM_LBUTTONDOWN:
		if mode == 1 {
			var rc RECT
			pGetClientRect.Call(h, uintptr(unsafe.Pointer(&rc)))
			w, hh := rc.Right, rc.Bottom
			pw := int32(math.Min(float64(w)*.66, 900))
			ph := int32(math.Min(float64(hh)*.58, 540))
			if pw < 640 { pw = int32(math.Min(float64(w)*.88, 640)) }
			if ph < 420 { ph = int32(math.Min(float64(hh)*.78, 420)) }
			px, py := (w-pw)/2, (hh-ph)/2
			ix := px + int32(float64(pw)*.11)
			iw := int32(float64(pw) * .60)
			iy := py + 286
			bx := ix + iw + 14
			bw := pw - (bx - px) - int32(float64(pw)*.11)
			x := int32(int16(lp & 0xffff))
			y := int32(int16((lp >> 16) & 0xffff))
			if x >= bx && x <= bx+bw && y >= iy && y <= iy+60 { checkPassword() }
		}
		return 0
	case WM_CLOSE:
		if !closing { closing = true; pDestroyWindow.Call(h) }
		return 0
	case WM_DESTROY:
		pKillTimer.Call(h, 1)
		stopSoundtrack()
		cleanupEmbeddedBitmaps()
		cleanupGDI()
		pPostQuitMessage.Call(0)
		return 0
	}
	ret, _, _ = pDefWindowProcW.Call(h, uintptr(msg), wp, lp)
	return ret
}

func main() {
	pSetProcessDPIAware.Call()
	inst, _, _ := kernel32.NewProc("GetModuleHandleW").Call(0)
	if inst == 0 { return }
	cls, _ := syscall.UTF16PtrFromString("SnapwaveV82Window")
	wc := WNDCLASSEX{CbSize: uint32(unsafe.Sizeof(WNDCLASSEX{})), LpfnWndProc: syscall.NewCallback(wndProc), HInstance: inst, LpszClassName: cls}
	if r, _, _ := pRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc))); r == 0 { return }
	sw, _, _ := pGetSystemMetrics.Call(SM_CXSCREEN)
	sh, _, _ := pGetSystemMetrics.Call(SM_CYSCREEN)
	if sw == 0 || sh == 0 { return }
	title, _ := syscall.UTF16PtrFromString("SNAPWAVE")
	hwnd, _, _ = pCreateWindowExW.Call(WS_EX_TOPMOST, uintptr(unsafe.Pointer(cls)), uintptr(unsafe.Pointer(title)), WS_POPUP, 0, 0, sw, sh, 0, 0, inst, 0)
	if hwnd == 0 { return }

	initEmbeddedBitmaps()
	started = time.Now()
	lastVisualKey = -1

	pSetWindowPos.Call(hwnd, HWND_TOPMOST, 0, 0, sw, sh, SWP_SHOWWINDOW)
	pShowWindow.Call(hwnd, SW_SHOW)
	pUpdateWindow.Call(hwnd)
	pSetFocus.Call(hwnd)
	playSoundtrack()
	if timer, _, _ := pSetTimer.Call(hwnd, 1, 50, 0); timer == 0 {
		mode = 1
		pInvalidateRect.Call(hwnd, 0, 0)
	}

	var m MSG
	for {
		r, _, _ := pGetMessageW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
		if int32(r) <= 0 { break }
		pTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
		pDispatchMessageW.Call(uintptr(unsafe.Pointer(&m)))
	}
}
