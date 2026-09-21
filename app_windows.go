package main

import (
	"math"
	"syscall"
	"time"
	"unsafe"
)

var hwnd uintptr
var started time.Time
var mode int
var password string

const correctPassword = "snapwave"

var successSince time.Time

func drawLock(hdc uintptr, w, h int32) {
	fill(hdc, RECT{0, 0, w, h}, rgb(4, 4, 5))
	fill(hdc, RECT{0, 0, 5, h}, rgb(145, 8, 16))
	drawText(hdc, "SYSTEM STATUS // COMPROMISED", RECT{28, 22, w - 28, 58}, 14, rgb(110, 110, 112), 500, "Consolas", DT_LEFT|DT_VCENTER|DT_SINGLELINE)
	drawText(hdc, "EXIT: ALT+F4", RECT{28, 22, w - 28, 58}, 12, rgb(70, 70, 72), 400, "Consolas", DT_RIGHT|DT_VCENTER|DT_SINGLELINE)
	pw := int32(math.Min(float64(w)*.66, 860))
	ph := int32(math.Min(float64(h)*.58, 520))
	px := (w - pw) / 2
	py := (h - ph) / 2
	fill(hdc, RECT{px, py, px + pw, py + ph}, rgb(8, 8, 10))
	drawText(hdc, "ВАС ПОЙМАЛ @SNAPWAVE", RECT{px + 40, py + 52, px + pw - 40, py + 145}, 48, rgb(237, 237, 239), 800, "Segoe UI", DT_CENTER|DT_VCENTER|DT_SINGLELINE)
	drawText(hdc, "ACCESS TO THIS SESSION HAS BEEN RESTRICTED", RECT{px + 40, py + 138, px + pw - 40, py + 185}, 13, rgb(107, 107, 110), 500, "Consolas", DT_CENTER|DT_VCENTER|DT_SINGLELINE)
	drawText(hdc, "PASSWORD: snapwave", RECT{px + 40, py + 198, px + pw - 40, py + 230}, 14, rgb(140, 140, 142), 500, "Consolas", DT_CENTER|DT_VCENTER|DT_SINGLELINE)
	ix := px + int32(float64(pw)*.13)
	iw := int32(float64(pw) * .57)
	iy := py + 270
	ih := int32(58)
	fill(hdc, RECT{ix, iy, ix + iw, iy + ih}, rgb(13, 13, 16))
	shown := ""
	for range password {
		shown += "•"
	}
	if int(time.Since(started).Seconds()*2)%2 == 0 {
		shown += "|"
	}
	drawText(hdc, shown, RECT{ix + 18, iy, ix + iw - 18, iy + ih}, 24, rgb(235, 235, 238), 500, "Consolas", DT_LEFT|DT_VCENTER|DT_SINGLELINE)
	bx := ix + iw + 14
	bw := pw - (bx - px) - int32(float64(pw)*.13)
	fill(hdc, RECT{bx, iy, bx + bw, iy + ih}, rgb(117, 7, 15))
	drawText(hdc, "UNLOCK", RECT{bx, iy, bx + bw, iy + ih}, 14, rgb(245, 245, 245), 700, "Consolas", DT_CENTER|DT_VCENTER|DT_SINGLELINE)
}
func drawDeath(hdc uintptr, w, h int32) {
	fill(hdc, RECT{0, 0, w, h}, rgb(24, 0, 2))
	drawText(hdc, "FATAL ERROR", RECT{0, int32(float64(h) * .31), w, int32(float64(h) * .49)}, 90, rgb(255, 33, 42), 900, "Segoe UI", DT_CENTER|DT_VCENTER|DT_SINGLELINE)
	drawText(hdc, "YOU SHOULD HAVE LISTENED", RECT{0, int32(float64(h) * .50), w, int32(float64(h) * .59)}, 20, rgb(185, 79, 85), 600, "Consolas", DT_CENTER|DT_VCENTER|DT_SINGLELINE)
}
func drawSuccess(hdc uintptr, w, h int32) {
	fill(hdc, RECT{0, 0, w, h}, rgb(2, 2, 2))
	drawText(hdc, "AUTHENTICATION ACCEPTED", RECT{0, int32(float64(h) * .39), w, int32(float64(h) * .48)}, 16, rgb(105, 105, 108), 500, "Consolas", DT_CENTER|DT_VCENTER|DT_SINGLELINE)
	drawText(hdc, "SNAPWAVE RELEASED THIS SESSION", RECT{0, int32(float64(h) * .48), w, int32(float64(h) * .62)}, 44, rgb(230, 230, 232), 700, "Segoe UI", DT_CENTER|DT_VCENTER|DT_SINGLELINE)
}
func paint(hdc uintptr, rc RECT) {
	w := rc.Right
	h := rc.Bottom
	if mode == 1 {
		drawLock(hdc, w, h)
		return
	}
	if mode == 2 {
		drawDeath(hdc, w, h)
		return
	}
	if mode == 3 {
		drawSuccess(hdc, w, h)
		return
	}
	t := time.Since(started).Seconds()
	if sc, ok := scareAt(t); ok {
		fill(hdc, RECT{0, 0, w, h}, rgb(0, 0, 0))
		drawBMP(hdc, scareImages[sc.face], 0, 0, w, h)
		return
	}
	switch {
	case t < 32.4:
		drawVignette(hdc, w, h)
		a, b := storyText(t)
		drawText(hdc, a, RECT{int32(float64(w) * .1), int32(float64(h) * .36), int32(float64(w) * .9), int32(float64(h) * .5)}, 38, rgb(226, 226, 228), 650, "Segoe UI", DT_CENTER|DT_VCENTER|DT_WORDBREAK)
		drawText(hdc, b, RECT{int32(float64(w) * .15), int32(float64(h) * .51), int32(float64(w) * .85), int32(float64(h) * .64)}, 20, rgb(108, 110, 112), 400, "Segoe UI", DT_CENTER|DT_VCENTER|DT_WORDBREAK)
	case t < 35:
		drawCorridor(hdc, w, h, t, 0)
	case t < 36.7:
		drawCorridor(hdc, w, h, t, 1)
	case t < 38.75:
		drawCorridor(hdc, w, h, t, 2)
	case t < 44.5:
		fill(hdc, RECT{0, 0, w, h}, 0)
	case t < 54.4:
		drawFalseEnd(hdc, w, h, t, 1)
	case t < 61.9:
		drawMonitorGrid(hdc, w, h, t)
	case t < 70.8:
		fill(hdc, RECT{0, 0, w, h}, 0)
	case t < 78.9:
		drawRecoveredFrame(hdc, w, h, t)
	case t < 83.75:
		fill(hdc, RECT{0, 0, w, h}, 0)
		drawText(hdc, "последний кадр был снят после отключения камер.", RECT{0, int32(float64(h) * .4), w, int32(float64(h) * .56)}, 27, rgb(176, 176, 178), 600, "Segoe UI", DT_CENTER|DT_VCENTER|DT_WORDBREAK)
	case t < 92.2:
		fill(hdc, RECT{0, 0, w, h}, 0)
	case t < 104.9:
		drawFalseEnd(hdc, w, h, t, 2)
	case t < 114.7:
		fill(hdc, RECT{0, 0, w, h}, 0)
	case t < 118.2:
		fill(hdc, RECT{0, 0, w, h}, 0)
		drawText(hdc, "@SNAPWAVE", RECT{0, int32(float64(h) * .31), w, int32(float64(h) * .59)}, 98, rgb(239, 239, 241), 900, "Segoe UI", DT_CENTER|DT_VCENTER|DT_SINGLELINE)
	default:
		mode = 1
		cur, _, _ := pLoadCursorW.Call(0, IDC_ARROW)
		pSetCursor.Call(cur)
		drawLock(hdc, w, h)
	}
}
func checkPassword() {
	if password == correctPassword {
		mode = 3
		successSince = time.Now()
	} else {
		mode = 2
	}
	pInvalidateRect.Call(hwnd, 0, 0)
}
func wndProc(h uintptr, msg uint32, wp, lp uintptr) uintptr {
	switch msg {
	case WM_PAINT:
		var ps PAINTSTRUCT
		dc, _, _ := pBeginPaint.Call(h, uintptr(unsafe.Pointer(&ps)))
		var rc RECT
		pGetClientRect.Call(h, uintptr(unsafe.Pointer(&rc)))
		paint(dc, rc)
		pEndPaint.Call(h, uintptr(unsafe.Pointer(&ps)))
		return 0
	case WM_TIMER:
		if mode == 3 && time.Since(successSince) > 2200*time.Millisecond {
			pPostQuitMessage.Call(0)
		}
		pInvalidateRect.Call(h, 0, 0)
		return 0
	case WM_SETCURSOR:
		if mode == 1 {
			c, _, _ := pLoadCursorW.Call(0, IDC_ARROW)
			pSetCursor.Call(c)
		} else {
			pSetCursor.Call(0)
		}
		return 1
	case WM_CHAR:
		if mode == 1 {
			ch := rune(wp)
			if ch == 8 && len(password) > 0 {
				password = password[:len(password)-1]
			} else if ch == 13 {
				checkPassword()
			} else if ch >= 32 && ch < 127 && len(password) < 32 {
				password += string(ch)
			}
			pInvalidateRect.Call(h, 0, 0)
		}
		return 0
	case WM_LBUTTONDOWN:
		if mode == 1 {
			var rc RECT
			pGetClientRect.Call(h, uintptr(unsafe.Pointer(&rc)))
			w, hh := rc.Right, rc.Bottom
			pw := int32(math.Min(float64(w)*.66, 860))
			ph := int32(math.Min(float64(hh)*.58, 520))
			px, py := (w-pw)/2, (hh-ph)/2
			ix := px + int32(float64(pw)*.13)
			iw := int32(float64(pw) * .57)
			iy := py + 270
			bx := ix + iw + 14
			bw := pw - (bx - px) - int32(float64(pw)*.13)
			x := int32(int16(lp & 0xffff))
			y := int32(int16((lp >> 16) & 0xffff))
			if x >= bx && x <= bx+bw && y >= iy && y <= iy+58 {
				checkPassword()
			}
		}
		return 0
	case WM_CLOSE, WM_DESTROY:
		pPostQuitMessage.Call(0)
		return 0
	}
	r, _, _ := pDefWindowProcW.Call(h, uintptr(msg), wp, lp)
	return r
}
func main() {
	started = time.Now()
	inst, _, _ := kernel32.NewProc("GetModuleHandleW").Call(0)
	cls, _ := syscall.UTF16PtrFromString("SnapwaveV8Window")
	wc := WNDCLASSEX{CbSize: uint32(unsafe.Sizeof(WNDCLASSEX{})), LpfnWndProc: syscall.NewCallback(wndProc), HInstance: inst, LpszClassName: cls}
	pRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))
	sw, _, _ := pGetSystemMetrics.Call(SM_CXSCREEN)
	sh, _, _ := pGetSystemMetrics.Call(SM_CYSCREEN)
	title, _ := syscall.UTF16PtrFromString("SNAPWAVE")
	hwnd, _, _ = pCreateWindowExW.Call(WS_EX_TOPMOST, uintptr(unsafe.Pointer(cls)), uintptr(unsafe.Pointer(title)), WS_POPUP, 0, 0, sw, sh, 0, 0, inst, 0)
	pSetWindowPos.Call(hwnd, HWND_TOPMOST, 0, 0, sw, sh, SWP_SHOWWINDOW)
	pShowWindow.Call(hwnd, SW_SHOW)
	pUpdateWindow.Call(hwnd)
	playSoundtrack()
	pSetTimer.Call(hwnd, 1, 33, 0)
	var m MSG
	for {
		r, _, _ := pGetMessageW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
		if int32(r) <= 0 {
			break
		}
		pTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
		pDispatchMessageW.Call(uintptr(unsafe.Pointer(&m)))
	}
}
