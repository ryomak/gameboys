# GameBoy Advance Games - Root Makefile

# ゲームディレクトリのリスト
GAMES = demo

# デフォルトターゲット: すべてのゲームをビルド
.PHONY: all
all: build-all

# すべてのゲームをビルド
.PHONY: build-all
build-all:
	@echo "Building all games..."
	@for game in $(GAMES); do \
		echo ""; \
		echo "Building $$game..."; \
		$(MAKE) -C $$game build || exit 1; \
	done
	@echo ""
	@echo "All games built successfully!"
	@echo "ROMs are in bin/ directory"

# 特定のゲームをビルド
.PHONY: build-%
build-%:
	@echo "Building $*..."
	@$(MAKE) -C $* build

# 特定のゲームを実行
.PHONY: run-%
run-%:
	@echo "Running $*..."
	@$(MAKE) -C $* run

# すべてクリーン
.PHONY: clean
clean:
	@echo "Cleaning all games..."
	@for game in $(GAMES); do \
		echo "Cleaning $$game..."; \
		$(MAKE) -C $$game clean; \
	done
	@rm -rf bin/
	@echo "Clean complete"

# binディレクトリのROMs一覧
.PHONY: list
list:
	@echo "Built ROMs in bin/:"
	@if [ -d bin ]; then \
		ls -lh bin/*.gba 2>/dev/null | awk '{print "  " $$9 " (" $$5 ")"}' || echo "  No ROMs found"; \
	else \
		echo "  bin/ directory does not exist"; \
	fi

# ディレクトリ構造を表示
.PHONY: tree
tree:
	@echo "Project structure:"
	@tree -L 2 -I 'bin' . 2>/dev/null || find . -maxdepth 2 -not -path '*/\.*' -not -path './bin/*' | sort

# 開発環境のチェック
.PHONY: check-env
check-env:
	@echo "Checking development environment..."
	@echo ""
	@echo "Go version:"
	@go version || echo "  Error: Go not found"
	@echo ""
	@echo "TinyGo version:"
	@tinygo version || echo "  Error: TinyGo not found"
	@echo ""
	@echo "mGBA emulator:"
	@if command -v mgba-qt >/dev/null 2>&1; then \
		echo "  mgba-qt found"; \
	elif command -v mgba >/dev/null 2>&1; then \
		echo "  mgba found"; \
	else \
		echo "  Warning: mGBA not found"; \
	fi
	@echo ""
	@echo "Checking Go workspace..."
	@if [ -f go.work ]; then \
		echo "  go.work found"; \
		go work sync 2>/dev/null && echo "  Workspace synced" || echo "  Warning: Could not sync workspace"; \
	else \
		echo "  Error: go.work not found"; \
	fi

# go.work を同期
.PHONY: sync
sync:
	@echo "Syncing Go workspace..."
	@go work sync
	@echo "Workspace synced"

# ヘルプ
.PHONY: help
help:
	@echo "GameBoy Advance Games - Makefile"
	@echo ""
	@echo "Available games: $(GAMES)"
	@echo ""
	@echo "Usage:"
	@echo "  make                  - Build all games"
	@echo "  make build-all        - Build all games"
	@echo "  make build-<game>     - Build specific game (e.g., make build-demo)"
	@echo "  make run-<game>       - Build and run specific game (e.g., make run-demo)"
	@echo "  make clean            - Clean all build files"
	@echo "  make list             - List built ROMs"
	@echo "  make tree             - Show project structure"
	@echo "  make check-env        - Check development environment"
	@echo "  make sync             - Sync Go workspace"
	@echo "  make help             - Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make build-demo       - Build demo game"
	@echo "  make run-demo         - Build and run demo game"
