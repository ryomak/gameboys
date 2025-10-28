# GBA共通ライブラリ

GameBoy Advance (GBA) ゲーム開発のための共通ライブラリです。
TinyGoでコンパイルし、複数のゲームプロジェクトで共有できます。

## パッケージ構成

### gba/display
ディスプレイ制御とVBlank管理

**主な機能:**
- `SetMode(mode)` - ディスプレイモード設定 (Mode 0-5)
- `EnableLayers(layers)` - レイヤー表示制御
- `WaitForVBlank()` - VBlank待機
- `IsVBlank()` - VBlank期間判定
- `SetFrameBuffer(frame)` - フレームバッファ切り替え（Mode 4, 5）

**使用例:**
```go
import "github.com/ryomak/gameboys/common/gba/display"

display.SetMode(display.Mode3)
display.EnableLayers(display.EnableBG2)
display.WaitForVBlank()
```

### gba/graphics
グラフィックス描画機能

**主な機能:**
- `RGB15(r, g, b)` - 15bitカラー作成
- `DrawPixel(x, y, color)` - ピクセル描画
- `ClearScreen(color)` - 画面クリア
- `FillRect(x, y, w, h, color)` - 矩形塗りつぶし
- `DrawLine(x0, y0, x1, y1, color)` - 線描画
- `DrawRect(x, y, w, h, color)` - 矩形枠描画
- `DrawCircle(cx, cy, r, color)` - 円描画
- `FillCircle(cx, cy, r, color)` - 円塗りつぶし

**使用例:**
```go
import "github.com/ryomak/gameboys/common/gba/graphics"

graphics.ClearScreen(graphics.ColorBlack)
graphics.DrawCircle(120, 80, 20, graphics.ColorRed)
graphics.FillRect(50, 50, 100, 60, graphics.ColorBlue)
```

### gba/input
キー入力処理

**主な機能:**
- `IsKeyDown(key)` - キーが押されているか（現在）
- `NewKeyState()` - 入力状態管理を初期化
- `KeyState.Update()` - 入力状態更新（毎フレーム）
- `KeyState.IsPressed(key)` - キーが押された瞬間
- `KeyState.IsHeld(key)` - キーが押され続けている

**キー定数:**
`KeyA`, `KeyB`, `KeySelect`, `KeyStart`, `KeyUp`, `KeyDown`, `KeyLeft`, `KeyRight`, `KeyL`, `KeyR`

**使用例:**
```go
import "github.com/ryomak/gameboys/common/gba/input"

keys := input.NewKeyState()

// メインループ内
keys.Update()
if keys.IsPressed(input.KeyA) {
    // Aボタンが押された瞬間の処理
}
if keys.IsHeld(input.KeyRight) {
    // 右キーが押され続けている間の処理
}
```

### gba/memory
メモリ操作とDMA転送

**主な機能:**
- `DMA3Copy16(dst, src, count)` - 16bit DMA転送
- `DMA3Copy32(dst, src, count)` - 32bit DMA転送
- `DMA3Fill16(dst, value, count)` - 16bit値で埋める
- `DMA3Fill32(dst, value, count)` - 32bit値で埋める
- `Copy16(dst, src, count)` - CPUコピー（16bit）
- `Fill32(dst, value, count)` - CPU塗りつぶし（32bit）

**使用例:**
```go
import "github.com/ryomak/gameboys/common/gba/memory"

// 高速なメモリクリア
memory.DMA3Fill16(unsafe.Pointer(&buffer), 0, len(buffer))
```

### math
数学関数（固定小数点演算）

**主な機能:**
- `Fixed` - 固定小数点型（16.16形式）
- `NewFixed(x)` - 整数から固定小数点へ変換
- `Fixed.Mul(other)` - 乗算
- `Fixed.Div(other)` - 除算
- `Vec2` - 2次元ベクトル
- `FixedSqrt(x)` - 平方根
- `Rand()`, `RandInt(n)` - 乱数生成

**使用例:**
```go
import "github.com/ryomak/gameboys/common/math"

// 固定小数点演算
a := math.NewFixed(10)  // 10.0
b := math.NewFixed(3)   // 3.0
c := a.Div(b)           // 3.333...
result := c.ToInt()     // 3

// ベクトル演算
v1 := math.NewVec2(10, 20)
v2 := math.NewVec2(5, 10)
v3 := v1.Add(v2)        // (15, 30)

// 乱数
math.SetSeed(12345)
randomNum := math.RandInt(100)  // 0-99
```

### util
ユーティリティ関数

**主な機能:**
- `Rect` - 矩形と衝突判定
- `Circle` - 円と衝突判定
- `Point` - 点と距離計算
- `Min(a, b)`, `Max(a, b)` - 最小/最大値
- `Clamp(value, min, max)` - 範囲制限
- `Abs(x)` - 絶対値

**使用例:**
```go
import "github.com/ryomak/gameboys/common/util"

rect1 := util.NewRect(10, 10, 50, 50)
rect2 := util.NewRect(40, 40, 50, 50)

if rect1.Intersects(rect2) {
    // 衝突している
}

circle := util.NewCircle(100, 100, 20)
if circle.IntersectsRect(rect1) {
    // 円と矩形が衝突
}
```

## 使用方法

### 1. go.workの設定

プロジェクトルートで`go.work`を作成：

```go
go 1.21

use (
    ./common
    ./your-game
)
```

### 2. ゲームプロジェクトのgo.mod

```go
module github.com/ryomak/gameboys/your-game

go 1.21

require github.com/ryomak/gameboys/common v0.0.0
```

### 3. インポートして使用

```go
package main

import (
    "github.com/ryomak/gameboys/common/gba/display"
    "github.com/ryomak/gameboys/common/gba/graphics"
    "github.com/ryomak/gameboys/common/gba/input"
)

func main() {
    // 初期化
    display.SetMode(display.Mode3)
    display.EnableLayers(display.EnableBG2)

    keys := input.NewKeyState()

    x, y := 120, 80

    // メインループ
    for {
        display.WaitForVBlank()
        keys.Update()

        if keys.IsHeld(input.KeyRight) {
            x++
        }
        if keys.IsHeld(input.KeyLeft) {
            x--
        }

        graphics.ClearScreen(graphics.ColorBlack)
        graphics.FillCircle(x, y, 5, graphics.ColorRed)
    }
}
```

## ビルド

```bash
tinygo build -o game.gba -target=gameboy-advance main.go
```

## 注意事項

- **メモリ制約**: IWRAM 32KB、EWRAM 256KBと限られているため、大きな配列は避ける
- **浮動小数点なし**: `math.Fixed`型を使用した固定小数点演算を使う
- **VBlank期間**: VRAM書き込みは`WaitForVBlank()`後に行う
- **DMA活用**: 大量データ転送にはDMAを使用すると高速

## パフォーマンスTips

1. **VBlank期間を活用** - 描画はVBlank期間中に
2. **DMAで転送** - 大量のデータはDMA使用
3. **固定小数点** - 浮動小数点の代わりに`math.Fixed`
4. **メモリアライメント** - 32bit/16bit境界に注意

## ライセンス

このライブラリはGameBoy Advanceゲーム開発用の共通ライブラリです。
