package math

import "testing"

func TestVec3Add(t *testing.T) {
	v1 := NewVec3(1, 2, 3)
	v2 := NewVec3(4, 5, 6)
	result := v1.Add(v2)

	expected := NewVec3(5, 7, 9)
	if result.X != expected.X || result.Y != expected.Y || result.Z != expected.Z {
		t.Errorf("Add failed: got %v, want %v", result, expected)
	}
}

func TestVec3Sub(t *testing.T) {
	v1 := NewVec3(5, 7, 9)
	v2 := NewVec3(1, 2, 3)
	result := v1.Sub(v2)

	expected := NewVec3(4, 5, 6)
	if result.X != expected.X || result.Y != expected.Y || result.Z != expected.Z {
		t.Errorf("Sub failed: got %v, want %v", result, expected)
	}
}

func TestVec3Mul(t *testing.T) {
	v := NewVec3(2, 3, 4)
	scalar := NewFixed(2)
	result := v.Mul(scalar)

	expected := NewVec3(4, 6, 8)
	if result.X != expected.X || result.Y != expected.Y || result.Z != expected.Z {
		t.Errorf("Mul failed: got %v, want %v", result, expected)
	}
}

func TestVec3Dot(t *testing.T) {
	v1 := NewVec3(1, 2, 3)
	v2 := NewVec3(4, 5, 6)
	result := v1.Dot(v2)

	// 1*4 + 2*5 + 3*6 = 4 + 10 + 18 = 32
	expected := NewFixed(32)
	if result != expected {
		t.Errorf("Dot failed: got %d, want %d", result, expected)
	}
}

func TestVec3Cross(t *testing.T) {
	v1 := NewVec3(1, 0, 0)
	v2 := NewVec3(0, 1, 0)
	result := v1.Cross(v2)

	// i × j = k
	expected := NewVec3(0, 0, 1)
	if result.X != expected.X || result.Y != expected.Y || result.Z != expected.Z {
		t.Errorf("Cross failed: got %v, want %v", result, expected)
	}
}

func TestVec3LengthSq(t *testing.T) {
	v := NewVec3(2, 3, 6)
	result := v.LengthSq()

	// 2^2 + 3^2 + 6^2 = 4 + 9 + 36 = 49
	expected := NewFixed(49)
	if result != expected {
		t.Errorf("LengthSq failed: got %d, want %d", result, expected)
	}
}

func TestDistance3(t *testing.T) {
	v1 := NewVec3(0, 0, 0)
	v2 := NewVec3(3, 4, 0)
	result := Distance3(v1, v2)

	// sqrt(3^2 + 4^2) = sqrt(25) = 5
	expected := NewFixed(5)
	// 固定小数点演算の誤差を考慮
	diff := result.Sub(expected).Abs()
	if diff > NewFixed(1) { // 1以内の誤差を許容
		t.Errorf("Distance3 failed: got %v, want %v", result.ToFloat(), expected.ToFloat())
	}
}
