# Reaction Bot - 設計ドキュメント

## プロジェクト概要

**プロジェクト名**: Reaction Bot
**目的**: Discord上で特定のカスタム絵文字のリアクションがつけられたメッセージを指定したチャンネルに転送するbot

## 機能仕様

### 基本動作
- 指定したカスタム絵文字のリアクションがメッセージに付けられると、そのメッセージを転送先チャンネルに転送
- リアクションが全て外されると、転送されたメッセージを削除
- カスタム絵文字IDで判定するため、サーバー固有の絵文字を使用

### 重複防止
- 同じメッセージに複数人が同じリアクションをつけても、転送は1度だけ
- 最後のリアクションが外された時のみ、転送メッセージを削除

## アーキテクチャ

### 採用アーキテクチャ
クリーンアーキテクチャの軽量版（entities/usecases/interfaces）

### ディレクトリ構成

```
reaction-bot/
├── main.go                           # エントリーポイント
├── go.mod                            # Go modules定義
├── go.sum                            # 依存関係のチェックサム
├── .env                              # 環境変数（gitignore対象）
├── .env.example                      # 環境変数のサンプル
├── .gitignore                        # Git除外設定
├── CLAUDE.md                         # このファイル
└── internal/                         # 内部パッケージ（外部からimport不可）
    ├── entities/                     # エンティティ層
    │   └── config.go                 # 設定とドメインモデル
    ├── usecases/                     # ユースケース層
    │   └── transfer_message.go       # メッセージ転送のビジネスロジック
    └── interfaces/                   # インターフェース層
        └── discord_handler.go        # Discordイベントハンドラー
```

### 各レイヤーの責務

#### entities/
- アプリケーションの設定を管理
- ドメインモデルの定義
- 環境変数の読み込みとバリデーション
- 外部ライブラリに依存しない純粋なビジネスルール

#### usecases/
- メッセージ転送のビジネスロジック
- Discord APIを使ったメッセージ取得・送信
- 転送メッセージの整形
- entitiesに依存、interfacesからは独立

#### interfaces/
- Discordイベントのハンドリング
- 外部ライブラリ（discordgo）とのやり取り
- usecasesの呼び出し

#### main.go
- アプリケーションのエントリーポイント
- 依存関係の注入（DI）
- Discord botの起動とシャットダウン処理

## 環境変数

環境変数の設定方法は [.env.example](.env.example) を参照してください。

## 開発ワークフロー

### Makefileの使用
このプロジェクトではMakefileを使用してビルドやテストを管理しています。

**重要**: ビルドやテストなどのコマンドを実行する際は、必ずMakefileの内容を確認し、定義されているターゲットを使用してください。

主要なコマンド:
- `make build` - バイナリをビルド（開発用）
- `make test` - テストを実行
- `make fmt` - コードフォーマット
- `make lint` - 静的解析
- `make run` - アプリケーションを実行

詳細は `make help` で確認できます。

## 実装規約

### Go言語ベストプラクティス

#### 命名規則
- **パッケージ名**: 小文字、単数形、短く簡潔（`entities`, `usecases`, `interfaces`）
- **変数名**: キャメルケース、明示的で説明的な名前
  - 一般的で伝わる略語はOK: `cfg`, `msg`, `ch`, `ctx`, `err`, `id`
  - 伝わりにくい略語はNG: `fwdCh`, `rctMap` など
  - **WHY**: 変数はスコープが限定的なため、略語による可読性の低下は少ない
- **関数名**: キャメルケース、動詞で始まる、**略語を使わず完全な単語を使用**
  - プライベート: `loadConfiguration()`, `validateReactionMapping()`
  - パブリック: `TransferMessage()`, `HandleReactionAdd()`
  - 良い例: `DeleteTransferredMessage()` - 検索しやすく、何をする関数か明確
  - 悪い例: `DeleteTransferredMsg()` - 検索で引っかからない可能性がある
  - **WHY**: 関数名は外部APIとして公開され、コードベース全体で検索される。略語を使うと検索性が低下し、他の開発者が機能を見つけにくくなる
- **メソッド名**: 関数名と同様、**略語を使わず完全な単語を使用**
  - 良い例: `GetTriggerReactionCount()`, `IsTriggerReactionEmoji()`
  - 悪い例: `GetTriggerReactionCnt()`, `IsTriggerReactionEmj()`
- **定数**: 大文字スネークケース
  - `DEFAULT_REACTION_EMOJI`, `MAX_RETRY_COUNT`
- **構造体**: パスカルケース、明示的な名前
  - `Config`, `ReactionMapping`, `TransferMessageUseCase`
- **インターフェース**: パスカルケース、`-er`で終わる（可能な場合）
  - `MessageTransferrer`, `ConfigLoader`

#### コーディング規約
- **エラーハンドリング**: 全てのエラーを適切に処理、ログ出力
- **コメント**: 日本語で記載
  - パブリック関数・構造体には必ずコメントを記載
  - **関数コメント**:
    - **パブリック関数（大文字始まり）**: 必須 - `// 関数名 - 説明` の形式で記載
      ```go
      // TransferMessage - 転送メッセージを作成し、転送メッセージIDを保存
      func (uc *TransferMessageUseCase) TransferMessage(...)
      ```
    - **プライベート関数（小文字始まり）**: 関数名が明示的で処理が自明な場合は省略可
      - ただし、複雑なロジックや初見で分かりにくい処理には記載すること
      ```go
      // コメント不要 - 名前で処理内容が明確
      func (h *DiscordHandler) isTriggerReactionEmoji(emoji discordgo.Emoji) bool
      ```
    - **WHY**: Goのドキュメント生成では `// 関数名[半角スペース]説明` で関数とコメントを紐付けるが、半角スペースだけだと区切りが分かりにくいため、ハイフンを入れて視認性を高める
  - 初見で伝わりにくい実装には `WHY: ` プレフィックスをつけて理由を説明
    ```go
    // WHY: Bot自身のリアクションに反応すると無限ループになるため除外
    if r.UserID == s.State.User.ID {
        return
    }
    ```
- **ログ出力**: 以下のルールで統一する
  - **エラーログ**: `log.Printf()` を使用し、エラー詳細を `%v` で出力
    ```go
    log.Printf("メッセージ取得に失敗: %v", err)
    ```
  - **情報ログ（複数の変数を含む）**: `log.Printf()` を使用
    ```go
    log.Printf("メッセージ %s をチャンネル %s に転送しました", msgID, channelID)
    ```
  - **情報ログ（単純な文字列のみ）**: `log.Println()` を使用
    ```go
    log.Println("bot起動完了")
    ```
  - **WHY**: `Printf` で統一することで、将来的な構造化ログへの移行が容易になる
- **importの順序**:
  1. 標準ライブラリ
  2. サードパーティライブラリ
  3. 内部パッケージ
- **依存関係の注入**: コンストラクタパターンを使用
- **エラーメッセージ**: 小文字で始める、末尾にピリオドをつけない

#### テスト
- `*_test.go`ファイルとして各ファイルと同じディレクトリに配置
- 各レイヤーごとにテストを記述
- テーブル駆動テストを活用

## 依存ライブラリ

依存ライブラリの詳細は [go.mod](go.mod) を参照してください。

主要な依存:
- `github.com/bwmarrin/discordgo` - Discord API クライアント
- `github.com/joho/godotenv` - 環境変数の読み込み
