package graphics

// RGB15 RGB値から15bitカラーに変換（各色5bit: 0-31）
func RGB15(r, g, b uint8) uint16 {
	return uint16(r&0x1F) | (uint16(g&0x1F) << 5) | (uint16(b&0x1F) << 10)
}

// RGB8to5 8bitカラー値（0-255）を5bitカラー値（0-31）に変換
func RGB8to5(value uint8) uint8 {
	return value >> 3
}

// RGB RGB値（0-255）から15bitカラーに変換
func RGB(r, g, b uint8) uint16 {
	return RGB15(RGB8to5(r), RGB8to5(g), RGB8to5(b))
}

// ExtractRGB 15bitカラーからRGB値を抽出（各色0-31）
func ExtractRGB(color uint16) (r, g, b uint8) {
	r = uint8(color & 0x1F)
	g = uint8((color >> 5) & 0x1F)
	b = uint8((color >> 10) & 0x1F)
	return
}

// 基本カラー定義（15bit RGB）
var (
	ColorBlack   = RGB15(0, 0, 0)
	ColorWhite   = RGB15(31, 31, 31)
	ColorRed     = RGB15(31, 0, 0)
	ColorGreen   = RGB15(0, 31, 0)
	ColorBlue    = RGB15(0, 0, 31)
	ColorYellow  = RGB15(31, 31, 0)
	ColorCyan    = RGB15(0, 31, 31)
	ColorMagenta = RGB15(31, 0, 31)
	ColorOrange  = RGB15(31, 15, 0)
	ColorPurple  = RGB15(15, 0, 31)
	ColorPink    = RGB15(31, 15, 20)
	ColorBrown   = RGB15(15, 7, 0)
	ColorGray    = RGB15(15, 15, 15)
	ColorDarkGray = RGB15(7, 7, 7)
	ColorLightGray = RGB15(23, 23, 23)
)
