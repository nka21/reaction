# ========================================
# ビルドステージ
# ========================================
FROM golang:1.21-alpine AS builder

# 作業ディレクトリを設定
WORKDIR /app

# 依存関係のキャッシュを活用するため、go.modとgo.sumを先にコピー
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# バイナリをビルド（静的リンクで軽量化）
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o reaction main.go

# ========================================
# 実行ステージ
# ========================================
FROM alpine:latest

# セキュリティアップデートとCA証明書をインストール
# WHY: Discord APIとのHTTPS通信にCA証明書が必要
RUN apk --no-cache add ca-certificates tzdata

# タイムゾーンを設定
ENV TZ=Asia/Tokyo

# 作業ディレクトリを設定
WORKDIR /root/

# ビルドステージからバイナリをコピー
COPY --from=builder /app/reaction .

# バイナリを実行
CMD ["./reaction"]
