package main

import (
	_ "embed"
	"encoding/binary"
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

func drawBMP(hdc uintptr, data []byte, x, y, w, h int32) bool {
	if len(data) < 54 || string(data[:2]) != "BM" {
		return false
	}
	off := binary.LittleEndian.Uint32(data[10:14])
	if int(off) >= len(data) {
		return false
	}
	sw := int32(binary.LittleEndian.Uint32(data[18:22]))
	sh := int32(binary.LittleEndian.Uint32(data[22:26]))
	if sw <= 0 || sh == 0 {
		return false
	}
	pStretchDIBits.Call(hdc, uintptr(x), uintptr(y), uintptr(w), uintptr(h), 0, 0, uintptr(sw), uintptr(sh), uintptr(unsafe.Pointer(&data[off])), uintptr(unsafe.Pointer(&data[14])), DIB_RGB_COLORS, SRCCOPY)
	return true
}
func playSoundtrack() {
	if len(soundtrack) > 0 {
		pPlaySoundW.Call(uintptr(unsafe.Pointer(&soundtrack[0])), 0, SND_MEMORY|SND_ASYNC|SND_NODEFAULT)
	}
}
