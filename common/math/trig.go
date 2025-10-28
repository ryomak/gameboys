package math

// 三角関数のルックアップテーブル
// GBAにはFPUがないため、ルックアップテーブルを使用して高速化

// 角度を256段階に分割（0-255が0-360度に対応）
const (
	AngleMax     = 256 // 1周を256段階に分割
	AngleQuarter = 64  // 90度
	AngleHalf    = 128 // 180度
)

// SinTable sin値のルックアップテーブル（固定小数点）
// 0-90度（0-64）の値を格納、他は対称性を利用
var SinTable = [65]Fixed{
	0, 402, 804, 1206, 1608, 2009, 2410, 2811, 3212, 3612, 4011, 4410, 4808, 5205, 5602, 5998,
	6393, 6786, 7179, 7571, 7962, 8351, 8739, 9126, 9512, 9896, 10278, 10659, 11039, 11417, 11793, 12167,
	12539, 12910, 13278, 13645, 14010, 14372, 14732, 15090, 15446, 15800, 16151, 16499, 16846, 17189, 17530, 17869,
	18204, 18537, 18868, 19195, 19519, 19841, 20159, 20475, 20787, 21096, 21403, 21706, 22005, 22301, 22594, 22884,
	23170,
}

// Sin 正弦関数（ルックアップテーブル使用）
// angle: 0-255 が 0-360度に対応
func Sin(angle int32) Fixed {
	// 角度を0-255の範囲に正規化
	angle = angle & (AngleMax - 1)

	if angle < AngleQuarter {
		// 0-90度
		return SinTable[angle]
	} else if angle < AngleHalf {
		// 90-180度
		return SinTable[AngleHalf-angle]
	} else if angle < AngleHalf+AngleQuarter {
		// 180-270度
		return -SinTable[angle-AngleHalf]
	} else {
		// 270-360度
		return -SinTable[AngleMax-angle]
	}
}

// Cos 余弦関数（sinを利用）
// angle: 0-255 が 0-360度に対応
func Cos(angle int32) Fixed {
	// cos(x) = sin(x + 90度)
	return Sin(angle + AngleQuarter)
}

// Tan 正接関数（近似）
// angle: 0-255 が 0-360度に対応
func Tan(angle int32) Fixed {
	cosVal := Cos(angle)
	if cosVal == 0 {
		// 無限大を防ぐ
		return FixedOne << 10 // 非常に大きな値
	}
	return Sin(angle).Div(cosVal)
}

// Atan2 2引数逆正接（近似版）
// y, xから角度（0-255）を返す
func Atan2(y, x Fixed) int32 {
	if x == 0 && y == 0 {
		return 0
	}

	// 簡易的な実装：まず象限を判定
	var angle int32
	absX := x.Abs()
	absY := y.Abs()

	// 小さい方を大きい方で割る
	if absX > absY {
		// x軸に近い
		ratio := absY.Div(absX)
		angle = atanApprox(ratio)
		if x < 0 {
			angle = AngleHalf - angle
		}
		if y < 0 {
			angle = -angle
		}
	} else {
		// y軸に近い
		ratio := absX.Div(absY)
		angle = AngleQuarter - atanApprox(ratio)
		if y < 0 {
			angle = AngleHalf + (AngleQuarter - angle)
		}
		if x < 0 {
			angle = AngleHalf - angle
		}
	}

	return angle & (AngleMax - 1)
}

// atanApprox 0-1の範囲でのatan近似（0-45度）
func atanApprox(x Fixed) int32 {
	// 簡易的な多項式近似
	// atan(x) ≈ x - x^3/3 + x^5/5（ただし0-1の範囲）
	// 結果を0-64（0-90度）にスケール

	// さらに簡易版：線形近似
	// x = 0 -> 0度, x = 1 -> 45度(64/4=16)
	result := x.Mul(NewFixed(AngleQuarter))
	return result.ToInt()
}

// DegToAngle 度数法からAngle形式に変換
func DegToAngle(degrees int32) int32 {
	// 360度 = 256
	return (degrees * AngleMax) / 360
}

// AngleToDeg Angle形式から度数法に変換
func AngleToDeg(angle int32) int32 {
	// 256 = 360度
	return (angle * 360) / AngleMax
}

// RadToAngle ラジアンからAngle形式に変換
func RadToAngle(radians Fixed) int32 {
	// 2π = 256
	// angle = rad * 256 / (2π) = rad * 256 / 6.28318... ≈ rad * 40.74
	factor := NewFixedFloat(40.74) // 256 / (2 * π)
	return radians.Mul(factor).ToInt()
}

// AngleToRad Angle形式からラジアンに変換
func AngleToRad(angle int32) Fixed {
	// angle * 2π / 256 = angle * 6.28318 / 256 ≈ angle * 0.0245
	factor := NewFixedFloat(0.024543693) // 2π / 256
	return NewFixed(angle).Mul(factor)
}
