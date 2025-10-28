package math

import "testing"

func TestProjectSimple(t *testing.T) {
	// 画面中央に近い位置のオブジェクト
	pos := NewVec3Fixed(0, 0, NewFixed(100))
	screenWidth := int32(240)
	screenHeight := int32(160)
	baseDepth := NewFixed(200)

	result := ProjectSimple(pos, screenWidth, screenHeight, baseDepth)

	if !result.Visible {
		t.Error("Object should be visible")
	}

	// 画面中央付近に表示されるはず
	if result.ScreenX < 100 || result.ScreenX > 140 {
		t.Errorf("ScreenX out of expected range: %d", result.ScreenX)
	}

	if result.ScreenY < 70 || result.ScreenY > 90 {
		t.Errorf("ScreenY out of expected range: %d", result.ScreenY)
	}

	// スケールは1未満になるはず（遠いので）
	if result.Scale > FixedOne {
		t.Errorf("Scale should be less than 1: %v", result.Scale.ToFloat())
	}
}

func TestProjectSimple_BehindCamera(t *testing.T) {
	// カメラより後ろ
	pos := NewVec3Fixed(0, 0, NewFixed(-10))
	screenWidth := int32(240)
	screenHeight := int32(160)
	baseDepth := NewFixed(200)

	result := ProjectSimple(pos, screenWidth, screenHeight, baseDepth)

	if result.Visible {
		t.Error("Object behind camera should not be visible")
	}
}

func TestCalculateSpriteScale(t *testing.T) {
	baseDistance := NewFixed(100)
	baseScale := int32(256) // 通常サイズ

	tests := []struct {
		name     string
		distance Fixed
		wantMin  int32
		wantMax  int32
	}{
		{
			name:     "Same distance",
			distance: NewFixed(100),
			wantMin:  200,
			wantMax:  300,
		},
		{
			name:     "Twice as far",
			distance: NewFixed(200),
			wantMin:  100,
			wantMax:  150,
		},
		{
			name:     "Half distance",
			distance: NewFixed(50),
			wantMin:  400,
			wantMax:  512,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateSpriteScale(tt.distance, baseDistance, baseScale)

			if result < tt.wantMin || result > tt.wantMax {
				t.Errorf("CalculateSpriteScale() = %d, want between %d and %d",
					result, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestCompareDepth(t *testing.T) {
	v1 := NewVec3(0, 0, 100) // 遠い
	v2 := NewVec3(0, 0, 50)  // 近い

	result := CompareDepth(v1, v2)
	if result != -1 {
		t.Errorf("CompareDepth failed: v1 should be 'less' (drawn first) than v2, got %d", result)
	}

	result = CompareDepth(v2, v1)
	if result != 1 {
		t.Errorf("CompareDepth failed: v2 should be 'greater' than v1, got %d", result)
	}

	v3 := NewVec3(0, 0, 100)
	result = CompareDepth(v1, v3)
	if result != 0 {
		t.Errorf("CompareDepth failed: v1 and v3 should be equal, got %d", result)
	}
}

func TestCamera_Project(t *testing.T) {
	camera := NewCamera()
	camera.Position = NewVec3(0, 0, -5)
	camera.Target = NewVec3(0, 0, 0)

	// カメラの前方にあるオブジェクト
	worldPos := NewVec3(0, 0, 2)
	screenWidth := int32(240)
	screenHeight := int32(160)

	result := camera.Project(worldPos, screenWidth, screenHeight)

	if !result.Visible {
		t.Error("Object in front of camera should be visible")
	}

	// 画面中央付近に表示されるはず
	centerX := screenWidth / 2
	centerY := screenHeight / 2

	if result.ScreenX < centerX-20 || result.ScreenX > centerX+20 {
		t.Errorf("ScreenX should be near center: got %d, want near %d", result.ScreenX, centerX)
	}

	if result.ScreenY < centerY-20 || result.ScreenY > centerY+20 {
		t.Errorf("ScreenY should be near center: got %d, want near %d", result.ScreenY, centerY)
	}
}
