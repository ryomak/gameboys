# GameBoy Advance 開発ドキュメント

このディレクトリには、Go言語（TinyGo）を使ったGameBoy Advance（GBA）ゲーム開発のためのドキュメントが含まれています。

## ドキュメント一覧

### 1. [hardware-spec.md](./hardware-spec.md)
GBAのハードウェア仕様書です。以下の情報が含まれています：

- CPU仕様（ARM7TDMI、16.78 MHz）
- メモリ構成（IWRAM、EWRAM、VRAM等）
- グラフィックスシステム（解像度、カラー、描画モード）
- サウンドシステム
- ROM形式とヘッダ構造

### 2. [memory-map.md](./memory-map.md)
GBAのメモリマップとI/Oレジスタの詳細情報です：

- メモリマップ全体図
- 各メモリ領域のアドレスとサイズ
- 主要I/Oレジスタの説明
  - ディスプレイ制御レジスタ
  - DMA転送制御レジスタ
  - キー入力レジスタ
  - タイマー、割り込み制御
- VRAMとパレットRAMの構成
- OAM（スプライト属性）の構造

### 3. [tinygo-implementation.md](./tinygo-implementation.md)
TinyGoを使ったGBA開発の実装方針です：

- TinyGoの環境構築方法
- ビルド方法とコマンド
- プロジェクト構造（Go Workspace使用）
- 基本的なGBAプログラムの構造
- メモリアクセス方法（unsafe.Pointer、volatileパッケージ）
- グラフィックス描画の実装例
- 入力処理の実装例
- パフォーマンス最適化の考慮事項
- デバッグ方法

### 4. [common-design.md](./common-design.md)
共通ライブラリの設計書です：

- 共通ライブラリのディレクトリ構造
- 各パッケージの詳細設計
  - `gba/display`: ディスプレイ制御、VBlank管理
  - `gba/graphics`: 描画機能（ピクセル、図形、色変換）
  - `gba/input`: キー入力、入力状態管理
  - `gba/memory`: DMA転送、メモリ操作
  - `math`: 固定小数点演算、ベクトル演算
  - `util`: 当たり判定などのユーティリティ
- サンプルコード
- パフォーマンス考慮事項

## 開発の始め方

1. **環境構築**: `tinygo-implementation.md` の「環境構築」セクションを参照
2. **ハードウェア理解**: `hardware-spec.md` と `memory-map.md` でGBAの仕様を把握
3. **共通ライブラリ**: `common-design.md` を基に共通ライブラリを実装
4. **ゲーム開発**: 各ゲームディレクトリで共通ライブラリを使用してゲームを作成

## プロジェクト構造

```
.
├── docs/                      # 本ドキュメント群
│   ├── README.md             # このファイル
│   ├── hardware-spec.md      # ハードウェア仕様
│   ├── memory-map.md         # メモリマップ
│   ├── tinygo-implementation.md  # 実装方針
│   └── common-design.md      # 共通ライブラリ設計
├── bin/                      # ビルド成果物（.gbaファイル）
├── common/                   # 共通ライブラリ
│   ├── go.mod
│   ├── gba/                 # GBA固有機能
│   ├── math/                # 数学関数
│   └── util/                # ユーティリティ
├── {game1}/                 # ゲーム1
│   ├── docs/               # ゲーム仕様
│   ├── main.go            # エントリーポイント
│   ├── go.mod
│   └── Makefile
└── go.work                  # Goワークスペース設定
```

## ビルド方法

各ゲームディレクトリで以下のコマンドを実行：

```bash
# ビルド
tinygo build -o ../../bin/game.gba -target=gameboy-advance main.go

# ビルド＆実行
tinygo run -target=gameboy-advance main.go
```

または、Makefileがある場合：

```bash
make build    # ビルドのみ
make run      # ビルド＆エミュレータで実行
```

## テスト環境

- **エミュレータ**: [mGBA](https://mgba.io/) を推奨
- **実機**: フラッシュカートリッジを使用して実機でテスト可能

## 参考資料

- [Tonc: GBA Programming Tutorial](https://www.coranac.com/tonc/text/)
- [GBATEK - GBA/NDS Technical Info](http://problemkaputt.de/gbatek.htm)
- [TinyGo - GameBoy Advance](https://tinygo.org/docs/reference/microcontrollers/gameboy-advance/)
- [awesome-gbadev](https://github.com/gbadev-org/awesome-gbadev)
- [Game Boy Advance Architecture](https://www.copetti.org/writings/consoles/game-boy-advance/)

## 次のステップ

1. 共通ライブラリ（`common/`）の実装
2. 最初のゲームプロジェクトの作成
3. Go Workspaceの設定（`go.work`）
4. サンプルゲームの作成とテスト
