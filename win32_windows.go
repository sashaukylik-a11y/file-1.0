package main

import (
	"syscall"
	"unsafe"
)

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	gdi32    = syscall.NewLazyDLL("gdi32.dll")
	winmm    = syscall.NewLazyDLL("winmm.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")

	pRegisterClassExW = user32.NewProc("RegisterClassExW")
	pCreateWindowExW  = user32.NewProc("CreateWindowExW")
	pDefWindowProcW   = user32.NewProc("DefWindowProcW")
	pShowWindow       = user32.NewProc("ShowWindow")
	pUpdateWindow     = user32.NewProc("UpdateWindow")
	pGetMessageW      = user32.NewProc("GetMessageW")
	pTranslateMessage = user32.NewProc("TranslateMessage")
	pDispatchMessageW = user32.NewProc("DispatchMessageW")
	pPostQuitMessage  = user32.NewProc("PostQuitMessage")
	pBeginPaint       = user32.NewProc("BeginPaint")
	pEndPaint         = user32.NewProc("EndPaint")
	pGetClientRect    = user32.NewProc("GetClientRect")
	pInvalidateRect   = user32.NewProc("InvalidateRect")
	pSetTimer         = user32.NewProc("SetTimer")
	pGetSystemMetrics = user32.NewProc("GetSystemMetrics")
	pSetWindowPos     = user32.NewProc("SetWindowPos")
	pLoadCursorW      = user32.NewProc("LoadCursorW")
	pSetCursor        = user32.NewProc("SetCursor")
	pFillRect         = user32.NewProc("FillRect")
	pDrawTextW        = user32.NewProc("DrawTextW")

	pCreateSolidBrush = gdi32.NewProc("CreateSolidBrush")
	pDeleteObject     = gdi32.NewProc("DeleteObject")
	pSetTextColor     = gdi32.NewProc("SetTextColor")
	pSetBkMode        = gdi32.NewProc("SetBkMode")
	pCreateFontW      = gdi32.NewProc("CreateFontW")
	pSelectObject     = gdi32.NewProc("SelectObject")
	pCreatePen        = gdi32.NewProc("CreatePen")
	pMoveToEx         = gdi32.NewProc("MoveToEx")
	pLineTo           = gdi32.NewProc("LineTo")
	pEllipse          = gdi32.NewProc("Ellipse")
	pRectangle        = gdi32.NewProc("Rectangle")
	pPolygon          = gdi32.NewProc("Polygon")
	pStretchDIBits    = gdi32.NewProc("StretchDIBits")
	pPlaySoundW       = winmm.NewProc("PlaySoundW")
)

const (
	WM_DESTROY     = 0x0002
	WM_CLOSE       = 0x0010
	WM_PAINT       = 0x000F
	WM_ERASEBKGND  = 0x0014
	WM_TIMER       = 0x0113
	WM_CHAR        = 0x0102
	WM_SETCURSOR   = 0x0020
	WM_LBUTTONDOWN = 0x0201
	WS_POPUP       = 0x80000000
	WS_EX_TOPMOST  = 0x00000008
	SW_SHOW        = 5
	HWND_TOPMOST   = ^uintptr(0)
	SWP_SHOWWINDOW = 0x0040
	SM_CXSCREEN    = 0
	SM_CYSCREEN    = 1
	TRANSPARENT    = 1
	PS_SOLID       = 0
	IDC_ARROW      = 32512
	DIB_RGB_COLORS = 0
	SRCCOPY        = 0x00CC0020
	DT_CENTER      = 0x1
	DT_VCENTER     = 0x4
	DT_SINGLELINE  = 0x20
	DT_WORDBREAK   = 0x10
	DT_LEFT        = 0
	DT_RIGHT       = 0x2
	SND_ASYNC      = 0x1
	SND_NODEFAULT  = 0x2
	SND_MEMORY     = 0x4
)

type WNDCLASSEX struct {
	CbSize, Style                            uint32
	LpfnWndProc                              uintptr
	CbClsExtra, CbWndExtra                   int32
	HInstance, HIcon, HCursor, HbrBackground uintptr
	LpszMenuName, LpszClassName              *uint16
	HIconSm                                  uintptr
}
type POINT struct{ X, Y int32 }
type RECT struct{ Left, Top, Right, Bottom int32 }
type PAINTSTRUCT struct {
	Hdc                  uintptr
	FErase               int32
	RcPaint              RECT
	FRestore, FIncUpdate int32
	RgbReserved          [32]byte
}
type MSG struct {
	Hwnd           uintptr
	Message        uint32
	WParam, LParam uintptr
	Time           uint32
	Pt             POINT
	LPrivate       uint32
}

func rgb(r, g, b byte) uintptr { return uintptr(uint32(r) | uint32(g)<<8 | uint32(b)<<16) }
func fill(hdc uintptr, r RECT, c uintptr) {
	b, _, _ := pCreateSolidBrush.Call(c)
	pFillRect.Call(hdc, uintptr(unsafe.Pointer(&r)), b)
	pDeleteObject.Call(b)
}
func pen(hdc uintptr, c uintptr, w int32) func() {
	p, _, _ := pCreatePen.Call(PS_SOLID, uintptr(w), c)
	old, _, _ := pSelectObject.Call(hdc, p)
	return func() { pSelectObject.Call(hdc, old); pDeleteObject.Call(p) }
}
func line(hdc uintptr, x1, y1, x2, y2 int32) {
	pMoveToEx.Call(hdc, uintptr(x1), uintptr(y1), 0)
	pLineTo.Call(hdc, uintptr(x2), uintptr(y2))
}
func ellipse(hdc uintptr, l, t, r, b int32, c uintptr) {
	br, _, _ := pCreateSolidBrush.Call(c)
	old, _, _ := pSelectObject.Call(hdc, br)
	pEllipse.Call(hdc, uintptr(l), uintptr(t), uintptr(r), uintptr(b))
	pSelectObject.Call(hdc, old)
	pDeleteObject.Call(br)
}
func rect(hdc uintptr, l, t, r, b int32, c uintptr) {
	br, _, _ := pCreateSolidBrush.Call(c)
	old, _, _ := pSelectObject.Call(hdc, br)
	pRectangle.Call(hdc, uintptr(l), uintptr(t), uintptr(r), uintptr(b))
	pSelectObject.Call(hdc, old)
	pDeleteObject.Call(br)
}
func polygon(hdc uintptr, pts []POINT, c uintptr) {
	br, _, _ := pCreateSolidBrush.Call(c)
	old, _, _ := pSelectObject.Call(hdc, br)
	pPolygon.Call(hdc, uintptr(unsafe.Pointer(&pts[0])), uintptr(len(pts)))
	pSelectObject.Call(hdc, old)
	pDeleteObject.Call(br)
}
func drawText(hdc uintptr, text string, r RECT, size int32, color uintptr, weight int32, face string, flags uintptr) {
	fp, _ := syscall.UTF16PtrFromString(face)
	f, _, _ := pCreateFontW.Call(uintptr(-size), 0, 0, 0, uintptr(weight), 0, 0, 0, 1, 0, 0, 5, 0, uintptr(unsafe.Pointer(fp)))
	old, _, _ := pSelectObject.Call(hdc, f)
	pSetTextColor.Call(hdc, color)
	pSetBkMode.Call(hdc, TRANSPARENT)
	tp, _ := syscall.UTF16PtrFromString(text)
	pDrawTextW.Call(hdc, uintptr(unsafe.Pointer(tp)), ^uintptr(0), uintptr(unsafe.Pointer(&r)), flags)
	pSelectObject.Call(hdc, old)
	pDeleteObject.Call(f)
}
