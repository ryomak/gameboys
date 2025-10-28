package graphics

import "unsafe"

const (
	ScreenWidth  = 240
	ScreenHeight = 160
	VRAMBase     = 0x06000000
)

// VideoBuffer Mode 3用のビデオバッファ（240x160、16bit color）
var VideoBuffer = (*[ScreenWidth * ScreenHeight]uint16)(unsafe.Pointer(uintptr(VRAMBase)))

// DrawPixel ピクセルを描画（Mode 3用）
func DrawPixel(x, y int, color uint16) {
	if x >= 0 && x < ScreenWidth && y >= 0 && y < ScreenHeight {
		VideoBuffer[y*ScreenWidth+x] = color
	}
}

// GetPixel ピクセルの色を取得（Mode 3用）
func GetPixel(x, y int) uint16 {
	if x >= 0 && x < ScreenWidth && y >= 0 && y < ScreenHeight {
		return VideoBuffer[y*ScreenWidth+x]
	}
	return 0
}

// ClearScreen 画面をクリア
func ClearScreen(color uint16) {
	for i := range VideoBuffer {
		VideoBuffer[i] = color
	}
}

// FillRect 矩形を塗りつぶす
func FillRect(x, y, width, height int, color uint16) {
	// クリッピング
	if x < 0 {
		width += x
		x = 0
	}
	if y < 0 {
		height += y
		y = 0
	}
	if x+width > ScreenWidth {
		width = ScreenWidth - x
	}
	if y+height > ScreenHeight {
		height = ScreenHeight - y
	}
	if width <= 0 || height <= 0 {
		return
	}

	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			VideoBuffer[(y+j)*ScreenWidth+(x+i)] = color
		}
	}
}

// DrawHLine 水平線を描画
func DrawHLine(x, y, length int, color uint16) {
	if y < 0 || y >= ScreenHeight {
		return
	}
	if x < 0 {
		length += x
		x = 0
	}
	if x+length > ScreenWidth {
		length = ScreenWidth - x
	}
	if length <= 0 {
		return
	}

	offset := y * ScreenWidth
	for i := 0; i < length; i++ {
		VideoBuffer[offset+x+i] = color
	}
}

// DrawVLine 垂直線を描画
func DrawVLine(x, y, length int, color uint16) {
	if x < 0 || x >= ScreenWidth {
		return
	}
	if y < 0 {
		length += y
		y = 0
	}
	if y+length > ScreenHeight {
		length = ScreenHeight - y
	}
	if length <= 0 {
		return
	}

	for i := 0; i < length; i++ {
		VideoBuffer[(y+i)*ScreenWidth+x] = color
	}
}
