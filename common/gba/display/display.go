package display

import (
	"runtime/volatile"
	"unsafe"
)

// レジスタアドレス
const (
	RegDISPCNT  = 0x04000000
	RegDISPSTAT = 0x04000004
	RegVCOUNT   = 0x04000006
)

// ディスプレイモード
const (
	Mode0 = 0x0000
	Mode1 = 0x0001
	Mode2 = 0x0002
	Mode3 = 0x0003 // ビットマップモード
	Mode4 = 0x0004
	Mode5 = 0x0005
)

// レイヤー表示フラグ
const (
	EnableBG0 = 1 << 8
	EnableBG1 = 1 << 9
	EnableBG2 = 1 << 10
	EnableBG3 = 1 << 11
	EnableOBJ = 1 << 12
)

// その他のディスプレイ制御フラグ
const (
	FrameSelect         = 1 << 4  // フレームバッファ選択（Mode 4, 5）
	OBJVRAMMapping      = 1 << 6  // OBJ VRAM 1次元マッピング
	ForcedBlank         = 1 << 7  // 強制ブランク
	EnableWin0          = 1 << 13 // Window 0表示
	EnableWin1          = 1 << 14 // Window 1表示
	EnableOBJWin        = 1 << 15 // OBJ Window表示
)

// レジスタアクセス用の変数
var (
	DISPCNT  = (*volatile.Register16)(unsafe.Pointer(uintptr(RegDISPCNT)))
	DISPSTAT = (*volatile.Register16)(unsafe.Pointer(uintptr(RegDISPSTAT)))
	VCOUNT   = (*volatile.Register16)(unsafe.Pointer(uintptr(RegVCOUNT)))
)

// SetMode ディスプレイモードを設定
func SetMode(mode uint16) {
	DISPCNT.Set(mode)
}

// EnableLayers レイヤーを有効化
func EnableLayers(layers uint16) {
	current := DISPCNT.Get()
	DISPCNT.Set(current | layers)
}

// DisableLayers レイヤーを無効化
func DisableLayers(layers uint16) {
	current := DISPCNT.Get()
	DISPCNT.Set(current &^ layers)
}

// SetControl ディスプレイ制御レジスタに直接値を設定
func SetControl(value uint16) {
	DISPCNT.Set(value)
}

// GetControl ディスプレイ制御レジスタの値を取得
func GetControl() uint16 {
	return DISPCNT.Get()
}

// WaitForVBlank VBlank期間まで待機
func WaitForVBlank() {
	// VBlank期間が終わるまで待つ
	for VCOUNT.Get() >= 160 {
	}
	// VBlank期間が始まるまで待つ
	for VCOUNT.Get() < 160 {
	}
}

// IsVBlank VBlank期間中かどうか
func IsVBlank() bool {
	return VCOUNT.Get() >= 160
}

// GetVCount 現在の垂直カウンタ値を取得（0-227）
func GetVCount() uint16 {
	return VCOUNT.Get()
}

// SetFrameBuffer フレームバッファを選択（Mode 4, 5用）
func SetFrameBuffer(frame uint16) {
	current := DISPCNT.Get()
	if frame == 1 {
		DISPCNT.Set(current | FrameSelect)
	} else {
		DISPCNT.Set(current &^ FrameSelect)
	}
}

// GetFrameBuffer 現在のフレームバッファを取得
func GetFrameBuffer() uint16 {
	if (DISPCNT.Get() & FrameSelect) != 0 {
		return 1
	}
	return 0
}
