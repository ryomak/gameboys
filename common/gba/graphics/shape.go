package graphics

// DrawLine 線を描画（ブレゼンハムのアルゴリズム）
func DrawLine(x0, y0, x1, y1 int, color uint16) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx := -1
	if x0 < x1 {
		sx = 1
	}
	sy := -1
	if y0 < y1 {
		sy = 1
	}
	err := dx - dy

	for {
		DrawPixel(x0, y0, color)
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

// DrawRect 矩形の枠を描画
func DrawRect(x, y, width, height int, color uint16) {
	if width <= 0 || height <= 0 {
		return
	}

	// 上下の線
	DrawHLine(x, y, width, color)
	if height > 1 {
		DrawHLine(x, y+height-1, width, color)
	}

	// 左右の線（角は既に描画済み）
	if height > 2 {
		DrawVLine(x, y+1, height-2, color)
		if width > 1 {
			DrawVLine(x+width-1, y+1, height-2, color)
		}
	}
}

// DrawCircle 円を描画（中点円描画アルゴリズム）
func DrawCircle(cx, cy, radius int, color uint16) {
	if radius <= 0 {
		return
	}

	x := radius
	y := 0
	err := 0

	for x >= y {
		DrawPixel(cx+x, cy+y, color)
		DrawPixel(cx+y, cy+x, color)
		DrawPixel(cx-y, cy+x, color)
		DrawPixel(cx-x, cy+y, color)
		DrawPixel(cx-x, cy-y, color)
		DrawPixel(cx-y, cy-x, color)
		DrawPixel(cx+y, cy-x, color)
		DrawPixel(cx+x, cy-y, color)

		if err <= 0 {
			y++
			err += 2*y + 1
		}
		if err > 0 {
			x--
			err -= 2*x + 1
		}
	}
}

// FillCircle 塗りつぶした円を描画
func FillCircle(cx, cy, radius int, color uint16) {
	if radius <= 0 {
		return
	}

	x := radius
	y := 0
	err := 0

	for x >= y {
		DrawHLine(cx-x, cy+y, 2*x+1, color)
		DrawHLine(cx-x, cy-y, 2*x+1, color)
		DrawHLine(cx-y, cy+x, 2*y+1, color)
		DrawHLine(cx-y, cy-x, 2*y+1, color)

		if err <= 0 {
			y++
			err += 2*y + 1
		}
		if err > 0 {
			x--
			err -= 2*x + 1
		}
	}
}

// DrawTriangle 三角形の枠を描画
func DrawTriangle(x0, y0, x1, y1, x2, y2 int, color uint16) {
	DrawLine(x0, y0, x1, y1, color)
	DrawLine(x1, y1, x2, y2, color)
	DrawLine(x2, y2, x0, y0, color)
}

// abs 絶対値を返す
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// min 最小値を返す
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max 最大値を返す
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
