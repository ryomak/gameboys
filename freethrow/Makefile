# フリースローゲーム用 Makefile

GAME_NAME = freethrow
OUTPUT = ../bin/$(GAME_NAME).gba
TINYGO = tinygo
TARGET = gameboy-advance
MAIN = main.go

# デフォルトターゲット
.PHONY: all
all: build

# ビルド
.PHONY: build
build:
	@echo "Building $(GAME_NAME)..."
	@mkdir -p ../bin
	$(TINYGO) build -o $(OUTPUT) -target=$(TARGET) -opt=z $(MAIN)
	@echo "Build complete: $(OUTPUT)"

# ビルドして即座にエミュレータで実行
.PHONY: run
run: build
	@echo "Running $(GAME_NAME) with mGBA..."
	@if command -v mgba-qt >/dev/null 2>&1; then \
		mgba-qt $(OUTPUT); \
	elif command -v mgba >/dev/null 2>&1; then \
		mgba $(OUTPUT); \
	else \
		echo "Error: mGBA not found. Please install mGBA."; \
		exit 1; \
	fi

# クイックビルド（最適化なし、高速コンパイル）
.PHONY: quick
quick:
	@echo "Quick building $(GAME_NAME)..."
	@mkdir -p ../bin
	$(TINYGO) build -o $(OUTPUT) -target=$(TARGET) $(MAIN)
	@echo "Quick build complete: $(OUTPUT)"

# クリーンアップ
.PHONY: clean
clean:
	@echo "Cleaning $(GAME_NAME)..."
	@rm -f $(OUTPUT)
	@echo "Clean complete"

# サイズ確認
.PHONY: size
size: build
	@echo "ROM size:"
	@ls -lh $(OUTPUT) | awk '{print $$5 " " $$9}'
	@echo ""
	@echo "Detailed size:"
	@du -h $(OUTPUT)

# ヘルプ
.PHONY: help
help:
	@echo "GBA Free Throw Game Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make build    - Build the game ROM"
	@echo "  make run      - Build and run with mGBA emulator"
	@echo "  make quick    - Quick build (no optimization)"
	@echo "  make clean    - Remove built files"
	@echo "  make size     - Show ROM size"
	@echo "  make help     - Show this help message"
