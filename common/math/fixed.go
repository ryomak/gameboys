package math

// Fixed 固定小数点数（16.16形式）
// 上位16bitが整数部、下位16bitが小数部
type Fixed int32

const (
	FixedShift = 16                // 小数点位置
	FixedOne   = Fixed(1 << FixedShift) // 1.0を表す値
	FixedHalf  = Fixed(1 << (FixedShift - 1)) // 0.5を表す値
)

// NewFixed 整数から固定小数点数に変換
func NewFixed(x int32) Fixed {
	return Fixed(x << FixedShift)
}

// NewFixedFloat float64から固定小数点数に変換（近似）
func NewFixedFloat(x float64) Fixed {
	return Fixed(x * float64(FixedOne))
}

// ToInt 固定小数点数から整数に変換（小数部切り捨て）
func (f Fixed) ToInt() int32 {
	return int32(f >> FixedShift)
}

// ToFloat 固定小数点数からfloat64に変換
func (f Fixed) ToFloat() float64 {
	return float64(f) / float64(FixedOne)
}

// Mul 乗算
func (f Fixed) Mul(other Fixed) Fixed {
	return Fixed((int64(f) * int64(other)) >> FixedShift)
}

// Div 除算
func (f Fixed) Div(other Fixed) Fixed {
	if other == 0 {
		return 0 // ゼロ除算回避
	}
	return Fixed((int64(f) << FixedShift) / int64(other))
}

// Add 加算
func (f Fixed) Add(other Fixed) Fixed {
	return f + other
}

// Sub 減算
func (f Fixed) Sub(other Fixed) Fixed {
	return f - other
}

// Neg 符号反転
func (f Fixed) Neg() Fixed {
	return -f
}

// Abs 絶対値
func (f Fixed) Abs() Fixed {
	if f < 0 {
		return -f
	}
	return f
}

// Min 最小値
func (f Fixed) Min(other Fixed) Fixed {
	if f < other {
		return f
	}
	return other
}

// Max 最大値
func (f Fixed) Max(other Fixed) Fixed {
	if f > other {
		return f
	}
	return other
}

// Clamp 値を範囲内に制限
func (f Fixed) Clamp(min, max Fixed) Fixed {
	if f < min {
		return min
	}
	if f > max {
		return max
	}
	return f
}

// Round 四捨五入
func (f Fixed) Round() int32 {
	return int32((f + FixedHalf) >> FixedShift)
}

// Floor 切り捨て
func (f Fixed) Floor() int32 {
	return int32(f >> FixedShift)
}

// Ceil 切り上げ
func (f Fixed) Ceil() int32 {
	if (f & (FixedOne - 1)) == 0 {
		return int32(f >> FixedShift)
	}
	return int32((f >> FixedShift) + 1)
}

// Frac 小数部を取得
func (f Fixed) Frac() Fixed {
	return f & (FixedOne - 1)
}

// Lerp 線形補間 (t: 0.0-1.0)
func Lerp(a, b, t Fixed) Fixed {
	// a + (b - a) * t
	return a.Add(b.Sub(a).Mul(t))
}

// FixedMin 2つの固定小数点数の最小値
func FixedMin(a, b Fixed) Fixed {
	if a < b {
		return a
	}
	return b
}

// FixedMax 2つの固定小数点数の最大値
func FixedMax(a, b Fixed) Fixed {
	if a > b {
		return a
	}
	return b
}

// FixedAbs 固定小数点数の絶対値
func FixedAbs(f Fixed) Fixed {
	if f < 0 {
		return -f
	}
	return f
}
