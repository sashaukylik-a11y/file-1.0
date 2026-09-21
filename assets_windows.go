package main

import (
	_ "embed"
	"encoding/binary"
	"sync"
	"unsafe"
)

//go:embed assets/soundtrack.wav
var soundtrack []byte
//go:embed assets/scare_1.bmp
var scare1 []byte
//go:embed assets/scare_2.bmp
var scare2 []byte
//go:embed assets/scare_3.bmp
var scare3 []byte
//go:embed assets/scare_4.bmp
var scare4 []byte
//go:embed assets/scare_5.bmp
var scare5 []byte
//go:embed assets/scare_6.bmp
var scare6 []byte
//go:embed assets/corridor_figure.bmp
var corridorFigure []byte

var scareImages = [][]byte{scare1, scare2, scare3, scare4, scare5, scare6}

type cachedDIB struct {
	dc, bitmap, old uintptr
	w, h            int32
	valid           bool
}

var (
	assetOnce sync.Once
	cachedScares [6]cachedDIB
	cachedCorridor cachedDIB
)

func parseBMP(data []byte) (off uint32, w, h int32, bitCount uint16, ok bool) {
	if len(data) < 54 || data[0] != 'B' || data[1] != 'M' { return 0,0,0,0,false }
	off = binary.LittleEndian.Uint32(data[10:14])
	headerSize := binary.LittleEndian.Uint32(data[14:18])
	if headerSize < 40 || int(off) >= len(data) { return 0,0,0,0,false }
	w = int32(binary.LittleEndian.Uint32(data[18:22]))
	h = int32(binary.LittleEndian.Uint32(data[22:26]))
	planes := binary.LittleEndian.Uint16(data[26:28])
	bitCount = binary.LittleEndian.Uint16(data[28:30])
	compression := binary.LittleEndian.Uint32(data[30:34])
	if w <= 0 || h == 0 || planes != 1 || (bitCount != 24 && bitCount != 32) || compression != 0 { return 0,0,0,0,false }
	return off,w,h,bitCount,true
}

func cacheBMP(data []byte) cachedDIB {
	off,w,h,bitCount,ok := parseBMP(data)
	if !ok { return cachedDIB{} }
	var bits uintptr
	bmp,_,_ := pCreateDIBSection.Call(0, uintptr(unsafe.Pointer(&data[14])), DIB_RGB_COLORS, uintptr(unsafe.Pointer(&bits)), 0,0)
	if bmp == 0 || bits == 0 { return cachedDIB{} }
	absH := h
	if absH < 0 { absH = -absH }
	stride := ((int(w)*int(bitCount)+31)/32)*4
	pixelBytes := stride*int(absH)
	if pixelBytes <= 0 || int(off)+pixelBytes > len(data) { pDeleteObject.Call(bmp); return cachedDIB{} }
	pRtlMoveMemory.Call(bits, uintptr(unsafe.Pointer(&data[int(off)])), uintptr(pixelBytes))
	dc,_,_ := pCreateCompatibleDC.Call(0)
	if dc == 0 { pDeleteObject.Call(bmp); return cachedDIB{} }
	old,_,_ := pSelectObject.Call(dc,bmp)
	return cachedDIB{dc:dc, bitmap:bmp, old:old, w:w, h:h, valid:true}
}

func initEmbeddedBitmaps() {
	assetOnce.Do(func() {
		for i := range cachedScares { cachedScares[i] = cacheBMP(scareImages[i]) }
		cachedCorridor = cacheBMP(corridorFigure)
	})
}

func (b *cachedDIB) draw(hdc uintptr, x,y,w,h int32) bool {
	if !b.valid || b.dc == 0 || w <= 0 || h <= 0 { return false }
	pSetStretchBltMode.Call(hdc, COLORONCOLOR)
	pStretchBlt.Call(hdc, uintptr(x), uintptr(y), uintptr(w), uintptr(h), b.dc, 0,0, uintptr(b.w), uintptr(b.h), SRCCOPY)
	return true
}

func drawScare(hdc uintptr, index int, x,y,w,h int32) bool {
	initEmbeddedBitmaps()
	if index < 0 || index >= len(cachedScares) { return false }
	return cachedScares[index].draw(hdc,x,y,w,h)
}
func drawCorridorImage(hdc uintptr, x,y,w,h int32) bool {
	initEmbeddedBitmaps()
	return cachedCorridor.draw(hdc,x,y,w,h)
}
func cleanupEmbeddedBitmaps() {
	cleanup := func(b *cachedDIB) {
		if b.dc == 0 { return }
		if b.old != 0 { pSelectObject.Call(b.dc,b.old) }
		if b.bitmap != 0 { pDeleteObject.Call(b.bitmap) }
		pDeleteDC.Call(b.dc)
		*b = cachedDIB{}
	}
	for i := range cachedScares { cleanup(&cachedScares[i]) }
	cleanup(&cachedCorridor)
}
func playSoundtrack() bool {
	if len(soundtrack) < 44 || string(soundtrack[:4]) != "RIFF" || string(soundtrack[8:12]) != "WAVE" { return false }
	r,_,_ := pPlaySoundW.Call(uintptr(unsafe.Pointer(&soundtrack[0])),0,SND_MEMORY|SND_ASYNC|SND_NODEFAULT)
	return r != 0
}
func stopSoundtrack() { pPlaySoundW.Call(0,0,0) }
