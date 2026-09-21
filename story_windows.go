package main

import (
	"math"
	"time"
)

func drawVignette(hdc uintptr, w, h int32) {
	fill(hdc, RECT{0, 0, w, h}, rgb(3, 3, 4))
	for i := 0; i < 6; i++ {
		c := byte(3 + i*2)
		th := int32(18 + i*12)
		fill(hdc, RECT{0, 0, w, th}, rgb(c, c, c))
		fill(hdc, RECT{0, h - th, w, h}, rgb(c, c, c))
		fill(hdc, RECT{0, 0, th, h}, rgb(c, c, c))
		fill(hdc, RECT{w - th, 0, w, h}, rgb(c, c, c))
	}
}
func drawSilhouette(hdc uintptr, cx, cy, height int32) {
	head := height / 5
	body := height * 3 / 5
	ellipse(hdc, cx-head/2, cy-head, cx+head/2, cy, rgb(0, 0, 0))
	polygon(hdc, []POINT{{cx - head/2, cy}, {cx + head/2, cy}, {cx + height/5, cy + body}, {cx - height/5, cy + body}}, rgb(0, 0, 0))
}
func drawCameraHeader(hdc uintptr, w, h int32, cam string) {
	drawText(hdc, cam, RECT{28, 18, w - 28, 52}, 13, rgb(165, 168, 170), 600, "Consolas", DT_LEFT|DT_VCENTER|DT_SINGLELINE)
	fill(hdc, RECT{w - 118, 28, w - 106, 40}, rgb(190, 14, 20))
	drawText(hdc, "REC", RECT{w - 100, 18, w - 28, 52}, 12, rgb(190, 190, 192), 700, "Consolas", DT_LEFT|DT_VCENTER|DT_SINGLELINE)
	sec := int(time.Since(started).Seconds()) % 10
	drawText(hdc, "03:17:0"+string(rune('0'+sec)), RECT{28, h - 52, w - 28, h - 18}, 12, rgb(105, 108, 110), 500, "Consolas", DT_RIGHT|DT_VCENTER|DT_SINGLELINE)
}
func drawCorridor(hdc uintptr, w, h int32, t float64, stage int) {
	fill(hdc, RECT{0, 0, w, h}, rgb(4, 5, 6))
	polygon(hdc, []POINT{{0, 0}, {w, 0}, {int32(float64(w) * .68), int32(float64(h) * .62)}, {int32(float64(w) * .32), int32(float64(h) * .62)}}, rgb(17, 19, 20))
	polygon(hdc, []POINT{{0, h}, {w, h}, {int32(float64(w) * .68), int32(float64(h) * .62)}, {int32(float64(w) * .32), int32(float64(h) * .62)}}, rgb(8, 9, 10))
	p := pen(hdc, rgb(42, 44, 45), 2)
	for i := 0; i < 8; i++ {
		x := int32(float64(w) * (.06 + float64(i)*.125))
		line(hdc, x, h, w/2+int32((float64(i)-3.5)*38), int32(float64(h)*.62))
	}
	for i := 0; i < 6; i++ {
		y := int32(float64(h) * (.66 + float64(i)*.06))
		line(hdc, 0, y, w, y)
	}
	p()
	fw := int32(float64(w) * .13)
	fh := int32(float64(h) * .36)
	fx := w/2 - fw/2
	fy := int32(float64(h) * .26)
	rect(hdc, fx, fy, fx+fw, fy+fh, rgb(0, 0, 0))
	if stage > 0 {
		scale := float64(stage)*.47 + .45
		hh := int32(float64(h) * .19 * scale)
		cx := w/2 + int32(math.Sin(t*2.1)*3)
		drawSilhouette(hdc, cx, int32(float64(h)*.60)-hh/2, hh)
	}
	for i := 0; i < 20; i++ {
		if (i+int(t*10))%3 == 0 {
			y := int32((i*71 + int(t*55)) % int(h))
			fill(hdc, RECT{0, y, w, y + 1}, rgb(32, 34, 35))
		}
	}
	drawCameraHeader(hdc, w, h, "CAM-0"+string(rune('1'+stage))+" // SECTOR 09")
}
func drawMonitorGrid(hdc uintptr, w, h int32, t float64) {
	fill(hdc, RECT{0, 0, w, h}, rgb(2, 3, 4))
	gap := int32(18)
	mw := (w - gap*3) / 2
	mh := (h - gap*3) / 2
	for i := 0; i < 4; i++ {
		x := gap + int32(i%2)*(mw+gap)
		y := gap + int32(i/2)*(mh+gap)
		fill(hdc, RECT{x, y, x + mw, y + mh}, rgb(8, 9, 10))
		drawText(hdc, "CAM-0"+string(rune('1'+i)), RECT{x + 14, y + 10, x + mw - 14, y + 36}, 11, rgb(118, 120, 122), 600, "Consolas", DT_LEFT|DT_VCENTER|DT_SINGLELINE)
		if (i == 2 && t > 56.5) || (i == 0 && t > 58.2) || (i == 3 && t > 59.7) {
			drawSilhouette(hdc, x+mw/2, y+mh/2, int32(float64(mh)*.28))
		}
	}
	drawText(hdc, "ALL CHANNELS SYNCHRONIZED", RECT{0, h - 42, w, h - 10}, 11, rgb(89, 92, 94), 500, "Consolas", DT_CENTER|DT_VCENTER|DT_SINGLELINE)
}
func drawRecoveredFrame(hdc uintptr, w, h int32, t float64) {
	fill(hdc, RECT{0, 0, w, h}, rgb(1, 1, 1))
	drawCorridorImage(hdc, 0, 0, w, h)
	drawText(hdc, "RECOVERED FRAME // 03:17:11.042", RECT{28, 20, w - 28, 55}, 12, rgb(190, 190, 192), 600, "Consolas", DT_LEFT|DT_VCENTER|DT_SINGLELINE)
	if t > 77 {
		drawText(hdc, "FACE MATCH: VIEWER", RECT{0, int32(float64(h) * .79), w, int32(float64(h) * .87)}, 24, rgb(180, 12, 18), 800, "Consolas", DT_CENTER|DT_VCENTER|DT_SINGLELINE)
	}
}
func drawFalseEnd(hdc uintptr, w, h int32, t float64, n int) {
	fill(hdc, RECT{0, 0, w, h}, rgb(2, 2, 3))
	if n == 1 {
		drawText(hdc, "RECORDING ENDED", RECT{0, int32(float64(h) * .37), w, int32(float64(h) * .49)}, 39, rgb(218, 218, 220), 700, "Consolas", DT_CENTER|DT_VCENTER|DT_SINGLELINE)
		drawText(hdc, "No further anomalies detected.", RECT{0, int32(float64(h) * .51), w, int32(float64(h) * .59)}, 18, rgb(90, 92, 94), 400, "Segoe UI", DT_CENTER|DT_VCENTER|DT_SINGLELINE)
	} else {
		drawText(hdc, "CAMERA SERVICE RESTORED", RECT{0, int32(float64(h) * .34), w, int32(float64(h) * .45)}, 34, rgb(218, 218, 220), 700, "Consolas", DT_CENTER|DT_VCENTER|DT_SINGLELINE)
		drawText(hdc, "4 / 4 ONLINE", RECT{0, int32(float64(h) * .46), w, int32(float64(h) * .57)}, 54, rgb(232, 232, 234), 800, "Consolas", DT_CENTER|DT_VCENTER|DT_SINGLELINE)
		if t > 101.3 {
			drawText(hdc, "одна камера всё ещё смотрит.", RECT{0, int32(float64(h) * .64), w, int32(float64(h) * .72)}, 21, rgb(142, 21, 26), 600, "Segoe UI", DT_CENTER|DT_VCENTER|DT_SINGLELINE)
		}
	}
}
