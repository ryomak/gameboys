package math

// Camera 3Dカメラ
type Camera struct {
	Position Vec3  // カメラの位置
	Target   Vec3  // カメラが見ている点
	FOV      Fixed // 視野角（固定小数点、ラジアン相当）
	Near     Fixed // ニアクリップ
	Far      Fixed // ファークリップ
}

// NewCamera デフォルトカメラを作成
func NewCamera() Camera {
	return Camera{
		Position: NewVec3(0, 0, -5),
		Target:   NewVec3(0, 0, 0),
		FOV:      NewFixedFloat(60.0), // 視野角60度相当
		Near:     NewFixed(1),
		Far:      NewFixed(100),
	}
}

// ProjectionResult 射影変換の結果
type ProjectionResult struct {
	ScreenX int32 // 画面X座標
	ScreenY int32 // 画面Y座標
	Scale   Fixed // スケール係数（スプライトサイズ調整用）
	Visible bool  // 画面内に表示されるか
}

// Project 3D座標を2D画面座標に変換（簡易版）
func (c *Camera) Project(worldPos Vec3, screenWidth, screenHeight int32) ProjectionResult {
	// カメラ座標系に変換
	relPos := worldPos.Sub(c.Position)

	// カメラの前方ベクトル
	forward := c.Target.Sub(c.Position).Normalize()

	// Z軸での距離（深度）
	depth := relPos.Dot(forward)

	// 視錐台の外なら非表示
	if depth < c.Near || depth > c.Far {
		return ProjectionResult{Visible: false}
	}

	// 透視投影
	// スケール = 1 / (1 + depth * factor)
	// より簡易的に: scale = near / depth
	scale := c.Near.Div(depth)

	// 画面中心からのオフセット
	centerX := screenWidth / 2
	centerY := screenHeight / 2

	// X, Y座標を画面座標に変換
	screenX := relPos.X.Mul(scale).ToInt() + centerX
	screenY := centerY - relPos.Y.Mul(scale).ToInt() // Y軸は上が正

	// 画面外判定（マージン付き）
	margin := int32(32)
	visible := screenX >= -margin && screenX < screenWidth+margin &&
		screenY >= -margin && screenY < screenHeight+margin

	return ProjectionResult{
		ScreenX: screenX,
		ScreenY: screenY,
		Scale:   scale,
		Visible: visible,
	}
}

// ProjectSimple 簡易射影変換（カメラなし）
// Z座標から直接スケールを計算
func ProjectSimple(pos Vec3, screenWidth, screenHeight int32, baseDepth Fixed) ProjectionResult {
	// Z座標が大きいほど遠い（小さく表示）
	// scale = baseDepth / (baseDepth + z)
	if pos.Z < 0 {
		// カメラより後ろは非表示
		return ProjectionResult{Visible: false}
	}

	depth := baseDepth.Add(pos.Z)
	if depth <= 0 {
		depth = FixedOne
	}

	scale := baseDepth.Div(depth)

	// 画面中心からのオフセット
	centerX := screenWidth / 2
	centerY := screenHeight / 2

	screenX := pos.X.Mul(scale).ToInt() + centerX
	screenY := centerY - pos.Y.Mul(scale).ToInt()

	// 画面外判定
	margin := int32(32)
	visible := screenX >= -margin && screenX < screenWidth+margin &&
		screenY >= -margin && screenY < screenHeight+margin

	return ProjectionResult{
		ScreenX: screenX,
		ScreenY: screenY,
		Scale:   scale,
		Visible: visible,
	}
}

// CalculateSpriteScale スプライトのスケール係数を計算
// distance: カメラからの距離
// baseScale: 基準距離でのスケール
// 戻り値: 0-255のスケール値（GBAのアフィン変換用）
func CalculateSpriteScale(distance Fixed, baseDistance Fixed, baseScale int32) int32 {
	if distance <= 0 {
		return baseScale
	}

	// scale = baseScale * (baseDistance / distance)
	scale := NewFixed(baseScale).Mul(baseDistance.Div(distance))
	result := scale.ToInt()

	// 範囲制限（GBAのアフィン変換は8bitスケール）
	if result < 1 {
		return 1
	}
	if result > 512 {
		return 512
	}
	return result
}

// SortByDepth 深度ソート用の比較関数
// 遠い（Z値が大きい）ものを先に描画するため、大きい方を「小さい」とする
func CompareDepth(a, b Vec3) int {
	if a.Z > b.Z {
		return -1 // aが遠い（先に描画）
	} else if a.Z < b.Z {
		return 1 // bが遠い
	}
	return 0
}
