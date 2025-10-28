package input

import (
	"runtime/volatile"
	"unsafe"
)

const RegKEYINPUT = 0x04000130

var KEYINPUT = (*volatile.Register16)(unsafe.Pointer(uintptr(RegKEYINPUT)))

// キー定義（負論理: 0=押下、1=未押下）
const (
	KeyA      = 1 << 0
	KeyB      = 1 << 1
	KeySelect = 1 << 2
	KeyStart  = 1 << 3
	KeyRight  = 1 << 4
	KeyLeft   = 1 << 5
	KeyUp     = 1 << 6
	KeyDown   = 1 << 7
	KeyR      = 1 << 8
	KeyL      = 1 << 9
)

// キーの組み合わせ
const (
	KeyAny = 0x03FF // すべてのキー
)

// IsKeyDown キーが押されているか（現在の状態）
func IsKeyDown(key uint16) bool {
	return (KEYINPUT.Get() & key) == 0
}

// IsKeyUp キーが離されているか（現在の状態）
func IsKeyUp(key uint16) bool {
	return (KEYINPUT.Get() & key) != 0
}

// GetKeys 現在のキー状態を取得（生の値）
func GetKeys() uint16 {
	return KEYINPUT.Get()
}

// AreAllKeysDown 指定された全てのキーが押されているか
func AreAllKeysDown(keys uint16) bool {
	return (KEYINPUT.Get() & keys) == 0
}

// IsAnyKeyDown 指定されたキーのいずれかが押されているか
func IsAnyKeyDown(keys uint16) bool {
	return (KEYINPUT.Get() & keys) != keys
}
