package math

// Vec2 2次元ベクトル（固定小数点）
type Vec2 struct {
	X, Y Fixed
}

// NewVec2 整数からベクトルを作成
func NewVec2(x, y int32) Vec2 {
	return Vec2{
		X: NewFixed(x),
		Y: NewFixed(y),
	}
}

// NewVec2Fixed 固定小数点数からベクトルを作成
func NewVec2Fixed(x, y Fixed) Vec2 {
	return Vec2{X: x, Y: y}
}

// Add ベクトル加算
func (v Vec2) Add(other Vec2) Vec2 {
	return Vec2{
		X: v.X + other.X,
		Y: v.Y + other.Y,
	}
}

// Sub ベクトル減算
func (v Vec2) Sub(other Vec2) Vec2 {
	return Vec2{
		X: v.X - other.X,
		Y: v.Y - other.Y,
	}
}

// Mul スカラー倍
func (v Vec2) Mul(scalar Fixed) Vec2 {
	return Vec2{
		X: v.X.Mul(scalar),
		Y: v.Y.Mul(scalar),
	}
}

// Div スカラー除算
func (v Vec2) Div(scalar Fixed) Vec2 {
	if scalar == 0 {
		return v
	}
	return Vec2{
		X: v.X.Div(scalar),
		Y: v.Y.Div(scalar),
	}
}

// Dot 内積
func (v Vec2) Dot(other Vec2) Fixed {
	return v.X.Mul(other.X).Add(v.Y.Mul(other.Y))
}

// LengthSq ベクトルの長さの二乗
func (v Vec2) LengthSq() Fixed {
	return v.Dot(v)
}

// Length ベクトルの長さ（近似）
func (v Vec2) Length() Fixed {
	// 簡易的な長さ計算（完全な平方根ではなく近似）
	return FixedSqrt(v.LengthSq())
}

// Normalize 正規化（長さ1にする）
func (v Vec2) Normalize() Vec2 {
	length := v.Length()
	if length == 0 {
		return v
	}
	return v.Div(length)
}

// Neg ベクトルの符号反転
func (v Vec2) Neg() Vec2 {
	return Vec2{
		X: -v.X,
		Y: -v.Y,
	}
}

// ToInt 整数座標に変換
func (v Vec2) ToInt() (x, y int32) {
	return v.X.ToInt(), v.Y.ToInt()
}

// Distance 2点間の距離
func Distance(a, b Vec2) Fixed {
	return a.Sub(b).Length()
}

// DistanceSq 2点間の距離の二乗
func DistanceSq(a, b Vec2) Fixed {
	return a.Sub(b).LengthSq()
}

// Vec2Int 整数ベクトル
type Vec2Int struct {
	X, Y int32
}

// NewVec2Int 整数ベクトルを作成
func NewVec2Int(x, y int32) Vec2Int {
	return Vec2Int{X: x, Y: y}
}

// Add ベクトル加算
func (v Vec2Int) Add(other Vec2Int) Vec2Int {
	return Vec2Int{
		X: v.X + other.X,
		Y: v.Y + other.Y,
	}
}

// Sub ベクトル減算
func (v Vec2Int) Sub(other Vec2Int) Vec2Int {
	return Vec2Int{
		X: v.X - other.X,
		Y: v.Y - other.Y,
	}
}

// Mul スカラー倍
func (v Vec2Int) Mul(scalar int32) Vec2Int {
	return Vec2Int{
		X: v.X * scalar,
		Y: v.Y * scalar,
	}
}

// ToFixed 固定小数点ベクトルに変換
func (v Vec2Int) ToFixed() Vec2 {
	return NewVec2(v.X, v.Y)
}
