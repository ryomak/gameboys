package math

// FixedSqrt 固定小数点数の平方根（ニュートン法による近似）
func FixedSqrt(x Fixed) Fixed {
	if x <= 0 {
		return 0
	}

	// 初期値の推定
	result := x
	if x >= FixedOne {
		result = x >> 1
	}

	// ニュートン法で反復計算（8回）
	for i := 0; i < 8; i++ {
		if result == 0 {
			break
		}
		// result = (result + x/result) / 2
		result = (result + x.Div(result)) >> 1
	}

	return result
}

// IntSqrt 整数の平方根（近似）
func IntSqrt(x int32) int32 {
	if x <= 0 {
		return 0
	}

	result := x
	if x >= 2 {
		result = x >> 1
	}

	for i := 0; i < 8; i++ {
		if result == 0 {
			break
		}
		result = (result + x/result) >> 1
	}

	return result
}
