package input

// KeyState キー入力状態管理
type KeyState struct {
	current  uint16
	previous uint16
}

// NewKeyState 入力状態管理を初期化
func NewKeyState() *KeyState {
	return &KeyState{
		current:  KEYINPUT.Get(),
		previous: KEYINPUT.Get(),
	}
}

// Update 入力状態を更新（毎フレーム呼び出す）
func (ks *KeyState) Update() {
	ks.previous = ks.current
	ks.current = KEYINPUT.Get()
}

// IsPressed キーが押された瞬間（トリガー）
func (ks *KeyState) IsPressed(key uint16) bool {
	return (ks.current&key) == 0 && (ks.previous&key) != 0
}

// IsReleased キーが離された瞬間
func (ks *KeyState) IsReleased(key uint16) bool {
	return (ks.current&key) != 0 && (ks.previous&key) == 0
}

// IsHeld キーが押され続けている
func (ks *KeyState) IsHeld(key uint16) bool {
	return (ks.current & key) == 0
}

// IsDown キーが押されている（IsHeldのエイリアス）
func (ks *KeyState) IsDown(key uint16) bool {
	return (ks.current & key) == 0
}

// IsUp キーが離されている
func (ks *KeyState) IsUp(key uint16) bool {
	return (ks.current & key) != 0
}

// GetCurrent 現在のキー状態を取得
func (ks *KeyState) GetCurrent() uint16 {
	return ks.current
}

// GetPrevious 前フレームのキー状態を取得
func (ks *KeyState) GetPrevious() uint16 {
	return ks.previous
}

// IsAnyKeyPressed いずれかのキーが押された瞬間か
func (ks *KeyState) IsAnyKeyPressed() bool {
	// 前フレームで押されていなくて、現在フレームで押されているキーがあるか
	return (ks.current & ^ks.previous & KeyAny) != 0
}

// GetPressedKeys 押された瞬間のキーを取得
func (ks *KeyState) GetPressedKeys() uint16 {
	// 前フレームで押されていなくて、現在押されているキー
	return ^ks.current & ks.previous
}

// GetReleasedKeys 離された瞬間のキーを取得
func (ks *KeyState) GetReleasedKeys() uint16 {
	// 前フレームで押されていて、現在離されているキー
	return ks.current & ^ks.previous
}
