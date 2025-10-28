package memory

import "unsafe"

// Copy16 16bitワード単位でメモリコピー（CPUコピー）
func Copy16(dst, src unsafe.Pointer, count int) {
	dstPtr := (*uint16)(dst)
	srcPtr := (*uint16)(src)
	for i := 0; i < count; i++ {
		*(*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(dstPtr)) + uintptr(i*2))) =
			*(*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(srcPtr)) + uintptr(i*2)))
	}
}

// Copy32 32bitワード単位でメモリコピー（CPUコピー）
func Copy32(dst, src unsafe.Pointer, count int) {
	dstPtr := (*uint32)(dst)
	srcPtr := (*uint32)(src)
	for i := 0; i < count; i++ {
		*(*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(dstPtr)) + uintptr(i*4))) =
			*(*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(srcPtr)) + uintptr(i*4)))
	}
}

// Fill16 16bit値でメモリを埋める（CPUコピー）
func Fill16(dst unsafe.Pointer, value uint16, count int) {
	dstPtr := (*uint16)(dst)
	for i := 0; i < count; i++ {
		*(*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(dstPtr)) + uintptr(i*2))) = value
	}
}

// Fill32 32bit値でメモリを埋める（CPUコピー）
func Fill32(dst unsafe.Pointer, value uint32, count int) {
	dstPtr := (*uint32)(dst)
	for i := 0; i < count; i++ {
		*(*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(dstPtr)) + uintptr(i*4))) = value
	}
}

// Zero メモリをゼロクリア（32bitモード）
func Zero(dst unsafe.Pointer, bytes int) {
	// 32bit単位でゼロクリア
	count := bytes / 4
	if count > 0 {
		Fill32(dst, 0, count)
	}
	// 余りがあれば16bitでクリア
	remainder := bytes % 4
	if remainder >= 2 {
		offset := count * 4
		Fill16(unsafe.Pointer(uintptr(dst)+uintptr(offset)), 0, 1)
	}
}
