# GBA共通ライブラリ 設計書

## 概要

複数のGBAゲームで共有する共通ライブラリの設計書です。
Go Workspaceを使用して、各ゲームプロジェクトから共通機能を利用できるようにします。

## ディレクトリ構造

```
common/
├── go.mod                 # モジュール定義
├── gba/                   # GBA固有の機能
│   ├── display/          # ディスプレイ制御
│   │   ├── display.go    # モード設定、VBlank制御
│   │   ├── bg.go         # 背景レイヤー制御
│   │   └── mode3.go      # Mode 3専用機能
│   ├── graphics/         # グラフィックス描画
│   │   ├── pixel.go      # ピクセル描画
│   │   ├── shape.go      # 図形描画（線、矩形、円）
│   │   ├── sprite.go     # スプライト管理
│   │   └── color.go      # カラー変換
│   ├── input/            # 入力処理
│   │   ├── keys.go       # キー入力
│   │   └── state.go      # 入力状態管理
│   ├── memory/           # メモリ操作
│   │   ├── dma.go        # DMA転送
│   │   ├── copy.go       # メモリコピー
│   │   └── vram.go       # VRAM操作
│   └── interrupt/        # 割り込み処理
│       ├── interrupt.go  # 割り込み設定
│       └── timer.go      # タイマー
├── math/                 # 数学関数
│   ├── fixed.go          # 固定小数点演算
│   ├── vector.go         # ベクトル演算
│   ├── trig.go           # 三角関数（ルックアップテーブル）
│   └── random.go         # 乱数生成
└── util/                 # ユーティリティ
    ├── text.go           # テキスト描画（ビットマップフォント）
    └── collision.go      # 当たり判定
```

## モジュール設計

### 1. gba/display パッケージ

ディスプレイ制御とVBlank管理を提供します。

#### display.go

```go
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

// レジスタアクセス用の型定義
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

// WaitForVBlank VBlank期間まで待機
func WaitForVBlank() {
    for VCOUNT.Get() >= 160 {
    }
    for VCOUNT.Get() < 160 {
    }
}

// IsVBlank VBlank期間中かどうか
func IsVBlank() bool {
    return VCOUNT.Get() >= 160
}
```

### 2. gba/graphics パッケージ

グラフィックス描画機能を提供します。

#### color.go

```go
package graphics

// RGB15 RGB値から15bitカラーに変換（各色5bit）
func RGB15(r, g, b uint8) uint16 {
    return uint16(r&0x1F) | (uint16(g&0x1F) << 5) | (uint16(b&0x1F) << 10)
}

// 基本カラー定義
var (
    ColorBlack   = RGB15(0, 0, 0)
    ColorWhite   = RGB15(31, 31, 31)
    ColorRed     = RGB15(31, 0, 0)
    ColorGreen   = RGB15(0, 31, 0)
    ColorBlue    = RGB15(0, 0, 31)
    ColorYellow  = RGB15(31, 31, 0)
    ColorCyan    = RGB15(0, 31, 31)
    ColorMagenta = RGB15(31, 0, 31)
)
```

#### pixel.go (Mode 3用)

```go
package graphics

import "unsafe"

const (
    ScreenWidth  = 240
    ScreenHeight = 160
    VRAMBase     = 0x06000000
)

var VideoBuffer = (*[ScreenWidth * ScreenHeight]uint16)(unsafe.Pointer(uintptr(VRAMBase)))

// DrawPixel ピクセルを描画
func DrawPixel(x, y int, color uint16) {
    if x >= 0 && x < ScreenWidth && y >= 0 && y < ScreenHeight {
        VideoBuffer[y*ScreenWidth+x] = color
    }
}

// ClearScreen 画面をクリア
func ClearScreen(color uint16) {
    for i := range VideoBuffer {
        VideoBuffer[i] = color
    }
}

// FillRect 矩形を塗りつぶす
func FillRect(x, y, width, height int, color uint16) {
    for j := 0; j < height; j++ {
        for i := 0; i < width; i++ {
            DrawPixel(x+i, y+j, color)
        }
    }
}
```

#### shape.go

```go
package graphics

// DrawLine 線を描画（ブレゼンハムのアルゴリズム）
func DrawLine(x0, y0, x1, y1 int, color uint16) {
    dx := abs(x1 - x0)
    dy := abs(y1 - y0)
    sx := -1
    if x0 < x1 {
        sx = 1
    }
    sy := -1
    if y0 < y1 {
        sy = 1
    }
    err := dx - dy

    for {
        DrawPixel(x0, y0, color)
        if x0 == x1 && y0 == y1 {
            break
        }
        e2 := 2 * err
        if e2 > -dy {
            err -= dy
            x0 += sx
        }
        if e2 < dx {
            err += dx
            y0 += sy
        }
    }
}

// DrawRect 矩形の枠を描画
func DrawRect(x, y, width, height int, color uint16) {
    // 上下の線
    for i := 0; i < width; i++ {
        DrawPixel(x+i, y, color)
        DrawPixel(x+i, y+height-1, color)
    }
    // 左右の線
    for j := 0; j < height; j++ {
        DrawPixel(x, y+j, color)
        DrawPixel(x+width-1, y+j, color)
    }
}

// DrawCircle 円を描画（中点円描画アルゴリズム）
func DrawCircle(cx, cy, radius int, color uint16) {
    x := radius
    y := 0
    err := 0

    for x >= y {
        DrawPixel(cx+x, cy+y, color)
        DrawPixel(cx+y, cy+x, color)
        DrawPixel(cx-y, cy+x, color)
        DrawPixel(cx-x, cy+y, color)
        DrawPixel(cx-x, cy-y, color)
        DrawPixel(cx-y, cy-x, color)
        DrawPixel(cx+y, cy-x, color)
        DrawPixel(cx+x, cy-y, color)

        if err <= 0 {
            y++
            err += 2*y + 1
        }
        if err > 0 {
            x--
            err -= 2*x + 1
        }
    }
}

func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}
```

### 3. gba/input パッケージ

キー入力処理を提供します。

#### keys.go

```go
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

// IsKeyDown キーが押されているか（現在の状態）
func IsKeyDown(key uint16) bool {
    return (KEYINPUT.Get() & key) == 0
}

// IsKeyUp キーが離されているか（現在の状態）
func IsKeyUp(key uint16) bool {
    return (KEYINPUT.Get() & key) != 0
}
```

#### state.go

```go
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
    return (ks.current&key) == 0 && (ks.previous&key) == 0
}
```

### 4. gba/memory パッケージ

メモリ操作とDMA転送を提供します。

#### dma.go

```go
package memory

import (
    "runtime/volatile"
    "unsafe"
)

// DMAレジスタ
const (
    RegDMA3SAD = 0x040000D4 // DMA3 転送元
    RegDMA3DAD = 0x040000D8 // DMA3 転送先
    RegDMA3CNT = 0x040000DC // DMA3 制御
)

// DMA制御フラグ
const (
    DMAEnable    = 1 << 31
    DMA32        = 1 << 26 // 32bitモード
    DMA16        = 0       // 16bitモード
    DMAImmediate = 0       // 即座に転送
)

// DMA3Copy DMA3を使ってメモリコピー（高速）
func DMA3Copy(src, dst unsafe.Pointer, count uint32, mode uint32) {
    sadReg := (*volatile.Register32)(unsafe.Pointer(uintptr(RegDMA3SAD)))
    dadReg := (*volatile.Register32)(unsafe.Pointer(uintptr(RegDMA3DAD)))
    cntReg := (*volatile.Register32)(unsafe.Pointer(uintptr(RegDMA3CNT)))

    sadReg.Set(uint32(uintptr(src)))
    dadReg.Set(uint32(uintptr(dst)))
    cntReg.Set(count | mode | DMAEnable)
}

// FillMemory32 32bit値でメモリを埋める
func FillMemory32(dst unsafe.Pointer, value uint32, count uint32) {
    // DMAの固定転送元モードを使用
    const DMAFixedSrc = 2 << 23
    src := (*volatile.Register32)(unsafe.Pointer(uintptr(0x03000000)))
    src.Set(value)

    DMA3Copy(unsafe.Pointer(src), dst, count, DMA32|DMAFixedSrc)
}
```

### 5. math パッケージ

数学関数を提供します。GBAには浮動小数点ユニットがないため、固定小数点演算を使用します。

#### fixed.go

```go
package math

// Fixed 固定小数点数（16.16形式）
type Fixed int32

const FixedShift = 16

// NewFixed 整数から固定小数点数に変換
func NewFixed(x int32) Fixed {
    return Fixed(x << FixedShift)
}

// ToInt 固定小数点数から整数に変換
func (f Fixed) ToInt() int32 {
    return int32(f >> FixedShift)
}

// Mul 乗算
func (f Fixed) Mul(other Fixed) Fixed {
    return Fixed((int64(f) * int64(other)) >> FixedShift)
}

// Div 除算
func (f Fixed) Div(other Fixed) Fixed {
    return Fixed((int64(f) << FixedShift) / int64(other))
}

// Add 加算
func (f Fixed) Add(other Fixed) Fixed {
    return f + other
}

// Sub 減算
func (f Fixed) Sub(other Fixed) Fixed {
    return f - other
}
```

#### vector.go

```go
package math

// Vec2 2次元ベクトル
type Vec2 struct {
    X, Y Fixed
}

// NewVec2 ベクトルを作成
func NewVec2(x, y int32) Vec2 {
    return Vec2{
        X: NewFixed(x),
        Y: NewFixed(y),
    }
}

// Add ベクトル加算
func (v Vec2) Add(other Vec2) Vec2 {
    return Vec2{
        X: v.X + other.X,
        Y: v.Y + other.Y,
    }
}

// Sub ベクトル減算
func (v Vec2) Sub(other Vec2) Vec2 {
    return Vec2{
        X: v.X - other.X,
        Y: v.Y - other.Y,
    }
}

// Mul スカラー倍
func (v Vec2) Mul(scalar Fixed) Vec2 {
    return Vec2{
        X: v.X.Mul(scalar),
        Y: v.Y.Mul(scalar),
    }
}
```

### 6. util パッケージ

ユーティリティ機能を提供します。

#### collision.go

```go
package util

// Rect 矩形
type Rect struct {
    X, Y, Width, Height int
}

// Intersects 矩形同士の衝突判定
func (r Rect) Intersects(other Rect) bool {
    return r.X < other.X+other.Width &&
        r.X+r.Width > other.X &&
        r.Y < other.Y+other.Height &&
        r.Y+r.Height > other.Y
}

// Contains 点が矩形内にあるか
func (r Rect) Contains(x, y int) bool {
    return x >= r.X && x < r.X+r.Width &&
        y >= r.Y && y < r.Y+r.Height
}
```

## 使用例

### ゲームプロジェクトでの使用

```go
package main

import (
    "common/gba/display"
    "common/gba/graphics"
    "common/gba/input"
)

func main() {
    // Mode 3（ビットマップモード）を設定
    display.SetMode(display.Mode3)
    display.EnableLayers(display.EnableBG2)

    // 入力状態管理を初期化
    keys := input.NewKeyState()

    // プレイヤー位置
    px, py := 120, 80

    // メインループ
    for {
        display.WaitForVBlank()

        // 入力更新
        keys.Update()

        // キー入力で移動
        if keys.IsHeld(input.KeyUp) {
            py--
        }
        if keys.IsHeld(input.KeyDown) {
            py++
        }
        if keys.IsHeld(input.KeyLeft) {
            px--
        }
        if keys.IsHeld(input.KeyRight) {
            px++
        }

        // 画面クリア
        graphics.ClearScreen(graphics.ColorBlack)

        // プレイヤーを描画（赤い円）
        graphics.DrawCircle(px, py, 5, graphics.ColorRed)
    }
}
```

## パフォーマンス考慮事項

1. **VBlank期間にVRAM書き込み**: 描画はVBlank期間中に行う
2. **DMA活用**: 大量データ転送にはDMAを使用
3. **固定小数点演算**: 浮動小数点の代わりに固定小数点を使用
4. **ルックアップテーブル**: 三角関数などは事前計算テーブルを使用
5. **IWRAM配置**: クリティカルな関数はIWRAMに配置

## 今後の拡張

- スプライト管理システム
- タイルベース背景システム
- サウンド再生機能
- テキスト描画システム
- アニメーション管理
- パーティクルシステム

## 参考資料

- [Tonc: GBA Programming Tutorial](https://www.coranac.com/tonc/text/)
- [GBATEK - Technical Reference](http://problemkaputt.de/gbatek.htm)
- [TinyGo Documentation](https://tinygo.org/docs/)
