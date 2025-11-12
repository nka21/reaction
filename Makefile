.PHONY: help fmt fmt-check lint run build build-prod test clean docker-build

# デフォルトターゲット
help:
	@echo "使用可能なコマンド:"
	@echo ""
	@echo "  [開発用]"
	@echo "  make run          - アプリケーションを実行"
	@echo "  make build        - バイナリをビルド (デバッグ情報あり)"
	@echo "  make test         - テストを実行"
	@echo ""
	@echo "  [コード品質]"
	@echo "  make fmt          - コードフォーマット (自動修正)"
	@echo "  make fmt-check    - フォーマットチェック (修正しない)"
	@echo "  make lint         - 静的解析 (golangci-lint)"
	@echo ""
	@echo "  [本番用]"
	@echo "  make build-prod   - 本番用バイナリをビルド (最適化、デバッグ情報なし)"
	@echo "  make docker-build - Dockerイメージをビルド"
	@echo ""
	@echo "  [その他]"
	@echo "  make clean        - ビルド成果物を削除"
	@echo ""
	@echo "デプロイはGitHub Actionsで自動実行されます (mainブランチへのpush時)"

# コードフォーマット
fmt:
	@echo "コードをフォーマット中..."
	gofmt -s -w .
	@echo "フォーマット完了"

# フォーマットチェック (修正しない)
fmt-check:
	@echo "フォーマットチェック中..."
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "以下のファイルがフォーマットされていません:"; \
		gofmt -l .; \
		exit 1; \
	fi
	@echo "フォーマットチェック完了"

# 静的解析
lint:
	@echo "静的解析を実行中..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "golangci-lint がインストールされていません"; \
		exit 1; \
	fi
	golangci-lint run ./...
	@echo "静的解析完了"

# アプリケーションを実行
run:
	@echo "アプリケーションを起動中..."
	go run main.go

# バイナリをビルド (開発用)
build:
	@echo "バイナリをビルド中..."
	go build -o reaction main.go
	@echo "ビルド完了: ./reaction"

# 本番用バイナリをビルド (最適化)
build-prod:
	@echo "本番用バイナリをビルド中..."
	go build -ldflags="-s -w" -trimpath -o reaction main.go
	@echo "ビルド完了: ./reaction"

# テストを実行
test:
	@echo "テストを実行中..."
	go test -v ./...
	@echo "テスト完了"

# ビルド成果物を削除
clean:
	@echo "ビルド成果物を削除中..."
	rm -f reaction
	@echo "削除完了"

# Dockerイメージをビルド
docker-build:
	@echo "Dockerイメージをビルド中..."
	docker build -t reaction-bot:latest .
	@echo "ビルド完了"
