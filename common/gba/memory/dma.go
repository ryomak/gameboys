package memory

import (
	"runtime/volatile"
	"unsafe"
)

// DMAレジスタアドレス
const (
	// DMA0
	RegDMA0SAD = 0x040000B0
	RegDMA0DAD = 0x040000B4
	RegDMA0CNT_L = 0x040000B8
	RegDMA0CNT_H = 0x040000BA

	// DMA1
	RegDMA1SAD = 0x040000BC
	RegDMA1DAD = 0x040000C0
	RegDMA1CNT_L = 0x040000C4
	RegDMA1CNT_H = 0x040000C6

	// DMA2
	RegDMA2SAD = 0x040000C8
	RegDMA2DAD = 0x040000CC
	RegDMA2CNT_L = 0x040000D0
	RegDMA2CNT_H = 0x040000D2

	// DMA3
	RegDMA3SAD = 0x040000D4
	RegDMA3DAD = 0x040000D8
	RegDMA3CNT_L = 0x040000DC
	RegDMA3CNT_H = 0x040000DE
)

// DMA制御フラグ
const (
	DMAEnable    = 1 << 15 // DMA有効化
	DMAIRQEnable = 1 << 14 // IRQ有効化
	DMAStartNow  = 0 << 12 // 即座に開始
	DMAStartVBlank = 1 << 12 // VBlank時に開始
	DMAStartHBlank = 2 << 12 // HBlank時に開始
	DMARepeat    = 1 << 9  // リピート
	DMA32        = 1 << 10 // 32bitモード
	DMA16        = 0 << 10 // 16bitモード

	// 転送元アドレス制御
	DMASrcIncrement = 0 << 7  // 転送元アドレス増加
	DMASrcDecrement = 1 << 7  // 転送元アドレス減少
	DMASrcFixed     = 2 << 7  // 転送元アドレス固定

	// 転送先アドレス制御
	DMADstIncrement = 0 << 5  // 転送先アドレス増加
	DMADstDecrement = 1 << 5  // 転送先アドレス減少
	DMADstFixed     = 2 << 5  // 転送先アドレス固定
	DMADstReload    = 3 << 5  // 転送先アドレスリロード
)

// DMA3Copy DMA3を使ってメモリコピー（最も汎用的なDMAチャンネル）
func DMA3Copy(dst, src unsafe.Pointer, count uint32, mode uint16) {
	sadReg := (*volatile.Register32)(unsafe.Pointer(uintptr(RegDMA3SAD)))
	dadReg := (*volatile.Register32)(unsafe.Pointer(uintptr(RegDMA3DAD)))
	cntReg := (*volatile.Register32)(unsafe.Pointer(uintptr(RegDMA3CNT_L)))

	// 転送元・転送先アドレスを設定
	sadReg.Set(uint32(uintptr(src)))
	dadReg.Set(uint32(uintptr(dst)))

	// カウントと制御フラグを設定（上位16bitに制御、下位16bitにカウント）
	cntReg.Set((uint32(mode|DMAEnable) << 16) | (count & 0xFFFF))
}

// DMA3Copy16 16bitモードでメモリコピー
func DMA3Copy16(dst, src unsafe.Pointer, count uint32) {
	DMA3Copy(dst, src, count, DMA16|DMASrcIncrement|DMADstIncrement)
}

// DMA3Copy32 32bitモードでメモリコピー（より高速）
func DMA3Copy32(dst, src unsafe.Pointer, count uint32) {
	DMA3Copy(dst, src, count, DMA32|DMASrcIncrement|DMADstIncrement)
}

// DMA3Fill16 16bit値でメモリを埋める
func DMA3Fill16(dst unsafe.Pointer, value uint16, count uint32) {
	// 一時的な転送元を用意
	src := value
	DMA3Copy(dst, unsafe.Pointer(&src), count, DMA16|DMASrcFixed|DMADstIncrement)
}

// DMA3Fill32 32bit値でメモリを埋める
func DMA3Fill32(dst unsafe.Pointer, value uint32, count uint32) {
	// 一時的な転送元を用意
	src := value
	DMA3Copy(dst, unsafe.Pointer(&src), count, DMA32|DMASrcFixed|DMADstIncrement)
}

// WaitDMA3 DMA3の転送完了を待つ
func WaitDMA3() {
	cntReg := (*volatile.Register16)(unsafe.Pointer(uintptr(RegDMA3CNT_H)))
	for (cntReg.Get() & DMAEnable) != 0 {
		// DMA転送中は待機
	}
}
