package util

// Min 2つの整数の最小値
func Min(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

// Max 2つの整数の最大値
func Max(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

// Clamp 値を範囲内に制限
func Clamp(value, min, max int32) int32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Abs 絶対値
func Abs(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

// Sign 符号（-1, 0, 1）
func Sign(x int32) int32 {
	if x > 0 {
		return 1
	}
	if x < 0 {
		return -1
	}
	return 0
}

// Lerp 線形補間（t: 0-100）
func Lerp(a, b, t int32) int32 {
	return a + (b-a)*t/100
}
