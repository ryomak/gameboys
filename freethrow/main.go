package main

import (
	"github.com/ryomak/gameboys/common/gba/display"
	"github.com/ryomak/gameboys/common/gba/graphics"
	"github.com/ryomak/gameboys/common/gba/input"
	"github.com/ryomak/gameboys/common/math"
)

// GameState ゲームの状態
type GameState int

const (
	StateReady      GameState = iota // 待機状態
	StatePowerGauge                  // パワーゲージ調整
	StateAngleAdjust                 // 角度調整
	StateShooting                    // シュート中
	StateResult                      // 結果表示
)

// Game ゲーム全体の管理
type Game struct {
	state         GameState
	ball          Ball
	goal          Goal
	powerGauge    PowerGauge
	angle         int32 // 角度（0-255、0度-360度）
	score         int32 // スコア
	attempts      int32 // 試投数
	consecutiveHits int32 // 連続成功数
}

// Ball バスケットボール
type Ball struct {
	pos       math.Vec3  // 3D位置（メートル単位、固定小数点）
	velocity  math.Vec3  // 速度ベクトル
	isFlying  bool       // 飛んでいるか
	radius    math.Fixed // ボールの半径（メートル）
}

// Goal バスケットゴール
type Goal struct {
	pos    math.Vec3  // ゴールの位置
	radius math.Fixed // ゴールの半径
}

// PowerGauge パワーゲージ
type PowerGauge struct {
	power      int32 // 0-100
	increasing bool  // 増加中か
	speed      int32 // 変化速度
}

// 物理定数
const (
	Gravity        = 980  // 重力加速度 (cm/s^2) 固定小数点前の値
	DeltaTime      = 16   // フレーム時間（ミリ秒、約60fps）
	GoalHeight     = 305  // ゴールの高さ（cm）
	GoalDistance   = 422  // フリースローラインからゴールまでの距離（cm）
	GoalRadius     = 23   // ゴールの半径（cm）
	BallRadius     = 12   // ボールの半径（cm）
	PlayerHeight   = 200  // プレイヤーの手の高さ（cm）
	MaxPower       = 100  // 最大パワー
	MaxAngle       = 80   // 最大角度（度）
	MinAngle       = 30   // 最小角度（度）
	AngleDefault   = 55   // デフォルト角度（度）
)

// NewGame ゲームを初期化
func NewGame() *Game {
	return &Game{
		state: StateReady,
		ball: Ball{
			pos:      math.NewVec3Fixed(0, math.NewFixed(PlayerHeight), 0),
			velocity: math.NewVec3(0, 0, 0),
			isFlying: false,
			radius:   math.NewFixed(BallRadius),
		},
		goal: Goal{
			pos:    math.NewVec3Fixed(0, math.NewFixed(GoalHeight), math.NewFixed(GoalDistance)),
			radius: math.NewFixed(GoalRadius),
		},
		powerGauge: PowerGauge{
			power:      0,
			increasing: true,
			speed:      3,
		},
		angle:    math.DegToAngle(AngleDefault),
		score:    0,
		attempts: 0,
		consecutiveHits: 0,
	}
}

// Update ゲームの状態を更新
func (g *Game) Update(keys *input.KeyState) {
	switch g.state {
	case StateReady:
		g.updateReady(keys)
	case StatePowerGauge:
		g.updatePowerGauge(keys)
	case StateAngleAdjust:
		g.updateAngleAdjust(keys)
	case StateShooting:
		g.updateShooting()
	case StateResult:
		g.updateResult(keys)
	}
}

// updateReady 待機状態の更新
func (g *Game) updateReady(keys *input.KeyState) {
	if keys.IsPressed(input.KeyA) {
		g.state = StatePowerGauge
		g.powerGauge.power = 0
		g.powerGauge.increasing = true
	}
}

// updatePowerGauge パワーゲージの更新
func (g *Game) updatePowerGauge(keys *input.KeyState) {
	// パワーゲージを増減
	if g.powerGauge.increasing {
		g.powerGauge.power += g.powerGauge.speed
		if g.powerGauge.power >= MaxPower {
			g.powerGauge.power = MaxPower
			g.powerGauge.increasing = false
		}
	} else {
		g.powerGauge.power -= g.powerGauge.speed
		if g.powerGauge.power <= 0 {
			g.powerGauge.power = 0
			g.powerGauge.increasing = true
		}
	}

	// Aボタンでパワー決定
	if keys.IsPressed(input.KeyA) {
		g.state = StateAngleAdjust
	}

	// Bボタンでキャンセル
	if keys.IsPressed(input.KeyB) {
		g.state = StateReady
	}
}

// updateAngleAdjust 角度調整の更新
func (g *Game) updateAngleAdjust(keys *input.KeyState) {
	// 上下キーで角度調整
	angleDeg := math.AngleToDeg(g.angle)

	if keys.IsHeld(input.KeyUp) {
		angleDeg++
		if angleDeg > MaxAngle {
			angleDeg = MaxAngle
		}
	}
	if keys.IsHeld(input.KeyDown) {
		angleDeg--
		if angleDeg < MinAngle {
			angleDeg = MinAngle
		}
	}

	g.angle = math.DegToAngle(angleDeg)

	// Aボタンでシュート
	if keys.IsPressed(input.KeyA) {
		g.shoot()
		g.state = StateShooting
		g.attempts++
	}

	// Bボタンでキャンセル
	if keys.IsPressed(input.KeyB) {
		g.state = StatePowerGauge
	}
}

// shoot シュートを実行
func (g *Game) shoot() {
	// パワーから初速度を計算
	// power: 0-100 -> velocity: 500-1500 cm/s
	velocityMag := 500 + (g.powerGauge.power * 10)

	// 角度から速度ベクトルを計算
	angleDeg := math.AngleToDeg(g.angle)
	angleRad := math.NewFixedFloat(float64(angleDeg) * 0.0174533) // deg to rad

	// 簡易的な三角関数計算
	vz := math.NewFixed(velocityMag).Mul(math.Cos(g.angle))
	vy := math.NewFixed(velocityMag).Mul(math.Sin(g.angle))

	g.ball.velocity = math.NewVec3Fixed(0, vy, vz)
	g.ball.isFlying = true
	g.ball.pos = math.NewVec3Fixed(0, math.NewFixed(PlayerHeight), 0)
}

// updateShooting シュート中の更新
func (g *Game) updateShooting() {
	if !g.ball.isFlying {
		return
	}

	// 時間刻み（秒）
	dt := math.NewFixedFloat(float64(DeltaTime) / 1000.0)

	// 重力を適用
	gravity := math.NewFixed(Gravity)
	g.ball.velocity.Y = g.ball.velocity.Y.Sub(gravity.Mul(dt))

	// 位置を更新
	g.ball.pos = g.ball.pos.Add(g.ball.velocity.Mul(dt))

	// 地面に落ちたら終了
	if g.ball.pos.Y < 0 {
		g.ball.isFlying = false
		g.checkResult()
		g.state = StateResult
		return
	}

	// ゴールとの当たり判定
	if g.checkGoal() {
		g.ball.isFlying = false
		g.score++
		g.consecutiveHits++
		g.state = StateResult
	}
}

// checkGoal ゴールに入ったか判定
func (g *Game) checkGoal() bool {
	// ボールの中心がゴールの高さ付近にあるか
	goalY := g.goal.pos.Y
	ballY := g.ball.pos.Y

	// Y方向の許容範囲（ゴール通過の高さ）
	yDiff := ballY.Sub(goalY).Abs()
	if yDiff > math.NewFixed(50) { // 50cm以内
		return false
	}

	// Z方向の位置確認（ゴールの位置を通過しているか）
	if g.ball.pos.Z < g.goal.pos.Z || g.ball.pos.Z > g.goal.pos.Z.Add(math.NewFixed(50)) {
		return false
	}

	// XY平面での距離を計算
	dx := g.ball.pos.X.Sub(g.goal.pos.X)
	dy := g.ball.pos.Y.Sub(g.goal.pos.Y)
	distSq := dx.Mul(dx).Add(dy.Mul(dy))
	radiusSq := g.goal.radius.Mul(g.goal.radius)

	return distSq <= radiusSq
}

// checkResult 結果をチェック
func (g *Game) checkResult() {
	// ゴールに入らなかった場合
	g.consecutiveHits = 0
}

// updateResult 結果表示の更新
func (g *Game) updateResult(keys *input.KeyState) {
	if keys.IsPressed(input.KeyA) || keys.IsPressed(input.KeyB) {
		g.state = StateReady
	}
}

// Draw ゲームを描画
func (g *Game) Draw() {
	// 背景をクリア
	graphics.ClearScreen(graphics.ColorBlack)

	// コートを描画
	g.drawCourt()

	// ボールを描画
	g.drawBall()

	// ゴールを描画
	g.drawGoal()

	// UIを描画
	g.drawUI()
}

// drawCourt コートを描画
func (g *Game) drawCourt() {
	// 地面のライン
	graphics.DrawLine(0, 140, graphics.ScreenWidth, 140, graphics.ColorWhite)

	// フリースローライン
	graphics.DrawLine(60, 140, 180, 140, graphics.ColorRed)
}

// drawBall ボールを描画（3D→2D変換）
func (g *Game) drawBall() {
	// 簡易的な射影変換
	baseDepth := math.NewFixed(300)
	result := math.ProjectSimple(g.ball.pos, graphics.ScreenWidth, graphics.ScreenHeight, baseDepth)

	if !result.Visible {
		return
	}

	// スケールに応じたボールサイズ
	size := result.Scale.Mul(math.NewFixed(16)).ToInt()
	if size < 2 {
		size = 2
	}
	if size > 32 {
		size = 32
	}

	// ボールを描画（円）
	graphics.FillCircle(int(result.ScreenX), int(result.ScreenY), int(size/2), graphics.ColorOrange)
}

// drawGoal ゴールを描画
func (g *Game) drawGoal() {
	// ゴール（リム）の描画
	baseDepth := math.NewFixed(300)
	result := math.ProjectSimple(g.goal.pos, graphics.ScreenWidth, graphics.ScreenHeight, baseDepth)

	if !result.Visible {
		return
	}

	// ゴールの円
	goalSize := result.Scale.Mul(math.NewFixed(40)).ToInt()
	if goalSize < 10 {
		goalSize = 10
	}

	// リムを描画
	graphics.DrawCircle(int(result.ScreenX), int(result.ScreenY), int(goalSize/2), graphics.ColorWhite)

	// ネット（縦線）
	netDepth := g.goal.pos.Z.Add(math.NewFixed(30))
	netPos := math.NewVec3Fixed(g.goal.pos.X, g.goal.pos.Y.Sub(math.NewFixed(30)), netDepth)
	netResult := math.ProjectSimple(netPos, graphics.ScreenWidth, graphics.ScreenHeight, baseDepth)

	if netResult.Visible {
		graphics.DrawLine(
			int(result.ScreenX), int(result.ScreenY),
			int(netResult.ScreenX), int(netResult.ScreenY),
			graphics.ColorGray,
		)
	}
}

// drawUI UIを描画
func (g *Game) drawUI() {
	// スコア表示
	g.drawScore()

	// 状態に応じたUIを描画
	switch g.state {
	case StateReady:
		g.drawReadyUI()
	case StatePowerGauge:
		g.drawPowerGauge()
	case StateAngleAdjust:
		g.drawAngleIndicator()
	case StateResult:
		g.drawResultUI()
	}
}

// drawScore スコアを描画
func (g *Game) drawScore() {
	// 簡易的なスコア表示（バー）
	for i := int32(0); i < g.score && i < 10; i++ {
		graphics.FillRect(int(5+i*8), 5, 6, 6, graphics.ColorGreen)
	}

	// 試投数
	for i := int32(0); i < g.attempts && i < 10; i++ {
		graphics.FillRect(int(5+i*8), 15, 6, 4, graphics.ColorGray)
	}
}

// drawReadyUI 待機状態のUI
func (g *Game) drawReadyUI() {
	// "Press A" メッセージ（簡易的な表示）
	graphics.FillRect(100, 100, 40, 8, graphics.ColorBlue)
}

// drawPowerGauge パワーゲージを描画
func (g *Game) drawPowerGauge() {
	// ゲージの枠
	gaugeX := 10
	gaugeY := 60
	gaugeWidth := 10
	gaugeHeight := 80

	graphics.DrawRect(gaugeX, gaugeY, gaugeWidth, gaugeHeight, graphics.ColorWhite)

	// ゲージの中身
	fillHeight := (gaugeHeight * int(g.powerGauge.power)) / MaxPower
	if fillHeight > 0 {
		graphics.FillRect(gaugeX+1, gaugeY+gaugeHeight-fillHeight, gaugeWidth-2, fillHeight, graphics.ColorRed)
	}
}

// drawAngleIndicator 角度インジケーターを描画
func (g *Game) drawAngleIndicator() {
	// パワーゲージも表示
	g.drawPowerGauge()

	// 角度の弧（簡易版）
	centerX := graphics.ScreenWidth - 30
	centerY := 120

	angleDeg := math.AngleToDeg(g.angle)

	// 角度を示す線（簡易的）
	length := int32(20)
	endX := centerX + int(math.Cos(g.angle).Mul(math.NewFixed(length)).ToInt())
	endY := centerY - int(math.Sin(g.angle).Mul(math.NewFixed(length)).ToInt())

	graphics.DrawLine(centerX, centerY, endX, endY, graphics.ColorYellow)

	// 角度の値を表示（10度刻み）
	angleMarks := angleDeg / 10
	for i := int32(0); i < angleMarks; i++ {
		graphics.FillRect(centerX+25, 120-int(i*4), 8, 2, graphics.ColorYellow)
	}
}

// drawResultUI 結果表示
func (g *Game) drawResultUI() {
	// 成功/失敗の表示
	if g.consecutiveHits > 0 {
		// 成功
		graphics.FillRect(80, 80, 80, 20, graphics.ColorGreen)
	} else {
		// 失敗
		graphics.FillRect(80, 80, 80, 20, graphics.ColorRed)
	}
}

func main() {
	// ディスプレイ初期化
	display.SetMode(display.Mode3)
	display.EnableLayers(display.EnableBG2)

	// 入力初期化
	keys := input.NewKeyState()

	// ゲーム初期化
	game := NewGame()

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
		game.Update(keys)

		// 描画処理
		game.Draw()
	}
}
