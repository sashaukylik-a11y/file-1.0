package main

import (
	"syscall"
	"unsafe"
)

var (
	user32 = syscall.NewLazyDLL("user32.dll")
	gdi32 = syscall.NewLazyDLL("gdi32.dll")
	winmm = syscall.NewLazyDLL("winmm.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	pRtlMoveMemory = kernel32.NewProc("RtlMoveMemory")

	pRegisterClassExW = user32.NewProc("RegisterClassExW")
	pCreateWindowExW = user32.NewProc("CreateWindowExW")
	pDefWindowProcW = user32.NewProc("DefWindowProcW")
	pShowWindow = user32.NewProc("ShowWindow")
	pUpdateWindow = user32.NewProc("UpdateWindow")
	pGetMessageW = user32.NewProc("GetMessageW")
	pTranslateMessage = user32.NewProc("TranslateMessage")
	pDispatchMessageW = user32.NewProc("DispatchMessageW")
	pPostQuitMessage = user32.NewProc("PostQuitMessage")
	pBeginPaint = user32.NewProc("BeginPaint")
	pEndPaint = user32.NewProc("EndPaint")
	pGetClientRect = user32.NewProc("GetClientRect")
	pInvalidateRect = user32.NewProc("InvalidateRect")
	pSetTimer = user32.NewProc("SetTimer")
	pGetSystemMetrics = user32.NewProc("GetSystemMetrics")
	pSetWindowPos = user32.NewProc("SetWindowPos")
	pDestroyWindow = user32.NewProc("DestroyWindow")
	pKillTimer = user32.NewProc("KillTimer")
	pSetProcessDPIAware = user32.NewProc("SetProcessDPIAware")
	pLoadCursorW = user32.NewProc("LoadCursorW")
	pSetCursor = user32.NewProc("SetCursor")
	pSetFocus = user32.NewProc("SetFocus")
	pFillRect = user32.NewProc("FillRect")
	pDrawTextW = user32.NewProc("DrawTextW")

	pCreateSolidBrush = gdi32.NewProc("CreateSolidBrush")
	pDeleteObject = gdi32.NewProc("DeleteObject")
	pSetTextColor = gdi32.NewProc("SetTextColor")
	pSetBkMode = gdi32.NewProc("SetBkMode")
	pCreateFontW = gdi32.NewProc("CreateFontW")
	pSelectObject = gdi32.NewProc("SelectObject")
	pCreatePen = gdi32.NewProc("CreatePen")
	pMoveToEx = gdi32.NewProc("MoveToEx")
	pLineTo = gdi32.NewProc("LineTo")
	pEllipse = gdi32.NewProc("Ellipse")
	pRectangle = gdi32.NewProc("Rectangle")
	pPolygon = gdi32.NewProc("Polygon")
	pStretchDIBits = gdi32.NewProc("StretchDIBits")
	pCreateCompatibleDC = gdi32.NewProc("CreateCompatibleDC")
	pCreateCompatibleBitmap = gdi32.NewProc("CreateCompatibleBitmap")
	pDeleteDC = gdi32.NewProc("DeleteDC")
	pBitBlt = gdi32.NewProc("BitBlt")
	pStretchBlt = gdi32.NewProc("StretchBlt")
	pSetStretchBltMode = gdi32.NewProc("SetStretchBltMode")
	pCreateDIBSection = gdi32.NewProc("CreateDIBSection")
	pPlaySoundW = winmm.NewProc("PlaySoundW")
)

const (
	WM_DESTROY=0x0002
	WM_CLOSE=0x0010
	WM_PAINT=0x000F
	WM_ERASEBKGND=0x0014
	WM_TIMER=0x0113
	WM_CHAR=0x0102
	WM_SETCURSOR=0x0020
	WM_LBUTTONDOWN=0x0201
	WS_POPUP=0x80000000
	WS_EX_TOPMOST=0x00000008
	SW_SHOW=5
	HWND_TOPMOST=^uintptr(0)
	SWP_SHOWWINDOW=0x0040
	SM_CXSCREEN=0
	SM_CYSCREEN=1
	TRANSPARENT=1
	PS_SOLID=0
	IDC_ARROW=32512
	DIB_RGB_COLORS=0
	SRCCOPY=0x00CC0020
	COLORONCOLOR=3
	DT_CENTER=0x1
	DT_VCENTER=0x4
	DT_SINGLELINE=0x20
	DT_WORDBREAK=0x10
	DT_LEFT=0
	DT_RIGHT=0x2
	SND_ASYNC=0x1
	SND_NODEFAULT=0x2
	SND_MEMORY=0x4
)

type WNDCLASSEX struct {
	CbSize, Style uint32
	LpfnWndProc uintptr
	CbClsExtra, CbWndExtra int32
	HInstance, HIcon, HCursor, HbrBackground uintptr
	LpszMenuName, LpszClassName *uint16
	HIconSm uintptr
}
type POINT struct{ X,Y int32 }
type RECT struct{ Left,Top,Right,Bottom int32 }
type PAINTSTRUCT struct {
	Hdc uintptr
	FErase int32
	RcPaint RECT
	FRestore,FIncUpdate int32
	RgbReserved [32]byte
}
type MSG struct {
	Hwnd uintptr
	Message uint32
	WParam,LParam uintptr
	Time uint32
	Pt POINT
	LPrivate uint32
}
type fontKey struct{ size,weight int32; face string }
type penKey struct{ color uintptr; width int32 }

var (
	fontCache=map[fontKey]uintptr{}
	brushCache=map[uintptr]uintptr{}
	penCache=map[penKey]uintptr{}
	backDC uintptr
	backBitmap uintptr
	backOld uintptr
	backW int32
	backH int32
)

func rgb(r,g,b byte) uintptr { return uintptr(uint32(r)|uint32(g)<<8|uint32(b)<<16) }
func cachedBrush(c uintptr) uintptr {
	if b:=brushCache[c]; b!=0 { return b }
	b,_,_:=pCreateSolidBrush.Call(c); brushCache[c]=b; return b
}
func cachedPen(c uintptr,w int32) uintptr {
	k:=penKey{c,w}; if p:=penCache[k]; p!=0 { return p }
	p,_,_:=pCreatePen.Call(PS_SOLID,uintptr(w),c); penCache[k]=p; return p
}
func cachedFont(size,weight int32,face string) uintptr {
	k:=fontKey{size,weight,face}; if f:=fontCache[k]; f!=0 { return f }
	fp,_:=syscall.UTF16PtrFromString(face)
	f,_,_:=pCreateFontW.Call(uintptr(-size),0,0,0,uintptr(weight),0,0,0,1,0,0,5,0,uintptr(unsafe.Pointer(fp)))
	fontCache[k]=f; return f
}
func cleanupGDI() {
	cleanupBackBuffer()
	for _,f:=range fontCache { pDeleteObject.Call(f) }
	for _,b:=range brushCache { pDeleteObject.Call(b) }
	for _,p:=range penCache { pDeleteObject.Call(p) }
	fontCache=map[fontKey]uintptr{}; brushCache=map[uintptr]uintptr{}; penCache=map[penKey]uintptr{}
}
func ensureBackBuffer(screenDC uintptr,w,h int32) uintptr {
	if w<=0||h<=0 { return screenDC }
	if backDC!=0&&backBitmap!=0&&backW==w&&backH==h { return backDC }
	cleanupBackBuffer()
	dc,_,_:=pCreateCompatibleDC.Call(screenDC); if dc==0 { return screenDC }
	bmp,_,_:=pCreateCompatibleBitmap.Call(screenDC,uintptr(w),uintptr(h))
	if bmp==0 { pDeleteDC.Call(dc); return screenDC }
	old,_,_:=pSelectObject.Call(dc,bmp)
	if old==0 { pDeleteObject.Call(bmp); pDeleteDC.Call(dc); return screenDC }
	backDC,backBitmap,backOld=dc,bmp,old; backW,backH=w,h; return backDC
}
func cleanupBackBuffer() {
	if backDC==0 { return }
	if backOld!=0 { pSelectObject.Call(backDC,backOld) }
	if backBitmap!=0 { pDeleteObject.Call(backBitmap) }
	pDeleteDC.Call(backDC)
	backDC,backBitmap,backOld=0,0,0; backW,backH=0,0
}
func presentBackBuffer(screenDC uintptr,w,h int32) {
	if backDC==0||backBitmap==0 { return }
	pBitBlt.Call(screenDC,0,0,uintptr(w),uintptr(h),backDC,0,0,SRCCOPY)
}
func fill(hdc uintptr,r RECT,c uintptr) { pFillRect.Call(hdc,uintptr(unsafe.Pointer(&r)),cachedBrush(c)) }
func pen(hdc uintptr,c uintptr,w int32) func() {
	p:=cachedPen(c,w); old,_,_:=pSelectObject.Call(hdc,p); return func(){ pSelectObject.Call(hdc,old) }
}
func line(hdc uintptr,x1,y1,x2,y2 int32) { pMoveToEx.Call(hdc,uintptr(x1),uintptr(y1),0); pLineTo.Call(hdc,uintptr(x2),uintptr(y2)) }
func ellipse(hdc uintptr,l,t,r,b int32,c uintptr) {
	br:=cachedBrush(c); old,_,_:=pSelectObject.Call(hdc,br); pEllipse.Call(hdc,uintptr(l),uintptr(t),uintptr(r),uintptr(b)); pSelectObject.Call(hdc,old)
}
func rect(hdc uintptr,l,t,r,b int32,c uintptr) {
	br:=cachedBrush(c); old,_,_:=pSelectObject.Call(hdc,br); pRectangle.Call(hdc,uintptr(l),uintptr(t),uintptr(r),uintptr(b)); pSelectObject.Call(hdc,old)
}
func polygon(hdc uintptr,pts []POINT,c uintptr) {
	if len(pts)==0 { return }
	br:=cachedBrush(c); old,_,_:=pSelectObject.Call(hdc,br); pPolygon.Call(hdc,uintptr(unsafe.Pointer(&pts[0])),uintptr(len(pts))); pSelectObject.Call(hdc,old)
}
func drawText(hdc uintptr,text string,r RECT,size int32,color uintptr,weight int32,face string,flags uintptr) {
	f:=cachedFont(size,weight,face); old,_,_:=pSelectObject.Call(hdc,f)
	pSetTextColor.Call(hdc,color); pSetBkMode.Call(hdc,TRANSPARENT)
	tp,_:=syscall.UTF16PtrFromString(text)
	pDrawTextW.Call(hdc,uintptr(unsafe.Pointer(tp)),^uintptr(0),uintptr(unsafe.Pointer(&r)),flags)
	pSelectObject.Call(hdc,old)
}
