# TinyGoを使ったGBA開発 実装方針

## TinyGoとは

TinyGoは、組み込みシステム、WebAssembly、マイコン向けに最適化されたGo言語の代替コンパイラです。
GameBoy Advance (ARM7TDMI)をターゲットとしてサポートしており、Go言語でGBAゲームを開発できます。

## 環境構築

### 必要なツール

1. **TinyGo**: GBAターゲットへのコンパイラ
2. **mGBA**: GBAエミュレータ（テスト用）
3. **gbafix**: ROMヘッダー修正ツール（TinyGoに含まれる）

### インストール

```bash
# macOS
brew install tinygo

# mGBAエミュレータ
brew install mgba

# または公式サイトからダウンロード
# https://mgba.io
```

## ビルド方法

### 基本的なビルドコマンド

```bash
# GBAのROMファイル（.gba）を生成
tinygo build -o game.gba -target=gameboy-advance main.go

# ビルドして即座にエミュレータで実行
tinygo run -target=gameboy-advance main.go
```

### ビルドオプション

```bash
# 最適化レベルを指定（-opt=z でサイズ最適化）
tinygo build -o game.gba -target=gameboy-advance -opt=z main.go

# スタックサイズを指定
tinygo build -o game.gba -target=gameboy-advance -stack-size=4096 main.go
```

## プロジェクト構造

CLAUDE.mdで定義された構造に従います：

```
.
├── docs/                    # GBA開発ドキュメント
│   ├── hardware-spec.md    # ハードウェア仕様
│   ├── memory-map.md       # メモリマップ
│   ├── tinygo-implementation.md  # 本ドキュメント
│   └── common-design.md    # 共通処理設計
├── bin/                    # ビルド成果物（.gbaファイル）
├── common/                 # 共通ライブラリ（go.work管理）
│   ├── go.mod
│   ├── gba/               # GBA関連の共通関数
│   │   ├── display.go    # ディスプレイ制御
│   │   ├── input.go      # 入力処理
│   │   ├── memory.go     # メモリ操作
│   │   └── graphics.go   # グラフィックス処理
│   └── math/             # 数学関数
│       └── vector.go     # ベクトル演算
├── game1/                # ゲーム1
│   ├── docs/            # ゲーム仕様
│   ├── main.go          # エントリーポイント
│   ├── go.mod
│   └── Makefile         # ビルドスクリプト
├── game2/                # ゲーム2
│   ├── docs/
│   ├── main.go
│   ├── go.mod
│   └── Makefile
└── go.work               # Goワークスペース設定
```

## Go Workspaceの設定

複数のゲームプロジェクトで共通ライブラリを共有するために`go.work`を使用します。

### go.work ファイル

```go
go 1.21

use (
    ./common
    ./game1
    ./game2
)
```

これにより、各ゲームディレクトリから`common`モジュールを参照できます。

## 基本的なGBAプログラムの構造

### エントリーポイント

```go
package main

import (
    "unsafe"
)

// メインループ
func main() {
    // 初期化処理
    initDisplay()

    // ゲームループ
    for {
        waitForVBlank()
        update()
        render()
    }
}

// VBlank待機
func waitForVBlank() {
    // VCOUNTレジスタ監視
    // 実装は common/gba/ に配置
}
```

## メモリアクセス

GBAのハードウェアレジスタやメモリに直接アクセスする必要があります。

### unsafe.Pointerを使用したメモリアクセス

```go
package main

import "unsafe"

// ディスプレイ制御レジスタのアドレス
const (
    DISPCNT = 0x04000000
    VCOUNT  = 0x04000006
)

// レジスタへの書き込み
func setDisplayControl(value uint16) {
    *(*uint16)(unsafe.Pointer(uintptr(DISPCNT))) = value
}

// レジスタからの読み込み
func getVCount() uint16 {
    return *(*uint16)(unsafe.Pointer(uintptr(VCOUNT)))
}
```

### volatileパッケージの活用

TinyGoには`volatile`パッケージがあり、より安全にメモリアクセスできます。

```go
import "runtime/volatile"

type DisplayControl struct {
    value volatile.Register16
}

var DISPCNT = (*DisplayControl)(unsafe.Pointer(uintptr(0x04000000)))

func setMode3() {
    DISPCNT.value.Set(0x0003) // Mode 3
}
```

## グラフィックス描画

### Mode 3（ビットマップモード）の例

```go
package main

import "unsafe"

const (
    SCREEN_WIDTH  = 240
    SCREEN_HEIGHT = 160
    VRAM          = 0x06000000
)

// VRAMへの直接アクセス
var videoBuffer = (*[SCREEN_WIDTH * SCREEN_HEIGHT]uint16)(unsafe.Pointer(uintptr(VRAM)))

// ピクセル描画
func drawPixel(x, y int, color uint16) {
    if x >= 0 && x < SCREEN_WIDTH && y >= 0 && y < SCREEN_HEIGHT {
        videoBuffer[y*SCREEN_WIDTH+x] = color
    }
}

// RGB値から16bitカラーへ変換
func rgb15(r, g, b uint8) uint16 {
    return uint16(r&0x1F) | (uint16(g&0x1F) << 5) | (uint16(b&0x1F) << 10)
}

// 画面クリア
func clearScreen(color uint16) {
    for i := range videoBuffer {
        videoBuffer[i] = color
    }
}
```

## 入力処理

```go
package main

import "unsafe"

const KEYINPUT = 0x04000130

const (
    KEY_A      = 1 << 0
    KEY_B      = 1 << 1
    KEY_SELECT = 1 << 2
    KEY_START  = 1 << 3
    KEY_RIGHT  = 1 << 4
    KEY_LEFT   = 1 << 5
    KEY_UP     = 1 << 6
    KEY_DOWN   = 1 << 7
    KEY_R      = 1 << 8
    KEY_L      = 1 << 9
)

// キー状態取得（0=押下、1=未押下）
func getKeys() uint16 {
    return *(*uint16)(unsafe.Pointer(uintptr(KEYINPUT)))
}

// キーが押されているか確認
func isKeyPressed(key uint16) bool {
    return (getKeys() & key) == 0
}
```

## パフォーマンス最適化

### 1. Thumb命令の活用

16bitバス幅のROMアクセスでは、Thumb命令の方が効率的です。
TinyGoはデフォルトでThumb命令を生成します。

### 2. IWRAM配置

頻繁にアクセスする関数やデータはIWRAMに配置すると高速化できます。

```go
//go:section .iwram
func criticalFunction() {
    // 高速実行が必要な処理
}
```

### 3. DMAの活用

大量のデータ転送にはDMAを使用します（共通ライブラリに実装予定）。

### 4. VBlank期間の活用

VRAMへの書き込みはVBlank期間中に行うと、画面のちらつきを防げます。

## デバッグ方法

### 1. mGBAのデバッガ

mGBAにはデバッガが内蔵されており、メモリやレジスタの状態を確認できます。

### 2. ログ出力

TinyGoの`println`は使用できません。デバッグには以下の方法があります：

- mGBAのログウィンドウを使用
- 画面にテキスト描画（独自実装）
- メモリの特定アドレスに値を書き込んで監視

## 制約事項

### TinyGoの制約

- リフレクション機能が制限される
- 一部の標準ライブラリが使用不可
- ゴルーチンは使用可能だが、メモリが限られる
- `fmt`パッケージは重いため非推奨

### GBAの制約

- メモリが非常に限られる（IWRAM 32KB、EWRAM 256KB）
- 浮動小数点演算がハードウェアサポートされない（固定小数点を使用）
- ファイルシステムなし（全てROMに含める）

## 推奨ライブラリ構成

共通ライブラリとして以下のモジュールを実装予定：

- `common/gba/display`: ディスプレイ制御、モード設定
- `common/gba/input`: キー入力処理
- `common/gba/graphics`: 描画関数（ピクセル、線、矩形、スプライト）
- `common/gba/memory`: メモリ操作、DMA
- `common/math`: 固定小数点演算、三角関数、ベクトル演算

## ビルドフローの自動化

各ゲームディレクトリに`Makefile`を配置して、ビルドを自動化します。

```makefile
GAME_NAME = game1
OUTPUT = ../../bin/$(GAME_NAME).gba

.PHONY: build
build:
	tinygo build -o $(OUTPUT) -target=gameboy-advance -opt=z main.go

.PHONY: run
run: build
	mgba-qt $(OUTPUT)

.PHONY: clean
clean:
	rm -f $(OUTPUT)
```

## 参考資料

- [TinyGo公式ドキュメント - Game Boy Advance](https://tinygo.org/docs/reference/microcontrollers/gameboy-advance/)
- [Learning Go by examples: part 5 - Create a GBA game in Go](https://dev.to/aurelievache/learning-go-by-examples-part-5-create-a-game-boy-advance-gba-game-in-go-5944)
- [GitHub - tinygo-org/tinygba](https://github.com/tinygo-org/tinygba)
- [GBA Development Resources](https://github.com/gbadev-org/awesome-gbadev)
