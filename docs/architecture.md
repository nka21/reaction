# Reaction Bot - アーキテクチャ

## アーキテクチャ図

```mermaid
graph TB
    subgraph "Presentation Layer"
        MAIN[main.go<br/>エントリーポイント<br/>DI・起動・シャットダウン]
    end

    subgraph "Interface Adapters Layer"
        HANDLER[DiscordHandler<br/>イベントハンドリング<br/>HandleReactionAdd<br/>HandleReactionRemove]
        GATEWAY[DiscordGateway<br/>外部API橋渡し<br/>GetMessage<br/>SendMessageWithReference<br/>DeleteMessage]
    end

    subgraph "Use Cases Layer"
        TRANSFER[TransferMessageUseCase<br/>ビジネスロジック<br/>TransferMessage<br/>DeleteTransferredMessage]
        INTERFACE[DiscordClient<br/>インターフェース定義]
    end

    subgraph "Entities Layer"
        CONFIG[Config<br/>設定・ドメインモデル<br/>LoadConfig]
    end

    subgraph "External"
        DISCORD[Discord API<br/>discordgo]
        ENV[環境変数<br/>.env]
    end

    MAIN -->|creates & injects| HANDLER
    MAIN -->|creates & injects| GATEWAY
    MAIN -->|creates & injects| TRANSFER
    MAIN -->|loads| CONFIG

    HANDLER -->|uses| TRANSFER
    HANDLER -->|uses| GATEWAY
    HANDLER -->|references| CONFIG

    GATEWAY -.->|implements| INTERFACE
    TRANSFER -->|uses| INTERFACE
    TRANSFER -->|references| CONFIG

    CONFIG -->|reads| ENV
    HANDLER <-->|discordgo| DISCORD
    GATEWAY <-->|discordgo| DISCORD

    classDef presentation fill:#ff6b6b,stroke:#c92a2a,color:#fff
    classDef adapter fill:#ffd93d,stroke:#f08700,color:#000
    classDef usecase fill:#4ecdc4,stroke:#0a9396,color:#fff
    classDef entity fill:#95e1d3,stroke:#38b000,color:#000
    classDef external fill:#e0e0e0,stroke:#666,color:#000

    class MAIN presentation
    class HANDLER,GATEWAY adapter
    class TRANSFER,INTERFACE usecase
    class CONFIG entity
    class DISCORD,ENV external
```

## レイヤー責務

### Presentation Layer
- **main.go**: エントリーポイント、依存関係の注入（DI）、Bot起動・シャットダウン

### Interface Adapters Layer
- **DiscordHandler (handlers/)**: 入力側アダプター - Discordイベントのハンドリング、UseCaseの呼び出し
- **DiscordGateway (gateways/)**: 出力側アダプター - 外部API（discordgo）との橋渡し

### Use Cases Layer
- **TransferMessageUseCase**: メッセージ転送のビジネスロジック
- **DiscordClient**: Discord API通信のインターフェース定義
- **重要**: 外部ライブラリ（discordgo）に直接依存しない

### Entities Layer
- **Config**: 設定とドメインモデル、環境変数の読み込みとバリデーション

## 依存関係の方向

```
Presentation → Interface Adapters → Use Cases → Entities
                (handlers, gateways)
```

### クリーンアーキテクチャの原則
- 各レイヤーは内側のレイヤーのみに依存
- UseCasesは**インターフェース**のみ定義、Gatewaysが**実装**
- HandlersとGatewaysは同じInterface Adapters層（同等の立場）
- Handlersがビジネスフローを統制（orchestrate）、Gatewaysを呼び出す