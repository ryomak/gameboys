package math

// 線形合同法による乱数生成器
// GBAには標準ライブラリの乱数生成器が使えないため、独自実装

var seed uint32 = 1

// SetSeed 乱数シードを設定
func SetSeed(s uint32) {
	seed = s
}

// Rand 0以上0x7FFFFFFFの乱数を生成
func Rand() uint32 {
	// 線形合同法: X(n+1) = (a * X(n) + c) mod m
	// a = 1103515245, c = 12345, m = 2^31
	seed = seed*1103515245 + 12345
	return (seed / 65536) & 0x7FFFFFFF
}

// RandInt 0以上n未満の整数乱数を生成
func RandInt(n int32) int32 {
	if n <= 0 {
		return 0
	}
	return int32(Rand()) % n
}

// RandRange min以上max未満の整数乱数を生成
func RandRange(min, max int32) int32 {
	if min >= max {
		return min
	}
	return min + RandInt(max-min)
}

// RandFixed 0.0以上1.0未満の固定小数点乱数を生成
func RandFixed() Fixed {
	return Fixed(Rand() & 0xFFFF) // 0-65535の範囲
}

// RandFixedRange min以上max未満の固定小数点乱数を生成
func RandFixedRange(min, max Fixed) Fixed {
	if min >= max {
		return min
	}
	range_ := max - min
	return min + RandFixed().Mul(range_).Div(NewFixed(65536))
}

// RandBool ランダムなbool値を生成
func RandBool() bool {
	return (Rand() & 1) == 1
}
