package math

// Vec3 3次元ベクトル（固定小数点）
type Vec3 struct {
	X, Y, Z Fixed
}

// NewVec3 整数からベクトルを作成
func NewVec3(x, y, z int32) Vec3 {
	return Vec3{
		X: NewFixed(x),
		Y: NewFixed(y),
		Z: NewFixed(z),
	}
}

// NewVec3Fixed 固定小数点数からベクトルを作成
func NewVec3Fixed(x, y, z Fixed) Vec3 {
	return Vec3{X: x, Y: y, Z: z}
}

// NewVec3Float float64からベクトルを作成
func NewVec3Float(x, y, z float64) Vec3 {
	return Vec3{
		X: NewFixedFloat(x),
		Y: NewFixedFloat(y),
		Z: NewFixedFloat(z),
	}
}

// Add ベクトル加算
func (v Vec3) Add(other Vec3) Vec3 {
	return Vec3{
		X: v.X + other.X,
		Y: v.Y + other.Y,
		Z: v.Z + other.Z,
	}
}

// Sub ベクトル減算
func (v Vec3) Sub(other Vec3) Vec3 {
	return Vec3{
		X: v.X - other.X,
		Y: v.Y - other.Y,
		Z: v.Z - other.Z,
	}
}

// Mul スカラー倍
func (v Vec3) Mul(scalar Fixed) Vec3 {
	return Vec3{
		X: v.X.Mul(scalar),
		Y: v.Y.Mul(scalar),
		Z: v.Z.Mul(scalar),
	}
}

// Div スカラー除算
func (v Vec3) Div(scalar Fixed) Vec3 {
	if scalar == 0 {
		return v
	}
	return Vec3{
		X: v.X.Div(scalar),
		Y: v.Y.Div(scalar),
		Z: v.Z.Div(scalar),
	}
}

// Dot 内積
func (v Vec3) Dot(other Vec3) Fixed {
	return v.X.Mul(other.X).Add(v.Y.Mul(other.Y)).Add(v.Z.Mul(other.Z))
}

// Cross 外積
func (v Vec3) Cross(other Vec3) Vec3 {
	return Vec3{
		X: v.Y.Mul(other.Z).Sub(v.Z.Mul(other.Y)),
		Y: v.Z.Mul(other.X).Sub(v.X.Mul(other.Z)),
		Z: v.X.Mul(other.Y).Sub(v.Y.Mul(other.X)),
	}
}

// LengthSq ベクトルの長さの二乗
func (v Vec3) LengthSq() Fixed {
	return v.Dot(v)
}

// Length ベクトルの長さ（近似）
func (v Vec3) Length() Fixed {
	return FixedSqrt(v.LengthSq())
}

// Normalize 正規化（長さ1にする）
func (v Vec3) Normalize() Vec3 {
	length := v.Length()
	if length == 0 {
		return v
	}
	return v.Div(length)
}

// Neg ベクトルの符号反転
func (v Vec3) Neg() Vec3 {
	return Vec3{
		X: -v.X,
		Y: -v.Y,
		Z: -v.Z,
	}
}

// ToInt 整数座標に変換
func (v Vec3) ToInt() (x, y, z int32) {
	return v.X.ToInt(), v.Y.ToInt(), v.Z.ToInt()
}

// Distance3 2点間の距離（3D）
func Distance3(a, b Vec3) Fixed {
	return a.Sub(b).Length()
}

// DistanceSq3 2点間の距離の二乗（3D）
func DistanceSq3(a, b Vec3) Fixed {
	return a.Sub(b).LengthSq()
}

// Lerp3 線形補間（3D）
func Lerp3(a, b Vec3, t Fixed) Vec3 {
	return Vec3{
		X: Lerp(a.X, b.X, t),
		Y: Lerp(a.Y, b.Y, t),
		Z: Lerp(a.Z, b.Z, t),
	}
}
