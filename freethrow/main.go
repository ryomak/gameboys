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

	// 角度から速度ベクトルを計算（g.angleは0-255の角度）
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

	// ゴールを描画（ボールより先に）
	g.drawGoal()

	// ボールを描画
	g.drawBall()

	// プレイヤーの手を描画（一人称視点）
	if g.state != StateShooting {
		g.drawPlayerHands()
	}

	// UIを描画
	g.drawUI()
}

// drawPlayerHands プレイヤーの手を描画
func (g *Game) drawPlayerHands() {
	// 肌色
	skinColor := graphics.RGB(220, 180, 140)
	darkSkinColor := graphics.RGB(180, 140, 100)

	// 左手（画面左下）
	leftHandX := 20
	leftHandY := graphics.ScreenHeight - 30

	// 左腕
	graphics.FillRect(leftHandX, leftHandY, 25, 35, darkSkinColor)
	graphics.FillRect(leftHandX+2, leftHandY+2, 21, 31, skinColor)

	// 左手の指
	for i := 0; i < 4; i++ {
		fingerX := leftHandX + 5 + i*5
		fingerY := leftHandY + 30
		graphics.FillRect(fingerX, fingerY, 3, 8, skinColor)
	}

	// 右手（画面右下）- シュート準備の位置
	rightHandX := graphics.ScreenWidth - 50
	rightHandY := graphics.ScreenHeight - 40

	// 右腕
	graphics.FillRect(rightHandX, rightHandY, 30, 40, darkSkinColor)
	graphics.FillRect(rightHandX+2, rightHandY+2, 26, 36, skinColor)

	// 右手の指（開いた状態）
	for i := 0; i < 5; i++ {
		fingerX := rightHandX + 5 + i*5
		fingerY := rightHandY + 35
		graphics.FillRect(fingerX, fingerY, 3, 10, skinColor)
	}
}

// drawCourt コートを描画
func (g *Game) drawCourt() {
	// 背景（体育館の壁）- 上部は暗め
	graphics.FillRect(0, 0, graphics.ScreenWidth, 60, graphics.RGB(30, 30, 50))

	// 床（遠近感のある台形）
	// 木目調の床色
	floorColor := graphics.RGB(139, 90, 43)
	lightFloorColor := graphics.RGB(160, 110, 60)

	// 床を段階的に描画して遠近感を出す
	for y := 90; y < graphics.ScreenHeight; y++ {
		// 遠くほど幅が狭い
		ratio := float64(y-90) / float64(graphics.ScreenHeight-90)
		width := int(float64(graphics.ScreenWidth) * (0.3 + ratio*0.7))
		startX := (graphics.ScreenWidth - width) / 2

		// 交互に色を変えて木目風に
		color := floorColor
		if (y/4)%2 == 0 {
			color = lightFloorColor
		}

		graphics.DrawLine(startX, y, startX+width, y, color)
	}

	// フリースローライン（白線）
	graphics.DrawLine(80, 145, 160, 145, graphics.ColorWhite)
	graphics.DrawLine(80, 146, 160, 146, graphics.ColorWhite)

	// ペイントエリアの線
	graphics.DrawLine(60, 140, 60, 155, graphics.ColorWhite)
	graphics.DrawLine(180, 140, 180, 155, graphics.ColorWhite)
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

	radius := int(size / 2)

	// バスケットボールを描画（グラデーションで立体感）
	// 外側から内側に向かって描画
	ballColor := graphics.RGB(255, 120, 0) // オレンジ
	darkBallColor := graphics.RGB(200, 80, 0) // 暗いオレンジ
	lightBallColor := graphics.RGB(255, 160, 60) // 明るいオレンジ

	// 影の部分（下側）
	for r := radius; r >= radius*2/3; r-- {
		graphics.DrawCircle(int(result.ScreenX), int(result.ScreenY+1), r, darkBallColor)
	}

	// メインの色
	graphics.FillCircle(int(result.ScreenX), int(result.ScreenY), radius, ballColor)

	// ハイライト（上側）
	highlightRadius := radius / 3
	if highlightRadius > 0 {
		graphics.FillCircle(
			int(result.ScreenX-int32(radius/4)),
			int(result.ScreenY-int32(radius/4)),
			highlightRadius,
			lightBallColor,
		)
	}

	// バスケットボールの線（黒いライン）
	if radius >= 4 {
		// 縦線
		graphics.DrawLine(
			int(result.ScreenX),
			int(result.ScreenY-int32(radius)),
			int(result.ScreenX),
			int(result.ScreenY+int32(radius)),
			graphics.ColorBlack,
		)
		// 横線
		graphics.DrawLine(
			int(result.ScreenX-int32(radius)),
			int(result.ScreenY),
			int(result.ScreenX+int32(radius)),
			int(result.ScreenY),
			graphics.ColorBlack,
		)
		// 斜め線
		offset := int32(radius * 7 / 10)
		graphics.DrawLine(
			int(result.ScreenX-offset),
			int(result.ScreenY-offset),
			int(result.ScreenX+offset),
			int(result.ScreenY+offset),
			graphics.ColorBlack,
		)
	}
}

// drawGoal ゴールを描画
func (g *Game) drawGoal() {
	baseDepth := math.NewFixed(300)

	// バックボード（背板）を描画
	backboardPos := math.NewVec3Fixed(0, g.goal.pos.Y.Add(math.NewFixed(20)), g.goal.pos.Z.Add(math.NewFixed(20)))
	backboardResult := math.ProjectSimple(backboardPos, graphics.ScreenWidth, graphics.ScreenHeight, baseDepth)

	if backboardResult.Visible {
		// バックボードのサイズ（遠近感を考慮）
		boardWidth := backboardResult.Scale.Mul(math.NewFixed(60)).ToInt()
		boardHeight := backboardResult.Scale.Mul(math.NewFixed(45)).ToInt()

		// バックボード（半透明の白）
		backboardColor := graphics.RGB(200, 200, 200)
		graphics.FillRect(
			int(backboardResult.ScreenX-boardWidth/2),
			int(backboardResult.ScreenY-boardHeight/2),
			int(boardWidth),
			int(boardHeight),
			backboardColor,
		)

		// バックボードの枠（赤）
		graphics.DrawRect(
			int(backboardResult.ScreenX-boardWidth/2),
			int(backboardResult.ScreenY-boardHeight/2),
			int(boardWidth),
			int(boardHeight),
			graphics.ColorRed,
		)

		// 四角いターゲット（内側の四角）
		targetSize := boardWidth / 3
		graphics.DrawRect(
			int(backboardResult.ScreenX-targetSize/2),
			int(backboardResult.ScreenY-targetSize/4),
			int(targetSize),
			int(targetSize/2),
			graphics.ColorRed,
		)
	}

	// ゴール（リム）の描画
	result := math.ProjectSimple(g.goal.pos, graphics.ScreenWidth, graphics.ScreenHeight, baseDepth)

	if !result.Visible {
		return
	}

	// ゴールのサイズ
	goalSize := result.Scale.Mul(math.NewFixed(35)).ToInt()
	if goalSize < 8 {
		goalSize = 8
	}

	// リム（楕円で立体感）- オレンジ色
	rimColor := graphics.RGB(255, 100, 0)

	// リムの外側
	graphics.DrawCircle(int(result.ScreenX), int(result.ScreenY), int(goalSize/2+1), rimColor)
	graphics.DrawCircle(int(result.ScreenX), int(result.ScreenY), int(goalSize/2), rimColor)

	// ネットを描画（格子状）
	netDepth := g.goal.pos.Z.Add(math.NewFixed(20))
	for i := int32(-2); i <= 2; i++ {
		netX := g.goal.pos.X.Add(math.NewFixed(i * 8))
		netPos := math.NewVec3Fixed(netX, g.goal.pos.Y.Sub(math.NewFixed(30)), netDepth)
		netResult := math.ProjectSimple(netPos, graphics.ScreenWidth, graphics.ScreenHeight, baseDepth)

		if netResult.Visible {
			// 縦のネット線
			graphics.DrawLine(
				int(result.ScreenX+i*goalSize/5),
				int(result.ScreenY),
				int(netResult.ScreenX),
				int(netResult.ScreenY),
				graphics.ColorWhite,
			)
		}
	}

	// 横のネット線
	for i := int32(0); i < 3; i++ {
		offsetY := (i + 1) * goalSize / 4
		graphics.DrawCircle(
			int(result.ScreenX),
			int(result.ScreenY+offsetY),
			int(goalSize/2-(goalSize/8)*i),
			graphics.ColorWhite,
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
	// スコアボード背景
	graphics.FillRect(5, 5, 85, 20, graphics.RGB(40, 40, 60))
	graphics.DrawRect(5, 5, 85, 20, graphics.ColorWhite)

	// 成功数（緑の丸）
	for i := int32(0); i < g.score && i < 10; i++ {
		graphics.FillCircle(int(10+i*8), 12, 3, graphics.ColorGreen)
	}

	// 試投数の枠（グレー）
	for i := int32(0); i < 10; i++ {
		color := graphics.RGB(80, 80, 80)
		if i < g.attempts {
			color = graphics.ColorGray
		}
		graphics.DrawCircle(int(10+i*8), 12, 3, color)
	}

	// 連続成功数の表示
	if g.consecutiveHits > 0 {
		graphics.FillRect(95, 5, 40, 10, graphics.RGB(255, 200, 0))
		graphics.DrawRect(95, 5, 40, 10, graphics.ColorYellow)
		// "STREAK" 表示（簡易版）
		for i := int32(0); i < g.consecutiveHits && i < 5; i++ {
			graphics.FillRect(int(98+i*7), 8, 4, 4, graphics.ColorRed)
		}
	}
}

// drawReadyUI 待機状態のUI
func (g *Game) drawReadyUI() {
	// "READY - Press A to Start" メッセージ
	msgWidth := 120
	msgHeight := 30
	msgX := (graphics.ScreenWidth - msgWidth) / 2
	msgY := 70

	// 背景（半透明風）
	graphics.FillRect(msgX, msgY, msgWidth, msgHeight, graphics.RGB(0, 0, 100))
	graphics.DrawRect(msgX, msgY, msgWidth, msgHeight, graphics.ColorWhite)
	graphics.DrawRect(msgX+1, msgY+1, msgWidth-2, msgHeight-2, graphics.ColorCyan)

	// "READY" 文字（ブロック表示）
	graphics.FillRect(msgX+15, msgY+8, 30, 6, graphics.ColorYellow)
	graphics.FillRect(msgX+50, msgY+8, 25, 6, graphics.ColorYellow)
	graphics.FillRect(msgX+80, msgY+8, 25, 6, graphics.ColorYellow)

	// "Press A" 点滅風（フレームで切り替え）
	graphics.FillRect(msgX+35, msgY+18, 50, 8, graphics.ColorGreen)
}

// drawPowerGauge パワーゲージを描画
func (g *Game) drawPowerGauge() {
	// ゲージの枠（立体感）
	gaugeX := 10
	gaugeY := 50
	gaugeWidth := 20
	gaugeHeight := 100

	// 外枠（影）
	graphics.FillRect(gaugeX+2, gaugeY+2, gaugeWidth, gaugeHeight, graphics.RGB(30, 30, 30))
	// メイン枠
	graphics.FillRect(gaugeX, gaugeY, gaugeWidth, gaugeHeight, graphics.RGB(60, 60, 60))
	graphics.DrawRect(gaugeX, gaugeY, gaugeWidth, gaugeHeight, graphics.ColorWhite)

	// 目盛り
	for i := 0; i <= 4; i++ {
		markY := gaugeY + (gaugeHeight * i / 4)
		graphics.DrawLine(gaugeX, markY, gaugeX+4, markY, graphics.ColorWhite)
		graphics.DrawLine(gaugeX+gaugeWidth-4, markY, gaugeX+gaugeWidth, markY, graphics.ColorWhite)
	}

	// ゲージの中身（グラデーション）
	fillHeight := (gaugeHeight * int(g.powerGauge.power)) / MaxPower
	if fillHeight > 0 {
		for i := 0; i < fillHeight; i++ {
			// 下から上に向かって色が変わる（緑→黄→赤）
			ratio := float64(i) / float64(gaugeHeight)
			var color uint16
			if ratio < 0.33 {
				color = graphics.RGB(0, 255, 0) // 緑
			} else if ratio < 0.66 {
				color = graphics.RGB(255, 255, 0) // 黄
			} else {
				color = graphics.RGB(255, 0, 0) // 赤
			}

			graphics.DrawLine(
				gaugeX+2,
				gaugeY+gaugeHeight-i,
				gaugeX+gaugeWidth-2,
				gaugeY+gaugeHeight-i,
				color,
			)
		}
	}

	// "POWER" ラベル
	graphics.FillRect(gaugeX-2, gaugeY-12, 24, 8, graphics.RGB(40, 40, 60))
	graphics.DrawRect(gaugeX-2, gaugeY-12, 24, 8, graphics.ColorWhite)
}

// drawAngleIndicator 角度インジケーターを描画
func (g *Game) drawAngleIndicator() {
	// パワーゲージも表示
	g.drawPowerGauge()

	// 角度計の背景
	centerX := graphics.ScreenWidth - 40
	centerY := 100
	radius := int32(35)

	// 背景円
	graphics.FillCircle(centerX, centerY, int(radius+5), graphics.RGB(40, 40, 60))
	graphics.DrawCircle(centerX, centerY, int(radius+5), graphics.ColorWhite)

	// 角度の範囲を示す弧（30-80度）
	for a := int32(MinAngle); a <= MaxAngle; a += 2 {
		angle := math.DegToAngle(a)
		x := centerX + int(math.Cos(angle).Mul(math.NewFixed(radius)).ToInt())
		y := centerY - int(math.Sin(angle).Mul(math.NewFixed(radius)).ToInt())
		graphics.DrawPixel(x, y, graphics.ColorGreen)
	}

	// 現在の角度を示す線
	angleDeg := math.AngleToDeg(g.angle)
	endX := centerX + int(math.Cos(g.angle).Mul(math.NewFixed(radius-5)).ToInt())
	endY := centerY - int(math.Sin(g.angle).Mul(math.NewFixed(radius-5)).ToInt())

	// 角度の針（太め）
	graphics.DrawLine(centerX, centerY, endX, endY, graphics.ColorYellow)
	graphics.DrawLine(centerX+1, centerY, endX+1, endY, graphics.ColorYellow)
	graphics.DrawLine(centerX, centerY+1, endX, endY+1, graphics.ColorYellow)

	// 中心点
	graphics.FillCircle(centerX, centerY, 3, graphics.ColorRed)

	// 最適角度（45度）を表示
	optimalAngle := math.DegToAngle(45)
	optX := centerX + int(math.Cos(optimalAngle).Mul(math.NewFixed(radius)).ToInt())
	optY := centerY - int(math.Sin(optimalAngle).Mul(math.NewFixed(radius)).ToInt())
	graphics.FillCircle(optX, optY, 2, graphics.RGB(0, 255, 0))

	// "ANGLE" ラベル
	graphics.FillRect(centerX-20, centerY-50, 40, 8, graphics.RGB(40, 40, 60))
	graphics.DrawRect(centerX-20, centerY-50, 40, 8, graphics.ColorWhite)

	// 角度の数値表示（簡易版）
	digitX := centerX - 10
	digitY := centerY + 15
	graphics.FillRect(digitX, digitY, 20, 10, graphics.ColorBlack)
	graphics.DrawRect(digitX, digitY, 20, 10, graphics.ColorWhite)

	// 角度の10の位と1の位を簡易表示
	tens := angleDeg / 10
	ones := angleDeg % 10
	for i := int32(0); i < tens && i < 9; i++ {
		graphics.FillRect(digitX+2, digitY+2+int(i), 6, 1, graphics.ColorYellow)
	}
	for i := int32(0); i < ones && i < 9; i++ {
		graphics.FillRect(digitX+12, digitY+2+int(i), 6, 1, graphics.ColorYellow)
	}
}

// drawResultUI 結果表示
func (g *Game) drawResultUI() {
	msgWidth := 140
	msgHeight := 50
	msgX := (graphics.ScreenWidth - msgWidth) / 2
	msgY := 60

	// 最後のシュートが成功したか判定（直前の状態から）
	lastSuccess := g.score > 0 && g.consecutiveHits > 0

	if lastSuccess {
		// 成功！
		// 背景（緑）
		graphics.FillRect(msgX, msgY, msgWidth, msgHeight, graphics.RGB(0, 150, 0))
		graphics.DrawRect(msgX, msgY, msgWidth, msgHeight, graphics.ColorWhite)
		graphics.DrawRect(msgX+2, msgY+2, msgWidth-4, msgHeight-4, graphics.ColorYellow)

		// "GOOD!" 文字風
		graphics.FillRect(msgX+20, msgY+10, 100, 15, graphics.ColorYellow)
		graphics.FillRect(msgX+25, msgY+12, 90, 11, graphics.RGB(0, 200, 0))

		// 星（装飾）
		for i := 0; i < 5; i++ {
			starX := msgX + 30 + i*20
			starY := msgY + 30
			graphics.FillRect(starX-2, starY, 5, 1, graphics.ColorYellow)
			graphics.FillRect(starX, starY-2, 1, 5, graphics.ColorYellow)
		}

		// 連続成功ボーナス表示
		if g.consecutiveHits >= 3 {
			graphics.FillRect(msgX+30, msgY+35, 80, 10, graphics.RGB(255, 200, 0))
			graphics.DrawRect(msgX+30, msgY+35, 80, 10, graphics.ColorRed)
		}
	} else {
		// 失敗...
		// 背景（赤）
		graphics.FillRect(msgX, msgY, msgWidth, msgHeight, graphics.RGB(150, 0, 0))
		graphics.DrawRect(msgX, msgY, msgWidth, msgHeight, graphics.ColorWhite)
		graphics.DrawRect(msgX+2, msgY+2, msgWidth-4, msgHeight-4, graphics.RGB(200, 100, 0))

		// "MISS" 文字風
		graphics.FillRect(msgX+20, msgY+10, 100, 15, graphics.RGB(200, 0, 0))
		graphics.FillRect(msgX+25, msgY+12, 90, 11, graphics.RGB(100, 0, 0))

		// X マーク
		for i := 0; i < 20; i++ {
			graphics.DrawPixel(msgX+40+i, msgY+25+i, graphics.ColorRed)
			graphics.DrawPixel(msgX+60-i, msgY+25+i, graphics.ColorRed)
		}
	}

	// "Press A to Continue"
	graphics.FillRect(msgX+20, msgY+msgHeight+10, 100, 8, graphics.ColorBlue)
	graphics.DrawRect(msgX+20, msgY+msgHeight+10, 100, 8, graphics.ColorWhite)
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
