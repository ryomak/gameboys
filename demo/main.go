package main

import (
	"github.com/ryomak/gameboys/common/gba/display"
	"github.com/ryomak/gameboys/common/gba/graphics"
	"github.com/ryomak/gameboys/common/gba/input"
	"github.com/ryomak/gameboys/common/math"
)

const (
	ballRadius = 8
	ballSpeed  = 2
)

// Ball プレイヤーが操作するボール
type Ball struct {
	x, y     int32
	prevX, prevY int32  // 前フレームの位置
	vx, vy   int32
	color    uint16
}

// NewBall ボールを作成
func NewBall(x, y int32, color uint16) *Ball {
	return &Ball{
		x:     x,
		y:     y,
		prevX: x,
		prevY: y,
		vx:    0,
		vy:    0,
		color: color,
	}
}

// Update ボールの状態を更新
func (b *Ball) Update(keys *input.KeyState) {
	// 前フレームの位置を保存
	b.prevX = b.x
	b.prevY = b.y

	// キー入力で速度を変更
	b.vx = 0
	b.vy = 0

	if keys.IsHeld(input.KeyUp) {
		b.vy = -ballSpeed
	}
	if keys.IsHeld(input.KeyDown) {
		b.vy = ballSpeed
	}
	if keys.IsHeld(input.KeyLeft) {
		b.vx = -ballSpeed
	}
	if keys.IsHeld(input.KeyRight) {
		b.vx = ballSpeed
	}

	// 位置を更新
	b.x += b.vx
	b.y += b.vy

	// 画面端で反転
	if b.x < ballRadius {
		b.x = ballRadius
	}
	if b.x > graphics.ScreenWidth-ballRadius {
		b.x = graphics.ScreenWidth - ballRadius
	}
	if b.y < ballRadius {
		b.y = ballRadius
	}
	if b.y > graphics.ScreenHeight-ballRadius {
		b.y = graphics.ScreenHeight - ballRadius
	}
}

// Erase 前フレームのボールを消す
func (b *Ball) Erase() {
	// 移動していない場合は消さない
	if b.x == b.prevX && b.y == b.prevY {
		return
	}

	// 前の位置を黒で塗りつぶす
	graphics.FillCircle(int(b.prevX), int(b.prevY), ballRadius, graphics.ColorBlack)

	// その位置にあった星を再描画
	for i := 0; i < len(stars); i++ {
		dx := stars[i].x - b.prevX
		dy := stars[i].y - b.prevY
		dist := dx*dx + dy*dy
		// 星が消された範囲にある場合は再描画
		if dist <= int32(ballRadius*ballRadius) {
			graphics.DrawPixel(int(stars[i].x), int(stars[i].y), stars[i].color)
		}
	}
}

// Draw ボールを描画
func (b *Ball) Draw() {
	graphics.FillCircle(int(b.x), int(b.y), ballRadius, b.color)
}

// Star 背景の星
type Star struct {
	x, y  int32
	color uint16
}

var stars [50]Star

// InitStars 星を初期化
func InitStars() {
	math.SetSeed(12345)
	for i := 0; i < len(stars); i++ {
		stars[i] = Star{
			x:     math.RandInt(graphics.ScreenWidth),
			y:     math.RandInt(graphics.ScreenHeight),
			color: graphics.ColorWhite,
		}
	}
}

// DrawStars 星を描画
func DrawStars() {
	for i := 0; i < len(stars); i++ {
		graphics.DrawPixel(int(stars[i].x), int(stars[i].y), stars[i].color)
	}
}

// InitUI UIの初期描画
func InitUI() {
	// タイトルバー
	graphics.FillRect(0, 0, graphics.ScreenWidth, 10, graphics.ColorDarkGray)
}

// UpdateUI UIを更新（フレームカウンターのみ）
var lastDisplayFrame int32 = -1

func UpdateUI(frame uint32) {
	displayFrame := int32((frame / 60) % 10)

	// 前回と同じ場合は更新しない
	if displayFrame == lastDisplayFrame {
		return
	}

	// フレームカウンター領域をクリア
	graphics.FillRect(220, 3, 20, 4, graphics.ColorDarkGray)

	// 新しいフレームカウンターを描画
	for i := int32(0); i < displayFrame; i++ {
		graphics.FillRect(int(220+i*2), 3, 1, 4, graphics.ColorGreen)
	}

	lastDisplayFrame = displayFrame
}

func main() {
	// ディスプレイ初期化
	display.SetMode(display.Mode3)
	display.EnableLayers(display.EnableBG2)

	// 入力初期化
	keys := input.NewKeyState()

	// ボール初期化
	ball := NewBall(graphics.ScreenWidth/2, graphics.ScreenHeight/2, graphics.ColorRed)

	// 星を初期化
	InitStars()

	// フレームカウンター
	var frame uint32 = 0

	// 初回描画：画面全体を初期化
	graphics.ClearScreen(graphics.ColorBlack)
	DrawStars()
	InitUI()
	ball.Draw()

	// メインループ
	for {
		// VBlank待機
		display.WaitForVBlank()

		// 入力更新
		keys.Update()

		// ゲーム終了チェック（Start + Select）
		if keys.IsHeld(input.KeyStart) && keys.IsHeld(input.KeySelect) {
			break
		}

		// 更新処理
		ball.Update(keys)

		// 描画処理（差分描画のみ）
		// 前の位置のボールを消す
		ball.Erase()

		// 新しい位置にボールを描画
		ball.Draw()

		// UI更新（フレームカウンターのみ）
		UpdateUI(frame)

		// Aボタンで色変更
		if keys.IsPressed(input.KeyA) {
			colors := []uint16{
				graphics.ColorRed,
				graphics.ColorBlue,
				graphics.ColorGreen,
				graphics.ColorYellow,
				graphics.ColorMagenta,
				graphics.ColorCyan,
			}
			colorIndex := (frame / 10) % uint32(len(colors))
			ball.color = colors[colorIndex]
		}

		// Bボタンで星の色変更
		if keys.IsPressed(input.KeyB) {
			for i := 0; i < len(stars); i++ {
				stars[i].color = uint16(math.RandInt(32768))
				// 星を再描画
				graphics.DrawPixel(int(stars[i].x), int(stars[i].y), stars[i].color)
			}
		}

		frame++
	}
}
