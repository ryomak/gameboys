# GameBoy Advance Games in Go

Go言語（TinyGo）で作るGameBoy Advance（GBA）ゲーム開発プロジェクトです。

## プロジェクト概要

このプロジェクトは、TinyGoを使ってGBA向けのゲームを開発するためのフレームワークとサンプルゲーム集です。
共通ライブラリを作成し、複数のゲームで再利用できる構造になっています。

## ディレクトリ構成

```
.
├── README.md             # このファイル
├── CLAUDE.md             # プロジェクト方針
├── go.work               # Goワークスペース設定
├── Makefile              # ルートMakefile
├── docs/                 # GBA開発ドキュメント
│   ├── README.md        # ドキュメント概要
│   ├── hardware-spec.md # GBAハードウェア仕様
│   ├── memory-map.md    # メモリマップとI/Oレジスタ
│   ├── tinygo-implementation.md  # TinyGo実装方針
│   └── common-design.md # 共通ライブラリ設計
├── bin/                  # ビルド成果物（.gbaファイル）
├── common/               # 共通ライブラリ
│   ├── go.mod
│   ├── README.md
│   ├── gba/             # GBA固有機能
│   │   ├── display/    # ディスプレイ制御
│   │   ├── graphics/   # グラフィックス描画
│   │   ├── input/      # キー入力
│   │   └── memory/     # メモリ操作・DMA
│   ├── math/            # 数学関数（固定小数点演算）
│   └── util/            # ユーティリティ（衝突判定など）
└── demo/                 # デモゲーム
    ├── main.go
    ├── go.mod
    ├── Makefile
    └── docs/
        └── README.md    # ゲーム仕様

```

## 必要な環境

### 必須

- **Go**: 1.21以上
- **TinyGo**: 0.30以上（GBAターゲットサポート）

### 推奨

- **mGBA**: GBAエミュレータ（テスト用）

### インストール方法

```bash
# macOS
brew install go
brew install tinygo
brew install mgba

# Linux (Ubuntu/Debian)
# Go: https://golang.org/dl/
# TinyGo: https://tinygo.org/getting-started/install/
# mGBA: https://mgba.io/downloads.html
```

## クイックスタート

### 1. 環境確認

```bash
make check-env
```

### 2. デモゲームをビルド

```bash
make build-demo
```

ビルドが成功すると `bin/demo.gba` が生成されます。

### 3. デモゲームを実行

```bash
make run-demo
```

mGBAエミュレータでゲームが起動します。

## 使い方

### 全ゲームをビルド

```bash
make build-all
# または単に
make
```

### 特定のゲームをビルド

```bash
make build-demo
```

### 特定のゲームを実行

```bash
make run-demo
```

### ビルドファイルをクリーン

```bash
make clean
```

### ビルド済みROM一覧

```bash
make list
```

### その他のコマンド

```bash
make help  # ヘルプ表示
```

## ゲーム一覧

### demo
ボール移動デモゲーム。共通ライブラリの使い方を示すサンプルです。

**操作方法:**
- 十字キー: ボール移動
- Aボタン: ボールの色変更
- Bボタン: 星の色変更
- Start + Select: 終了

詳細は [`demo/docs/README.md`](demo/docs/README.md) を参照。

## 共通ライブラリ

`common/` ディレクトリには、GBAゲーム開発のための共通機能が実装されています。

### 主要パッケージ

- **gba/display**: ディスプレイ制御、VBlank管理
- **gba/graphics**: 描画機能（ピクセル、図形、色変換）
- **gba/input**: キー入力処理
- **gba/memory**: DMA転送、メモリ操作
- **math**: 固定小数点演算、ベクトル、乱数
- **util**: 衝突判定、ユーティリティ関数

詳細は [`common/README.md`](common/README.md) を参照。

## 新しいゲームの作り方

### 1. ゲームディレクトリを作成

```bash
mkdir your-game
mkdir your-game/docs
```

### 2. go.modを作成

```go
module github.com/ryomak/gameboys/your-game

go 1.21

require github.com/ryomak/gameboys/common v0.0.0

replace github.com/ryomak/gameboys/common => ../common
```

### 3. main.goを作成

```go
package main

import (
    "github.com/ryomak/gameboys/common/gba/display"
    "github.com/ryomak/gameboys/common/gba/graphics"
    "github.com/ryomak/gameboys/common/gba/input"
)

func main() {
    display.SetMode(display.Mode3)
    display.EnableLayers(display.EnableBG2)
    keys := input.NewKeyState()

    for {
        display.WaitForVBlank()
        keys.Update()

        graphics.ClearScreen(graphics.ColorBlack)
        // ゲームロジックをここに書く
    }
}
```

### 4. Makefileを作成

`demo/Makefile` をコピーして、`GAME_NAME` を変更。

### 5. go.workに追加

```go
use (
    ./common
    ./demo
    ./your-game  // 追加
)
```

### 6. ルートMakefileに追加

```makefile
GAMES = demo your-game
```

## ドキュメント

詳細なドキュメントは `docs/` ディレクトリにあります：

- **[hardware-spec.md](docs/hardware-spec.md)**: GBAのハードウェア仕様
- **[memory-map.md](docs/memory-map.md)**: メモリマップとI/Oレジスタ
- **[tinygo-implementation.md](docs/tinygo-implementation.md)**: TinyGoでの実装方針
- **[common-design.md](docs/common-design.md)**: 共通ライブラリの設計

## ビルド方法

### 基本的なビルド

```bash
cd your-game
tinygo build -o ../bin/your-game.gba -target=gameboy-advance main.go
```

### 最適化ビルド

```bash
tinygo build -o game.gba -target=gameboy-advance -opt=z main.go
```

### 実行

```bash
mgba-qt game.gba
```

## パフォーマンスTips

1. **VBlank期間を活用**: VRAM書き込みはVBlank期間中に
2. **DMAを使用**: 大量データ転送はDMAで高速化
3. **固定小数点演算**: `math.Fixed`を使用（浮動小数点なし）
4. **メモリ制約**: IWRAM 32KB、EWRAM 256KBに注意

## 制約事項

### TinyGoの制約
- リフレクション機能が制限される
- 一部の標準ライブラリが使用不可
- `fmt`パッケージは重い（非推奨）

### GBAの制約
- メモリが限られる（IWRAM 32KB、EWRAM 256KB）
- 浮動小数点演算のハードウェアサポートなし
- ファイルシステムなし

## トラブルシューティング

### TinyGoのビルドエラー

```bash
# TinyGoのバージョン確認
tinygo version

# GBAターゲットがサポートされているか確認
tinygo targets | grep gameboy
```

### エミュレータが起動しない

```bash
# mGBAがインストールされているか確認
which mgba-qt
which mgba

# 手動で起動
mgba-qt bin/demo.gba
```

### 共通ライブラリが見つからない

```bash
# Goワークスペースを同期
go work sync

# または環境確認
make check-env
```

## 参考資料

- [TinyGo公式ドキュメント](https://tinygo.org/docs/)
- [TinyGo GBA Support](https://tinygo.org/docs/reference/microcontrollers/gameboy-advance/)
- [Tonc: GBA Programming](https://www.coranac.com/tonc/text/)
- [GBATEK](http://problemkaputt.de/gbatek.htm)
- [awesome-gbadev](https://github.com/gbadev-org/awesome-gbadev)

## ライセンス

このプロジェクトはGameBoy Advanceゲーム開発の学習・研究用です。

## 貢献

新しいゲームや共通ライブラリの改善は歓迎します！
