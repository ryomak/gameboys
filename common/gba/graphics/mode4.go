package graphics

import (
	"runtime/volatile"
	"unsafe"
)

// Mode 4: 8bitカラー、240x160、ダブルバッファリング対応

const (
	Mode4Width  = 240
	Mode4Height = 160

	// VRAMアドレス
	VRAM4Base    = 0x06000000
	VRAM4Frame0  = 0x06000000 // フレーム0
	VRAM4Frame1  = 0x0600A000 // フレーム1（40KB = 0xA000バイト後）
	PaletteRAM   = 0x05000000 // パレットRAM
)

// currentDrawBuffer 現在の描画先バッファ（0 or 1）
// 表示バッファは逆のバッファになる
var currentDrawBuffer uint16 = 0

// GetMode4BackBuffer 現在のバックバッファ（描画先）のアドレスを取得
func GetMode4BackBuffer() uintptr {
	if currentDrawBuffer == 0 {
		return VRAM4Frame0
	}
	return VRAM4Frame1
}

// GetCurrentDrawBuffer 現在の描画先バッファ番号を取得（0 or 1）
func GetCurrentDrawBuffer() uint16 {
	return currentDrawBuffer
}

// SwapBuffers バッファを切り替える（内部状態のみ変更）
func SwapBuffers() {
	currentDrawBuffer = 1 - currentDrawBuffer
}

// SetMode4Pixel Mode 4でピクセルを設定（バックバッファに描画）
func SetMode4Pixel(x, y int, colorIndex uint8) {
	if x < 0 || x >= Mode4Width || y < 0 || y >= Mode4Height {
		return
	}

	addr := GetMode4BackBuffer()
	offset := uintptr(y*Mode4Width + x)
	ptr := (*volatile.Register8)(unsafe.Pointer(addr + offset))
	ptr.Set(colorIndex)
}

// GetMode4Pixel Mode 4のピクセルを取得（バックバッファから）
func GetMode4Pixel(x, y int) uint8 {
	if x < 0 || x >= Mode4Width || y < 0 || y >= Mode4Height {
		return 0
	}

	addr := GetMode4BackBuffer()
	offset := uintptr(y*Mode4Width + x)
	ptr := (*volatile.Register8)(unsafe.Pointer(addr + offset))
	return ptr.Get()
}

// ClearMode4Screen Mode 4の画面全体をクリア（バックバッファ）
func ClearMode4Screen(colorIndex uint8) {
	addr := GetMode4BackBuffer()

	// 16bit単位で高速クリア
	color16 := uint16(colorIndex) | (uint16(colorIndex) << 8)

	for i := uintptr(0); i < Mode4Width*Mode4Height/2; i++ {
		ptr := (*volatile.Register16)(unsafe.Pointer(addr + i*2))
		ptr.Set(color16)
	}
}

// SetMode4Palette パレットに色を設定
func SetMode4Palette(index uint8, color uint16) {
	addr := PaletteRAM + uintptr(index)*2
	ptr := (*volatile.Register16)(unsafe.Pointer(addr))
	ptr.Set(color)
}

// InitMode4Palette デフォルトパレットを初期化
func InitMode4Palette() {
	// 基本色を設定
	SetMode4Palette(0, ColorBlack)
	SetMode4Palette(1, ColorWhite)
	SetMode4Palette(2, ColorRed)
	SetMode4Palette(3, ColorGreen)
	SetMode4Palette(4, ColorBlue)
	SetMode4Palette(5, ColorYellow)
	SetMode4Palette(6, ColorCyan)
	SetMode4Palette(7, ColorMagenta)
	SetMode4Palette(8, ColorGray)
	SetMode4Palette(9, ColorDarkGray)
	SetMode4Palette(10, ColorOrange)

	// その他のカスタムカラー
	SetMode4Palette(11, RGB(220, 180, 140)) // 肌色
	SetMode4Palette(12, RGB(180, 140, 100)) // 暗い肌色
	SetMode4Palette(13, RGB(139, 90, 43))   // 床色（木目）
	SetMode4Palette(14, RGB(160, 110, 60))  // 明るい床色
	SetMode4Palette(15, RGB(30, 30, 50))    // 体育館の壁
	SetMode4Palette(16, RGB(200, 200, 200)) // バックボード
	SetMode4Palette(17, RGB(255, 100, 0))   // リム色（オレンジ）
	SetMode4Palette(18, RGB(255, 120, 0))   // ボール色
	SetMode4Palette(19, RGB(200, 80, 0))    // 暗いボール色
	SetMode4Palette(20, RGB(255, 160, 60))  // 明るいボール色
	SetMode4Palette(21, RGB(0, 150, 0))     // 成功色（濃い緑）
	SetMode4Palette(22, RGB(0, 200, 0))     // 成功色（明るい緑）
	SetMode4Palette(23, RGB(150, 0, 0))     // 失敗色（濃い赤）
	SetMode4Palette(24, RGB(100, 0, 0))     // 失敗色（暗い赤）
	SetMode4Palette(25, RGB(255, 200, 0))   // ゴールド
	SetMode4Palette(26, RGB(40, 40, 60))    // UI背景
	SetMode4Palette(27, RGB(60, 60, 60))    // ゲージ背景
	SetMode4Palette(28, RGB(0, 0, 100))     // 青い背景
	SetMode4Palette(29, RGB(200, 100, 0))   // オレンジ系
}

// パレットインデックス定数
const (
	PalBlack        = 0
	PalWhite        = 1
	PalRed          = 2
	PalGreen        = 3
	PalBlue         = 4
	PalYellow       = 5
	PalCyan         = 6
	PalMagenta      = 7
	PalGray         = 8
	PalDarkGray     = 9
	PalOrange       = 10
	PalSkin         = 11
	PalDarkSkin     = 12
	PalFloor        = 13
	PalLightFloor   = 14
	PalWall         = 15
	PalBackboard    = 16
	PalRim          = 17
	PalBall         = 18
	PalDarkBall     = 19
	PalLightBall    = 20
	PalSuccessDark  = 21
	PalSuccessLight = 22
	PalFailDark     = 23
	PalFailDarker   = 24
	PalGold         = 25
	PalUIBG         = 26
	PalGaugeBG      = 27
	PalBlueBG       = 28
	PalOrangeBG     = 29
)

// Mode 4用の描画関数

// DrawLineMode4 Mode 4で直線を描画
func DrawLineMode4(x0, y0, x1, y1 int, colorIndex uint8) {
	// Bresenhamの直線描画アルゴリズム
	dx := x1 - x0
	if dx < 0 {
		dx = -dx
	}
	dy := y1 - y0
	if dy < 0 {
		dy = -dy
	}

	var sx, sy int
	if x0 < x1 {
		sx = 1
	} else {
		sx = -1
	}
	if y0 < y1 {
		sy = 1
	} else {
		sy = -1
	}

	err := dx - dy

	for {
		SetMode4Pixel(x0, y0, colorIndex)

		if x0 == x1 && y0 == y1 {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

// FillRectMode4 Mode 4で矩形を塗りつぶし
func FillRectMode4(x, y, width, height int, colorIndex uint8) {
	if x < 0 {
		width += x
		x = 0
	}
	if y < 0 {
		height += y
		y = 0
	}
	if x >= Mode4Width || y >= Mode4Height {
		return
	}
	if x+width > Mode4Width {
		width = Mode4Width - x
	}
	if y+height > Mode4Height {
		height = Mode4Height - y
	}
	if width <= 0 || height <= 0 {
		return
	}

	addr := GetMode4BackBuffer()

	for row := 0; row < height; row++ {
		offset := uintptr((y+row)*Mode4Width + x)
		for col := 0; col < width; col++ {
			ptr := (*volatile.Register8)(unsafe.Pointer(addr + offset + uintptr(col)))
			ptr.Set(colorIndex)
		}
	}
}

// DrawRectMode4 Mode 4で矩形の枠を描画
func DrawRectMode4(x, y, width, height int, colorIndex uint8) {
	// 上
	DrawLineMode4(x, y, x+width-1, y, colorIndex)
	// 下
	DrawLineMode4(x, y+height-1, x+width-1, y+height-1, colorIndex)
	// 左
	DrawLineMode4(x, y, x, y+height-1, colorIndex)
	// 右
	DrawLineMode4(x+width-1, y, x+width-1, y+height-1, colorIndex)
}

// FillCircleMode4 Mode 4で円を塗りつぶし
func FillCircleMode4(cx, cy, radius int, colorIndex uint8) {
	if radius <= 0 {
		return
	}

	x := 0
	y := radius
	d := 3 - 2*radius

	for x <= y {
		// 水平線を描画して塗りつぶし
		DrawLineMode4(cx-x, cy+y, cx+x, cy+y, colorIndex)
		DrawLineMode4(cx-x, cy-y, cx+x, cy-y, colorIndex)
		DrawLineMode4(cx-y, cy+x, cx+y, cy+x, colorIndex)
		DrawLineMode4(cx-y, cy-x, cx+y, cy-x, colorIndex)

		if d < 0 {
			d = d + 4*x + 6
		} else {
			d = d + 4*(x-y) + 10
			y--
		}
		x++
	}
}

// DrawCircleMode4 Mode 4で円の輪郭を描画
func DrawCircleMode4(cx, cy, radius int, colorIndex uint8) {
	if radius <= 0 {
		return
	}

	x := 0
	y := radius
	d := 3 - 2*radius

	for x <= y {
		SetMode4Pixel(cx+x, cy+y, colorIndex)
		SetMode4Pixel(cx-x, cy+y, colorIndex)
		SetMode4Pixel(cx+x, cy-y, colorIndex)
		SetMode4Pixel(cx-x, cy-y, colorIndex)
		SetMode4Pixel(cx+y, cy+x, colorIndex)
		SetMode4Pixel(cx-y, cy+x, colorIndex)
		SetMode4Pixel(cx+y, cy-x, colorIndex)
		SetMode4Pixel(cx-y, cy-x, colorIndex)

		if d < 0 {
			d = d + 4*x + 6
		} else {
			d = d + 4*(x-y) + 10
			y--
		}
		x++
	}
}

// DrawPixelMode4 Mode 4でピクセルを描画
func DrawPixelMode4(x, y int, colorIndex uint8) {
	SetMode4Pixel(x, y, colorIndex)
}
